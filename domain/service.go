package domain

import "github.com/google/uuid"

type SensorService interface {
	Create(name string, location Coordinate, tags []Tag) (*Sensor, error)
	Update(name string, location Coordinate, tags []Tag) (*Sensor, error)
	GetByName(name string) (*Sensor, error)
	GetById(uuid uuid.UUID) (*Sensor, error)
	FindNearestSensor(coordinate Coordinate) (*Sensor, error)
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

func (s SensorServiceImpl) Update(name string, location Coordinate, tags []Tag) (sensor *Sensor, err error) {
	if sensor, err = s.repository.Update(name, location); err != nil {
		return nil, err
	}
	if err = s.repository.ClearTags(sensor); err != nil {
		return nil, err
	}
	if err = s.persistTags(sensor, tags); err != nil {
		return nil, err
	}
	return sensor, nil

}

func (s SensorServiceImpl) persistTags(sensor *Sensor, tags []Tag) error {
	result := make([]Tag, len(tags))
	for i, tag := range tags {
		dbTag, err := s.repository.CreateTag(sensor, tag.Name, tag.Value)
		if err != nil {
			return err
		}
		result[i] = *dbTag
	}
	//sensor.Tags = result
	return nil
}

func (s SensorServiceImpl) GetByName(name string) (*Sensor, error) {
	return s.repository.GetByName(name)
}

func (s SensorServiceImpl) GetById(uuid uuid.UUID) (*Sensor, error) {
	return s.repository.GetById(uuid)
}

func (s SensorServiceImpl) FindNearestSensor(coordinate Coordinate) (*Sensor, error) {
	return s.repository.FindNearestSensor(coordinate)
}
