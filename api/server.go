package api

import (
	"flag"
	"github.com/gin-gonic/gin"
	ginglog "github.com/szuecs/gin-glog"
	"github.com/valverdethiago/sensors-api/config"
	"net/http"
	"time"
)

const (
	sensorBasePath = "/sensor"
	sensorIdPath   = "/sensor/:id"
	sensorTagsPath = "/sensor/:id/tags"
	nearestSensor  = "/nearest"
)

type Server struct {
	router     *gin.Engine
	config     *config.ServerConfig
	controller *SensorController
}

func NewServer(router *gin.Engine, config *config.ServerConfig, controller *SensorController) *Server {
	return &Server{
		router:     router,
		config:     config,
		controller: controller,
	}
}

func (server *Server) Configure() {
	server.configureLogging()
	server.configureRoutes()
}

func (server *Server) configureLogging() {
	flag.Parse()
	server.router.Use(ginglog.Logger(3 * time.Second))
	server.router.Use(gin.Recovery())
}

func (server *Server) configureRoutes() {
	server.router.GET(sensorBasePath, server.controller.Search)
	server.router.GET(sensorIdPath, server.controller.GetById)
	server.router.PUT(sensorBasePath, server.controller.Create)
	server.router.POST(sensorBasePath, server.controller.Update)
	server.router.GET(nearestSensor, server.controller.NearestSensor)
	server.router.GET(sensorTagsPath, server.controller.GetTags)
}

// Start runs the HTTP Server on a specific address
func (server *Server) Start() error {
	var readTimeout time.Duration
	readTimeout = time.Duration(server.config.ReadTimeout) * time.Second
	var writeTimeout time.Duration
	writeTimeout = time.Duration(server.config.WriteTimeout) * time.Second

	s := &http.Server{
		Addr:         server.config.ServerAddress,
		Handler:      server.router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	return s.ListenAndServe()
}
