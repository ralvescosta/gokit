package env

type Environment int8

func NewEnvironment(env string) Environment {
	switch env {
	case "local":
		fallthrough
	case "LOCAL":
		return LOCAL_ENV
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
		return UNKNOWN_ENV
	}
}

func (e Environment) ToString() string {
	switch e {
	case LOCAL_ENV:
		return "local"
	case DEVELOPMENT_ENV:
		return "development"
	case PRODUCTION_ENV:
		return "production"
	case STAGING_ENV:
		return "staging"
	case QA_ENV:
		return "qa"
	default:
		return "unknown"
	}
}
