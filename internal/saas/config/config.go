package config

import "gorm.io/gorm"

var (
    ServerInfo = &ServerConfig{}
    Db         *gorm.DB // todo
)

type ServerConfig struct {
    ApiServiceInfo *ApiServiceConfig `yaml:"api_service_info"`
}

type ApiServiceConfig struct {
    ServiceName string `yaml:"service_name"`
    ServiceHost string `yaml:"service_host"`
    ServicePort int    `yaml:"service_port"`
}
