package env

type Environment int8

func NewEnvironment(env string) Environment {
	switch env {
	case "development":
		fallthrough
	case "DEVELOPMENT":
		fallthrough
	case "dev":
		fallthrough
	case "DEV":
		return DEVELOPMENT_ENV
	case "production":
		fallthrough
	case "PRODUCTION":
		fallthrough
	case "prod":
		fallthrough
	case "PROD":
		return PRODUCTION_ENV
	case "staging":
		fallthrough
	case "STAGING":
		fallthrough
	case "stg":
		fallthrough
	case "STG":
		return STAGING_ENV
	case "qa":
		fallthrough
	case "QA":
		return QA_ENV
	default:
		return DEVELOPMENT_ENV
	}
}
