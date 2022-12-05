package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/valverdethiago/sensors-api/adapters"
	"github.com/valverdethiago/sensors-api/api"
	"github.com/valverdethiago/sensors-api/config"
	"github.com/valverdethiago/sensors-api/domain"
	"log"
)

const DatabaseDsnLayout string = "host=%s user=%s dbname=%s port=%s sslmode=disable"

func main() {
	cfg := loadAppConfig()
	db := connectToDatabase(cfg.DatabaseConfig)
	repository := configureRepository(db)
	service := configureService(repository)
	controller := configureController(service)
	server := configureServer(cfg.ServerConfig, controller)
	server.Start()
	fmt.Println(db)
}

func configureController(service domain.SensorService) api.SensorController {
	return api.NewSensorController(service)
}

func configureService(repository domain.SensorRepository) domain.SensorService {
	return domain.NewSensorServiceImpl(repository)
}

func configureRepository(db *sql.DB) domain.SensorRepository {
	return adapters.NewSensorSqlAdapter(db)
}

func configureServer(config config.ServerConfig, controller api.SensorController) *api.Server {
	router := gin.Default()
	server := api.NewServer(router, &config, &controller)
	server.Configure()
	return server
}

func loadAppConfig() *config.AppConfig {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config file:", err)
	}
	return cfg
}
func connectToDatabase(cfg config.DatabaseConfig) *sql.DB {
	dsn := fmt.Sprintf(DatabaseDsnLayout, cfg.Host, cfg.Username, cfg.Name, cfg.Port)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}
	return db
}
