package data

import (
	"AltTreasury/internal/conf"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewTreasuryRepo)

// Data .
type Data struct {
	db  *gorm.DB
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)
	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		mysqlDSN = c.Database.Source
	}
	db, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		log.Errorf("failed opening connection to mysql: %v", err)
		return nil, nil, err
	}

	// 自动迁移数据库结构
	// if err := db.AutoMigrate(&WithdrawClaim{}); err != nil {
	// 	log.Errorf("failed auto migrating database: %v", err)
	// 	return nil, nil, err
	// }

	// 设置自定义表名
	db.Table("withdrawal_claim").AutoMigrate(&WithdrawClaim{})
	db.Table("withdrawal_claim_confirmation").AutoMigrate(&WithdrawClaimConfirmation{})

	d := &Data{
		db:  db,
		log: log,
	}
	return d, func() {
		log.Info("closing the data resources")
		sqlDB, err := d.db.DB()
		if err != nil {
			log.Error(err)
		}
		sqlDB.Close()
	}, nil
}
