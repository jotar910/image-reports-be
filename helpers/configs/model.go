package configs

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type GlobalConfig struct {
	Gin      ginConfig
	Services serviceListConfig
	Image    imageConfig
}

type ginConfig struct {
	Mode string
}

type serviceListConfig struct {
	ApiGateway serviceConfig
	Users      serviceConfig
	Reporter   serviceConfig
	Storage    serviceConfig
	Processing serviceConfig
	Realtime   serviceConfig
}

type serviceConfig struct {
	Host string
	Port int
}

type imageConfig struct {
	MaxSize    int64
	Extensions string
}
