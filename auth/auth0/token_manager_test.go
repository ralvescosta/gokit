package auth0

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/ralvescosta/gokit/auth/errors"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
)

const (
	validToken           = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImUxM2Y4YjE1MzU1YTRjNzA5NDEwMzUxNTJlMmFhZmY4In0.eyJ1c2VyX2RhdGEiOnsiYXBwX21ldGFkYXRhIjp7fSwiY3JlYXRlZF9hdCI6IjIwMjItMTItMDdUMDg6MDE6MjAuOTg2WiIsInVzZXJfbWV0YWRhdGEiOnt9fSwiaXNzIjoidGVzdElzc3VlciIsImF1ZCI6WyJ0ZXN0QXVkaWVuY2UiXSwic3ViIjoiNWJlODYzNTkwNzNjNDM0YmFkMmRhMzkzMjIyMmRhYmUiLCJpYXQiOjE2ODc3MTY4MTcsImp0aSI6Imp0aSIsInNjb3BlIjoic2NvcGUiLCJwZXJtaXNzaW9ucyI6WyJwZXJtaXNzaW9ucyJdfQ.lomxNYaleBZmQydZHJQHMswqJU8oPiG7Myuo8TJT9IdV3yOgRGP4m6o0BYOyY0IamLq4GRhOr003CSi-6jiY1ajLYTSKTO6xi6r8GaKAG3RMZqPBpri09fjdIuTTszd8iioqflZS9FsgTQX_nOa-1BLswBwfbZcvTLRKcVbnIa5yqU-GG2g5Vjp0ccWUimKElrHk7eWxdTXj9TBHAjfVDKAjLsHbQfenMxC-ZPUwwNNdjml84-KJ60CNSYWtEq7StX1sqaEzkgY6r4k38IcTCqEmBHYfv5MT1GwYnXMlF22MUVA40bL0BrAPCrKjveoh_azqo9Wuz1ooZtJa6mDNCA"
	invalidAudienceToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImUxM2Y4YjE1MzU1YTRjNzA5NDEwMzUxNTJlMmFhZmY4In0.eyJ1c2VyX2RhdGEiOnsiYXBwX21ldGFkYXRhIjp7fSwiY3JlYXRlZF9hdCI6IjIwMjItMTItMDdUMDg6MDE6MjAuOTg2WiIsInVzZXJfbWV0YWRhdGEiOnt9fSwiaXNzIjoidGVzdElzc3VlciIsImF1ZCI6WyJpbnZhbGlkIl0sInN1YiI6IjViZTg2MzU5MDczYzQzNGJhZDJkYTM5MzIyMjJkYWJlIiwiaWF0IjoxNjg3NzE2ODE3LCJqdGkiOiJqdGkiLCJzY29wZSI6InNjb3BlIiwicGVybWlzc2lvbnMiOlsicGVybWlzc2lvbnMiXX0.f3--ScOF-3l6Ee-OzqIKCxlkvOSr_XDq_ayGmQRG5Dq1iOAj3M_HzicriO-Azz7hZrs92E0FXdDV-y0NmDUCXgBo-MWjosL7crljXH4M0w1RH1HITlBOJhD6S5KaQHrbAW25LUg4EOF2w4BB1ZAjDGVo3l_fm6wcrRIKPCnZ3BME86IPc1FQ9RXW1UIlKtgriXVAsd76K9U2Cksspw05U4_JrmCubFTRaqA1p_Tx2NGVXkFz2JrIEhxTqiPsU8o_kAPkLifOShE-Y_o3xn1LszMJ9-v6HAQgRUU4hK0xL7-O9j6ksWyyz8Cn0uQ3b9-zzbhg3obeFp-twehJ_c4LwA"
	invalidIssuerToken   = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImUxM2Y4YjE1MzU1YTRjNzA5NDEwMzUxNTJlMmFhZmY4In0.eyJ1c2VyX2RhdGEiOnsiYXBwX21ldGFkYXRhIjp7fSwiY3JlYXRlZF9hdCI6IjIwMjItMTItMDdUMDg6MDE6MjAuOTg2WiIsInVzZXJfbWV0YWRhdGEiOnt9fSwiaXNzIjoiaW52YWxpZCIsImF1ZCI6WyJ0ZXN0QXVkaWVuY2UiXSwic3ViIjoiNWJlODYzNTkwNzNjNDM0YmFkMmRhMzkzMjIyMmRhYmUiLCJpYXQiOjE2ODc3MTY4MTcsImp0aSI6Imp0aSIsInNjb3BlIjoic2NvcGUiLCJwZXJtaXNzaW9ucyI6WyJwZXJtaXNzaW9ucyJdfQ.SwBkNb8Ri34PVgBo-MTXJne9VfWE9ezhJqOb3Fp3rNyK7KQyV2BV3maGw-usLa43jVP_bykn24h0WdpYoS0Yfl8ikLoJivRkxxHLUKCfbBjlQFLWs9wWBeKds85qesg7DxT0hwADL28qDkVWMiXTXc87Baln2_9_v25KDxGnBB5vv23D5lj1CTCUkEYHU1Ne6zOffacAHF171hVgRlIJJGtTBgToyWUu-CVKzHGPszFyv5KiNtye1xjYy_DuuOky54Gy9OLObPvnWEUTUfEO9oNua-zuYtcenT68W-VPj5zIOVCasAIKpuGpADpenrtMqMhzGQsE8rO896itXCRTug"
	ExpireidToken        = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImUxM2Y4YjE1MzU1YTRjNzA5NDEwMzUxNTJlMmFhZmY4In0.eyJ1c2VyX2RhdGEiOnsiYXBwX21ldGFkYXRhIjp7fSwiY3JlYXRlZF9hdCI6IjIwMjItMTItMDdUMDg6MDE6MjAuOTg2WiIsInVzZXJfbWV0YWRhdGEiOnt9fSwiaXNzIjoidGVzdElzc3VlciIsImF1ZCI6WyJ0ZXN0QXVkaWVuY2UiXSwic3ViIjoiNWJlODYzNTkwNzNjNDM0YmFkMmRhMzkzMjIyMmRhYmUiLCJpYXQiOjE2ODc3MTY4MTcsImp0aSI6Imp0aSIsInNjb3BlIjoic2NvcGUiLCJleHAiOjE2ODc3OTY2OTIsInBlcm1pc3Npb25zIjpbInBlcm1pc3Npb25zIl19.I5kXqObnYcJVoCDdm8yvgwqD0nNsozYeOs2JD0WYRbTMz5dyYjXW-BrNkV4Ti3n2hRVIoqCN6nQQFsSUg88ZpBG_bDlW86-iiVbGZ6ZctMuHUPdSNUS6aV-BK7AptCLm73IM_eemHpIy1bHPeVwySaSGbVFfnBA_dzJUI87dEnbt_8w53w2_izjmVkFw_niOpyJ56N5okdSo1AjpfJF8ATCkDctTzMqPrZ1uwUf3CkcScmiNrA7puqIyMhw27tdyzuc-bh9fNF9LJ4evF5fbwqINFRVRtPyhlwgX7VVGbjGlAX4Eo1QzdLkqmcUMYE0o5uKwqFuVI9bbsMcHa82zXA"

	issuer   = "testIssuer"
	audience = "testAudience"
)

