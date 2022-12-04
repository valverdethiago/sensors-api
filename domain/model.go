package domain

type Point struct {
	Latitude  float64
	Longitude float64
}

type Sensor struct {
	Name     string
	Location Point
	Tags     []Tag
}

type Tag struct {
	Name  string
	Value string
}
