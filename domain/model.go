package domain

import (
	"database/sql/driver"
	"errors"
	"github.com/google/uuid"
	"strconv"
)

type Sensor struct {
	SensorUuid *uuid.UUID `json:"sensor_uuid,omitempty" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Name       string     `gorm:"colum:name" json:"name" `
	Location   Coordinate `gorm:"colum:location" json:"location" `
	Tags       *[]Tag     `gorm:"foreignKey:sensor_uuid;references:sensor_uuid" json:"tags" `
}

type Tag struct {
	Id       *uuid.UUID `gorm:"colum:tag_uuid;primarykey;type:uuid;default:uuid_generate_v4()" json:"id,omitempty" `
	Name     string     `gorm:"colum:name" json:"name" `
	Value    string     `gorm:"colum:value" json:"value" `
	SensorId uuid.UUID  `gorm:"colum:sensor_uuid;type:uuid" json:"-" `
	Sensor   Sensor     `gorm:"foreignKey:sensor_uuid;references:sensor_uuid" json:"-" `
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

func (p *Coordinate) Value() (driver.Value, error) {
	out := []byte("POINT(")
	out = strconv.AppendFloat(out, p.Longitude, 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, p.Latitude, 'f', -1, 64)
	out = append(out, ')')
	return out, nil
}

func (p *Coordinate) Scan(src interface{}) (err error) {
	var data []byte
	switch src := src.(type) {
	case []byte:
		data = src
	case string:
		data = []byte(src)
	case nil:
		return nil
	default:
		return errors.New("(*Coordinate).Scan: unsupported data type")
	}

	if len(data) == 0 {
		return nil
	}

	data = data[1 : len(data)-1] // drop the surrounding parentheses
	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			if p.Latitude, err = strconv.ParseFloat(string(data[:i]), 64); err != nil {
				return err
			}
			if p.Longitude, err = strconv.ParseFloat(string(data[i+1:]), 64); err != nil {
				return err
			}
			break
		}
	}
	return nil
}