func setup() *auth0nManager {
	return &auth0nManager{
		logger: logging.NewMockLogger(),
		jwtConfigs: &configs.JWTConfigs{
			Issuer:   "testIssuer",
			Audience: "testAudience",
		},
		auth0Configs: &configs.Auth0Configs{
			MillisecondsBetweenJWK: 10,
		},
		jwks:            jwks(),
		lastJWKRetrieve: time.Now().UnixMilli(),
	}
}

func TestSouldValidTokenAndReturnTheClaims(t *testing.T) {
	manager := setup()

	s, err := manager.Validate(context.Background(), validToken)

	if err != nil || s == nil {
		t.Error("should not return an error and return the claims")
	}
}

func TestSouldReturnErrInvalidIssuer(t *testing.T) {
	manager := setup()

	claims, err := manager.Validate(context.Background(), invalidIssuerToken)

	if err == nil && err != errors.ErrInvalidIssuer && claims != nil {
		t.Error("should return ErrInvalidIssuer")
	}
}

func TestSouldReturnErrExpired(t *testing.T) {
	manager := setup()

	claims, err := manager.Validate(context.Background(), invalidIssuerToken)

	if err == nil && err != errors.ErrExpired && claims != nil {
		t.Error("should return ErrExpired")
	}
}

func TestSouldReturnErrInvalidAudience(t *testing.T) {
	manager := setup()

	claims, err := manager.Validate(context.Background(), invalidAudienceToken)

	if err == nil && err != errors.ErrInvalidAudience && claims != nil {
		t.Error("should return ErrInvalidAudience")
	}
}

