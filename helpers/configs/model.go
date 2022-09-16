package configs

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type GlobalConfig struct {
	Gin      GinConfig
	Services serviceListConfig
	Image    ImageConfig
}

type GinConfig struct {
	Mode string
}

type serviceListConfig struct {
	ApiGateway ServiceConfig
	Users      ServiceConfig
	Reporter   ServiceConfig
	Storage    ServiceConfig
	Processing ServiceConfig
	Realtime   ServiceConfig
}

type ServiceConfig struct {
	Host string
	Port int
}

type ImageConfig struct {
	MaxSize    int64
	Extensions string
}
