package adapters

import (
	"github.com/google/uuid"
	"github.com/valverdethiago/sensors-api/domain"
	"gorm.io/gorm"
)

type SensorGormRepository struct {
	db *gorm.DB
}

func NewSensorGormRepository(db *gorm.DB) domain.SensorRepository {
	return SensorGormRepository{db: db}
}

func (s SensorGormRepository) Create(name string, location domain.Coordinate) (*domain.Sensor, error) {
	sensor := domain.Sensor{
		Name:     name,
		Location: location,
		//Tags:     tags,
	}
	result := s.db.Create(&sensor)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sensor, nil
}

func (s SensorGormRepository) CreateTag(sensor *domain.Sensor, name string, value string) (*domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (s SensorGormRepository) Update(name string, location domain.Coordinate) (*domain.Sensor, error) {
	//TODO implement me
	panic("implement me")
}

func (s SensorGormRepository) GetByName(name string) (*domain.Sensor, error) {
	//TODO implement me
	panic("implement me")
}

func (s SensorGormRepository) GetById(uuid uuid.UUID) (*domain.Sensor, error) {
	//TODO implement me
	panic("implement me")
}

func (s SensorGormRepository) FindNearestSensor(location domain.Coordinate) (*domain.Sensor, error) {
	//TODO implement me
	panic("implement me")
}

func (s SensorGormRepository) ClearTags(sensor *domain.Sensor) error {
	//TODO implement me
	panic("implement me")
}