func TestSouldReturnErrClaimsRetrievingIfThereIsNoKey(t *testing.T) {
	manager := setup()
	manager.jwks = &jose.JSONWebKeySet{}

	claims, err := manager.Validate(context.Background(), invalidAudienceToken)

	if err == nil && err != errors.ErrClaimsRetrieving && claims != nil {
		t.Error("should return ErrClaimsRetrieving")
	}
}

func TestSouldReturnErrClaimsRetrievingIfThereIsWrongKey(t *testing.T) {
	manager := setup()
	manager.jwks = &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &rsa.PublicKey{N: big.NewInt(2048), E: 65543}}}}

	claims, err := manager.Validate(context.Background(), invalidAudienceToken)

	if err == nil && err != errors.ErrClaimsRetrieving && claims != nil {
		t.Error("should return ErrClaimsRetrieving")
	}
}

func jwks() *jose.JSONWebKeySet {
	// The JSON data
	jsonData := `{
		"alg": "RS256",
		"e": "AQAB",
		"key_ops": ["verify"],
		"kty": "RSA",
		"n": "n_jT5ubFYD2Vu6HT62RU21ZnimcRePKwdqb7K5SlS3k5MpijTlvArexwo9n3iuNxqymARyG1dMJEGwsGjFVHIL9uzSYq84n3IoFQHA6aE9R41f0pTMKb5Me-TAiHcR1lV8kx4YBRAPRUmREUnzOkKN-j8BlJzhuM7V7-_Swd41wWIyT0MYUg4-VxrO7Xw_ajBx6YUy2Oi-2-ljzc4ePuq2WCRIZ6ARXV2rCn4Gt8lUs_9jQ_ExmK5aGmiAfrJQRMWNH8JoRC1bgtI5w8Zs9Zhl2p5iZfa7OdDWMzz_Ic4X0rrNSmcvNTJFVRucCtTk3HVDgeBYDxzKDRoGqiBFLZWw",
		"use": "sig",
		"kid": "e13f8b15355a4c70941035152e2aaff8"
	}`

	// Create a struct to match the JSON structure
	type Key struct {
		Alg    string   `json:"alg"`
		E      string   `json:"e"`
		KeyOps []string `json:"key_ops"`
		Kty    string   `json:"kty"`
		N      string   `json:"n"`
		Use    string   `json:"use"`
		Kid    string   `json:"kid"`
	}

	// Parse the JSON data into the Key struct
	var keyData Key
	err := json.Unmarshal([]byte(jsonData), &keyData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	// Decode the base64-encoded "n" field into a big.Int
	nBytes, err := base64.RawURLEncoding.DecodeString(keyData.N)
	if err != nil {
		fmt.Println("Error decoding 'n' field:", err)
		return nil
	}
	n := new(big.Int).SetBytes(nBytes)

	// Decode the base64-encoded "e" field into an exponent value
	eBytes, err := base64.RawURLEncoding.DecodeString(keyData.E)
	if err != nil {
		fmt.Println("Error decoding 'e' field:", err)
		return nil
	}
	e := new(big.Int).SetBytes(eBytes).Int64()

	return &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &rsa.PublicKey{N: n, E: int(e)}}}}
}

