package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"study/internal/admin_api/global"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func DB() {
	postgresInfo := global.ServerConfig.PostgresInfo
	global.DB = connect(dsn(postgresInfo.Endpoint, postgresInfo.Username,
		postgresInfo.Password, postgresInfo.Database, postgresInfo.Port))
}

func dsn(endpoint, username, password, database string, port int) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		endpoint, username, password, database, port)
}

func connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，不自动给表名加s
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second, // 慢 sql 阈值
				Colorful:      true,        // 禁用彩色打印
				LogLevel:      logger.Info,
			},
		),
	})
	if err != nil {
		zap.S().Panic("failed to connect database", zap.String("dsn", dsn), zap.Error(err))
	}

	return db
}
