package mysqlCli

import (
	"FantasticLife/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 定义一个模型
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func NewMysqlCli(logger *zap.Logger, config *config.Config) (mysqlClient *gorm.DB) {
	mysqlCofig := config.Mysql
	logger.Info("初始化mysql", zap.Any("mysqlCofig", mysqlCofig))
	dsn := mysqlCofig.Username + ":" + mysqlCofig.Password + "@tcp(" + mysqlCofig.Addr + ")/" + mysqlCofig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("初始化mysql失败", zap.Error(err), zap.String("dsn", dsn))
		panic("failed to connect database")
	}
	logger.Info("初始化DB成功 connect database success:", zap.String("dsn", dsn))
	return db
}