/*
Token Generator
https://www.scottbrady91.com/tools/jwt

KEY
Private
{
  "alg": "RS256",
  "d": "Dv704kbMxtZPFH05j-3iVINXhm5eAXACocTKc83l5trQxVDwshZAzC0HbByxK1hh3fEwgLqEt5LEbqKMdRhDaCr52IpU6WqL-7SSjWbDA8vdnfWy6uqtUXd1-8uq4qwmRWHrZp-wOD4vNgAXZkshfuFkDUxZklQb1F6c2Z_Kl3by-vCWD9X4h-MweQ0kXFylrl4Qg_nVFI8XA39Ls6oJj1RsOVYxdUXWXtFS76IyxfmncnYet0IBHEDrkX_0VFWh3RVbVCn65k9zS4wc9Fog7xtV4Z8-JlyuqVb3tKjLtPvGVx3cVCAfM7XE9axSOFc8REZdvPXdwi5Jj56yu1nryQ",
  "dp": "aavSrXgI-eyF_Jv9rq6yugbTtLimp-kU3Y6Jht3VESfLvZvVE61NX8QOWBpnRLGmpKEJl4-QGq87vMpGXFcXFaMNA5_Zfce7BbnVNBe8X1g7mXZwY-xSMNI8dHJP0QQSeOcdN89wy8Rscy9QTSUE5JHdvXBcVTRzLqbUKLGdCmE",
  "dq": "kAAyanC825ix0rgvwMEswIm5c2IFD2NNmWBGxknhgcjSer1SLJZyZNRWvhec3p4Icr0HnlSPm-5dpEs37UZD1z9ETRfWxVid8vcjhbMpWQuNoXkxwmptIB3242906mKmpfri_NSSiHqK9tbEHJILbw8-4pyz29j7d4VVFsSl_VE",
  "e": "AQAB",
  "key_ops": [
    "sign"
  ],
  "kty": "RSA",
  "n": "n_jT5ubFYD2Vu6HT62RU21ZnimcRePKwdqb7K5SlS3k5MpijTlvArexwo9n3iuNxqymARyG1dMJEGwsGjFVHIL9uzSYq84n3IoFQHA6aE9R41f0pTMKb5Me-TAiHcR1lV8kx4YBRAPRUmREUnzOkKN-j8BlJzhuM7V7-_Swd41wWIyT0MYUg4-VxrO7Xw_ajBx6YUy2Oi-2-ljzc4ePuq2WCRIZ6ARXV2rCn4Gt8lUs_9jQ_ExmK5aGmiAfrJQRMWNH8JoRC1bgtI5w8Zs9Zhl2p5iZfa7OdDWMzz_Ic4X0rrNSmcvNTJFVRucCtTk3HVDgeBYDxzKDRoGqiBFLZWw",
  "p": "zv_ZgX1qxgKxWQV3NJ3x3Wzx4vTS0-xi-UrbzV7aCgPk3K0DbJz2bZXW1ZV7R5EqaJ0QOfqr-CjRckXPFC5tTfFUKP1oK1lw3Emp6o0qK39pDrUGsMNSX4di2gmVBalYvPtk_v64GOSjgGVl9OZOLvE0Jk5OpG3CzxrVZBGsePM",
  "q": "xdcefGp4eO4CA8R3BBKpnekxEV_XDJ0TkglhuBGk-tQMbUOSL-U0Ry0slqTQpVnzZzroAfRVI6Jp_rH1l2Iu-h_EpZKa_siqgcr9p2Z3OgSxr0H-tby-HzUin1Z1aSbe7x_v0AqO4ypXAvYAbP9d498lIu9HsNBv4vs_mGlXN_k",
  "qi": "RhbnVowqMSFeYkLqbZjHzJf363PS52Il-Z2cE-k-hqz1AhH3wlR9sA9upT_VwFacd2mnYEr0TbEJsKiyU-90IjmL3xdCmDcosPW-DLcZpK7UKoHKKGNe_E7bAXbx_NnySa5uUT-3yrnbYa36Qb8MtSUvkGuJEhfjIc96OwPuV10",
  "use": "sig",
  "kid": "e13f8b15355a4c70941035152e2aaff8"
}
Public
{
  "alg": "RS256",
  "e": "AQAB",
  "key_ops": [
    "verify"
  ],
  "kty": "RSA",
  "n": "n_jT5ubFYD2Vu6HT62RU21ZnimcRePKwdqb7K5SlS3k5MpijTlvArexwo9n3iuNxqymARyG1dMJEGwsGjFVHIL9uzSYq84n3IoFQHA6aE9R41f0pTMKb5Me-TAiHcR1lV8kx4YBRAPRUmREUnzOkKN-j8BlJzhuM7V7-_Swd41wWIyT0MYUg4-VxrO7Xw_ajBx6YUy2Oi-2-ljzc4ePuq2WCRIZ6ARXV2rCn4Gt8lUs_9jQ_ExmK5aGmiAfrJQRMWNH8JoRC1bgtI5w8Zs9Zhl2p5iZfa7OdDWMzz_Ic4X0rrNSmcvNTJFVRucCtTk3HVDgeBYDxzKDRoGqiBFLZWw",
  "use": "sig",
  "kid": "e13f8b15355a4c70941035152e2aaff8"
}
RS256/2048

TOKEN
Header
{
  "alg": "RS256",
  "typ": "JWT",
  "kid": "e13f8b15355a4c70941035152e2aaff8"
}
Payload
{
  "user_data": {
    "app_metadata": {},
    "created_at": "2022-12-07T08:01:20.986Z",
    "user_metadata": {}
  },
  "iss": "testIssuer",
  "aud": ["testAudience"],
  "sub": "5be86359073c434bad2da3932222dabe",
  "iat": 1687716817,
  "jti": "jti",
  "scope": "scope",
  "exp": 1687796692,
  "permissions": ["permissions"]
}
*/
