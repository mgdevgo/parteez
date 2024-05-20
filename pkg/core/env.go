package core

type AppEnv string

const (
	AppEnvLocal AppEnv = "local"
	AppEnvDev   AppEnv = "dev"
	AppEnvProd  AppEnv = "prod"
)
