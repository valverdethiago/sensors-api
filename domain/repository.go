package domain

import "github.com/google/uuid"

type SensorRepository interface {
	Create(name string, location Point, tags *[]Tag) (*Sensor, error)
	Update(name string, location Point, tags *[]Tag) (*Sensor, error)
	GetByName(name string) (*Sensor, error)
	GetById(uuid uuid.UUID) (*Sensor, error)
	FindNearestSensor(location Point) (*Sensor, error)
}
