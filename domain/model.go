package domain

import (
	"github.com/google/uuid"
)

type Sensor struct {
	Id       *uuid.UUID `json:"id,omitempty"`
	Name     string     `json:"name" `
	Location Coordinate `json:"location" `
	Tags     *[]Tag     `json:"tags,omitempty" `
}

type NearestResponse struct {
	Sensor   Sensor  `json:"sensor" `
	Distance float64 `json:"distance" `
}

type Tag struct {
	Id    *uuid.UUID `json:"id,omitempty" `
	Name  string     `json:"name" `
	Value string     `json:"value" `
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}
