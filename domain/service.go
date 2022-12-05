package domain

import "github.com/google/uuid"

type SensorService interface {
	Create(name string, location Coordinate, tags []Tag) (*Sensor, error)
	Update(name string, newName string, location Coordinate, tags []Tag) (*Sensor, error)
	GetByName(name string) (*Sensor, error)
	GetById(uuid uuid.UUID) (*Sensor, error)
	FindNearestSensor(coordinate Coordinate) (*Sensor, float64, error)
	SearchByName(name string) (*Sensor, error)
	GetAll() ([]Sensor, error)
}

type SensorServiceImpl struct {
	repository SensorRepository
}

func NewSensorServiceImpl(repository SensorRepository) SensorService {
	return &SensorServiceImpl{
		repository: repository,
	}
}

func (s SensorServiceImpl) Create(name string, location Coordinate, tags []Tag) (sensor *Sensor, err error) {
	if sensor, err = s.repository.Create(name, location, tags); err != nil {
		return nil, err
	}
	return sensor, nil
}

func (s SensorServiceImpl) Update(name string, newName string, location Coordinate, tags []Tag) (sensor *Sensor, err error) {
	return s.repository.Update(name, newName, location, tags)
}

func (s SensorServiceImpl) GetByName(name string) (*Sensor, error) {
	return s.repository.GetByName(name)
}

func (s SensorServiceImpl) GetById(uuid uuid.UUID) (*Sensor, error) {
	return s.repository.GetById(uuid)
}

func (s SensorServiceImpl) FindNearestSensor(coordinate Coordinate) (*Sensor, float64, error) {
	return s.repository.FindNearestSensor(coordinate)
}

func (s SensorServiceImpl) SearchByName(name string) (*Sensor, error) {
	return s.repository.GetByName(name)
}

func (s SensorServiceImpl) GetAll() ([]Sensor, error) {
	return s.repository.GetAll()
}
