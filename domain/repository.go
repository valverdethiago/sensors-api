package domain

import "github.com/google/uuid"

type SensorRepository interface {
	Create(name string, location Coordinate) (*Sensor, error)
	CreateTag(sensor *Sensor, name string, value string) (*Tag, error)
	Update(name string, location Coordinate) (*Sensor, error)
	GetByName(name string) (*Sensor, error)
	GetById(uuid uuid.UUID) (*Sensor, error)
	FindNearestSensor(location Coordinate) (*Sensor, error)
	ClearTags(sensor *Sensor) error
}
