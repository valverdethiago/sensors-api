package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valverdethiago/sensors-api/domain"
	"net/http"
	"strconv"
	"strings"
)

const (
	invalidCoordinatesErrorMsg     = "invalid coordinates input"
	idIsRequiredErrorMsg           = "sensor id is required"
	idIsNotAllowedInSensorCreation = "id is not allowed in sensor creation"
)

type SensorController struct {
	service domain.SensorService
}

func NewSensorController(service domain.SensorService) SensorController {
	controller := SensorController{
		service: service,
	}
	return controller
}

func (controller *SensorController) Search(context *gin.Context) {
	var coordinate *domain.Coordinate
	var result *domain.Sensor
	var err error
	if coordinate, err = extractSearchParameters(context); err != nil {
		return
	}
	if result, err = controller.service.FindNearestSensor(*coordinate); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result != nil {
		context.JSON(http.StatusOK, result)
	}
	context.Status(http.StatusNotFound)

}

func extractSearchParameters(context *gin.Context) (*domain.Coordinate, error) {
	latitude, _ := context.GetQuery("lat")
	longitude, _ := context.GetQuery("long")
	lat, errLat := parse(&latitude)
	lon, errLong := parse(&longitude)
	if errLat != nil || errLong != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": invalidCoordinatesErrorMsg})
		return nil, errors.New(invalidCoordinatesErrorMsg)
	}

	return &domain.Coordinate{
		Latitude:  lat,
		Longitude: lon,
	}, nil

}

func (controller *SensorController) GetById(context *gin.Context) {
	var idHash uuid.UUID
	var err error
	var result *domain.Sensor
	idHash, err = controller.extractIdFromContext(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if result, err = controller.service.GetById(idHash); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if result == nil {
		context.Status(http.StatusNotFound)
	} else {
		context.JSON(http.StatusOK, result)
	}
}

func (controller *SensorController) Create(context *gin.Context) {

	var err error
	var sensor *domain.Sensor

	if err = json.NewDecoder(context.Request.Body).Decode(&sensor); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sensor.SensorUuid != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": idIsNotAllowedInSensorCreation})
		return
	}
	tags := make([]domain.Tag, 0)
	if sensor, err = controller.service.Create(sensor.Name, sensor.Location, tags); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, sensor)
}

func (controller *SensorController) Update(context *gin.Context) {

}

func (controller *SensorController) GetTags(context *gin.Context) {

}

func (controller *SensorController) extractIdFromContext(context *gin.Context) (idHash uuid.UUID, err error) {
	id := context.Param("id")
	if strings.Trim(id, " ") == "" {
		return idHash, errors.New(idIsRequiredErrorMsg)
	}
	if idHash, err = uuid.Parse(id); err != nil {
		return idHash, errors.New(idIsRequiredErrorMsg)
	}
	return idHash, nil
}

func parse(number *string) (float64, error) {
	if s, err := strconv.ParseFloat(*number, 64); err != nil {
		return s, err
	}
	return 0, nil
}
