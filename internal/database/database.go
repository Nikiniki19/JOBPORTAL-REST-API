package database

import (
	"fmt"
	"job-portal-api/config"
	"job-portal-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Implementing openconn func to connect to the db
func OpenConnection() (*gorm.DB, error) {
	//fmt.Println("hi")
	cfg := config.GetConfig()
	fmt.Println(cfg.PostgresConfig)
	dsn := fmt.Sprintf("host=%s user=%s password=%s  dbname=%s  port=%s  sslmode=%s TimeZone=%s", cfg.PostgresConfig.Host, cfg.PostgresConfig.User, cfg.PostgresConfig.Password, cfg.PostgresConfig.Db, cfg.PostgresConfig.DbPort, cfg.PostgresConfig.SslMode, cfg.PostgresConfig.TimeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//fmt.Println("DB:====", db.Migrator().AutoMigrate(&models.User{}))
	err = db.Migrator().AutoMigrate(&models.User{})
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return nil, err
	}
	err = db.Migrator().AutoMigrate(&models.Company{})
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return nil, err
	}
	err = db.Migrator().AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Job{},
		&models.Location{},
		&models.Skill{},
		&models.WorkMode{},
		&models.Qualification{},
		&models.Shift{},
		&models.JobType{},
	)
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return nil, err
	}
	return db, nil
}
