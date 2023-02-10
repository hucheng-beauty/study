package config

type ServerConfig struct {
	Name         string         `yaml:"name"`
	Host         string         `yaml:"host"`
	Port         string         `yaml:"port"`
	PostgresInfo PostgresConfig `yaml:"postgres_info"`
	JWTInfo      JWTConfig      `yaml:"jwt"`
}

type PostgresConfig struct {
	Endpoint string `yaml:"endpoint"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type JWTConfig struct {
	SigningKey string `yaml:"key" json:"key"`
}
