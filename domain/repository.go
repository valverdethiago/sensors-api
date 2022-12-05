package domain

import "github.com/google/uuid"

type SensorRepository interface {
	Create(name string, location Coordinate, tags []Tag) (*Sensor, error)
	Update(name string, newName string, location Coordinate, tags []Tag) (*Sensor, error)
	GetByName(name string) (*Sensor, error)
	GetById(uuid uuid.UUID) (*Sensor, error)
	FindNearestSensor(location Coordinate) (*Sensor, float64, error)
	GetAll() ([]Sensor, error)
}
