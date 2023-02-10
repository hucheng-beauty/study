package initialize

import (
	"flag"
)

var (
	cfg = flag.String("config", defaultConfigPath /*TODO 加载正确的配置文件路径*/, "config file")
)

func init() {
	Logger()
	flag.Parse()
	Config(*cfg)
	DB()
	Cache()
}
