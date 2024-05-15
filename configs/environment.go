package configs

type Environment int8

const (
	UnknownEnv     Environment = 0
	LocalEnv       Environment = 1
	DevelopmentEnv Environment = 2
	StagingEnv     Environment = 3
	QaEnv          Environment = 4
	ProductionEnv  Environment = 5
)

var (
	EnvironmentMapping = map[Environment]string{
		UnknownEnv:     "unknown",
		LocalEnv:       "local",
		DevelopmentEnv: "development",
		StagingEnv:     "staging",
		QaEnv:          "qa",
		ProductionEnv:  "production",
	}
)

func NewEnvironment(env string) Environment {
	switch env {
	case "local":
		fallthrough
	case "LOCAL":
		return LocalEnv
	case "development":
		fallthrough
	case "DEVELOPMENT":
		fallthrough
	case "dev":
		fallthrough
	case "DEV":
		return DevelopmentEnv
	case "production":
		fallthrough
	case "PRODUCTION":
		fallthrough
	case "prod":
		fallthrough
	case "PROD":
		return ProductionEnv
	case "staging":
		fallthrough
	case "STAGING":
		fallthrough
	case "stg":
		fallthrough
	case "STG":
		return StagingEnv
	case "qa":
		fallthrough
	case "QA":
		return QaEnv
	default:
		return UnknownEnv
	}
}

func (e Environment) ToString() string {
	switch e {
	case LocalEnv:
		return "local"
	case DevelopmentEnv:
		return "development"
	case ProductionEnv:
		return "production"
	case StagingEnv:
		return "staging"
	case QaEnv:
		return "qa"
	default:
		return "unknown"
	}
}
