package env

type DeploymentType string

const (
	DEV  DeploymentType = "dev"
	PROD DeploymentType = "prod"
	MOCK DeploymentType = "mock"
)
