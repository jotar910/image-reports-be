package configs

type GlobalConfig struct {
	Gin      ginConfig
	Services serviceListConfig
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
