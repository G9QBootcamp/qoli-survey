package db

import (
	"fmt"

	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDb struct {
	dbClient *gorm.DB
}

func (p *postgresDb) Init(cfg *config.Config) error {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)
	var err error

	p.dbClient, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDb, _ := p.dbClient.DB()
	err = sqlDb.Ping()
	if err != nil {
		return err
	}
	sqlDb.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime * time.Minute)

	return nil
}
func (p *postgresDb) Close() {
	conn, _ := p.dbClient.DB()
	conn.Close()
}

func (p *postgresDb) GetDb() *gorm.DB {
	return p.dbClient
}

func newPostgresDb() *postgresDb {
	return &postgresDb{}
}
