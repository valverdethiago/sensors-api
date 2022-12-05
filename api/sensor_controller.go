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
	name, hasName := context.GetQuery("name")
	if hasName {
		sensor, err := controller.service.SearchByName(name)
		handleResult(context, sensor, err)
	} else {
		result, err := controller.service.GetAll()
		handleResult(context, result, err)
	}

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
	if sensor.Id != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": idIsNotAllowedInSensorCreation})
		return
	}
	if sensor, err = controller.service.Create(sensor.Name, sensor.Location, *sensor.Tags); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, sensor)
}

func (controller *SensorController) NearestSensor(context *gin.Context) {
	coordinate, err := extractCoordinatesFromQuery(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sensor, distance, err := controller.service.FindNearestSensor(*coordinate)

	handleResult(context, domain.NearestResponse{
		Sensor:   *sensor,
		Distance: distance,
	}, err)
}

func (controller *SensorController) Update(context *gin.Context) {
	name, hasName := context.GetQuery("name")
	if !hasName {
		context.Status(http.StatusBadRequest)
	}
	var sensor *domain.Sensor

	if err := json.NewDecoder(context.Request.Body).Decode(&sensor); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := controller.service.Update(name, sensor.Name, sensor.Location, *sensor.Tags)
	handleResult(context, result, err)
}

func (controller *SensorController) GetTags(context *gin.Context) {
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
		context.JSON(http.StatusOK, result.Tags)
	}
}

func handleResult(context *gin.Context, result interface{}, err error) {
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil {
		context.Status(http.StatusNotFound)
		return
	}
	context.JSON(http.StatusOK, result)
	return
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

func extractCoordinatesFromQuery(context *gin.Context) (*domain.Coordinate, error) {
	latitude, _ := context.GetQuery("lat")
	longitude, _ := context.GetQuery("lon")
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

func parse(number *string) (float64, error) {
	if s, err := strconv.ParseFloat(*number, 64); err != nil {
		return s, err
	}
	return 0, nil
}
