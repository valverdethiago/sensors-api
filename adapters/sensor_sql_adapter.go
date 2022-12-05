package adapters

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/valverdethiago/sensors-api/domain"
)

type SensorSqlAdapter struct {
	db *sql.DB
}

func NewSensorSqlAdapter(db *sql.DB) domain.SensorRepository {
	return SensorSqlAdapter{db: db}
}

func (s SensorSqlAdapter) Create(name string, location domain.Coordinate, tags []domain.Tag) (*domain.Sensor, error) {
	var id uuid.UUID
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	insertStmt, err := tx.Prepare(`INSERT INTO sensors (name, location) 
													 VALUES   ($1, POINT($2, $3 ) ) 
										  RETURNING sensor_uuid`)
	if err != nil {
		return nil, err
	}
	defer insertStmt.Close()
	err = insertStmt.QueryRow(name, location.Longitude, location.Latitude).Scan(&id)
	if err != nil {
		return nil, err
	}
	dbTags, err := s.persistTags(tx, id, tags)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()

	if err != nil {
		return nil, err
	}
	return &domain.Sensor{
		Id:       &id,
		Name:     name,
		Location: location,
		Tags:     dbTags,
	}, nil
}

func (s SensorSqlAdapter) Update(name string, newName string, location domain.Coordinate, tags []domain.Tag) (*domain.Sensor, error) {
	var id uuid.UUID
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(
		`  UPDATE sensors SET name = $2, location = POINT($3, $4) 
                   WHERE name = $1 
               RETURNING sensor_uuid`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(name, newName, location.Longitude, location.Latitude).Scan(&id)
	if err != nil {
		return nil, err
	}
	err = s.dropOldTags(tx, id)
	if err != nil {
		return nil, err
	}
	dbTags, err := s.persistTags(tx, id, tags)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return &domain.Sensor{
		Id:       &id,
		Name:     name,
		Location: location,
		Tags:     dbTags,
	}, nil
}

func (s SensorSqlAdapter) GetByName(term string) (sensor *domain.Sensor, err error) {
	var id uuid.UUID
	var name string
	var longitude float64
	var latitude float64
	selectStmt, err := s.db.Prepare(`SELECT sensor_uuid, name, location[0], location[1] FROM sensors WHERE name = $1 `)
	if err != nil {
		return nil, err
	}
	defer selectStmt.Close()
	err = selectStmt.QueryRow(term).Scan(&id, &name, &longitude, &latitude)
	if err != nil {
		return nil, err
	}
	tags, err := s.findTags(id)
	if err != nil {
		return nil, err
	}
	return &domain.Sensor{
		Id:   &id,
		Name: name,
		Location: domain.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		},
		Tags: tags,
	}, nil
}

func (s SensorSqlAdapter) GetById(id uuid.UUID) (*domain.Sensor, error) {
	var dbId uuid.UUID
	var name string
	var longitude float64
	var latitude float64
	selectStmt, err := s.db.Prepare("SELECT sensor_uuid, name, location[0], location[1] FROM sensors WHERE sensor_uuid = $1")
	if err != nil {
		return nil, err
	}
	defer selectStmt.Close()
	err = selectStmt.QueryRow(id.String()).Scan(&dbId, &name, &longitude, &latitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	tags, err := s.findTags(id)
	if err != nil {
		return nil, err
	}
	return &domain.Sensor{
		Id:   &dbId,
		Name: name,
		Location: domain.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		},
		Tags: tags,
	}, nil
}

func (s SensorSqlAdapter) FindNearestSensor(location domain.Coordinate) (*domain.Sensor, float64, error) {
	distanceQuery := `
		SELECT sensor_uuid, name, location[0], location[1], 
				 (location <-> point($1, $2)) as distance
			FROM sensors
		   WHERE location != point($1, $2)
		ORDER BY distance
		   LIMIT 1`
	var distance float64
	selectStmt, err := s.db.Prepare(distanceQuery)
	if err != nil {
		return nil, distance, err
	}
	sensor := domain.Sensor{
		Location: domain.Coordinate{},
	}
	defer selectStmt.Close()
	err = selectStmt.QueryRow(location.Longitude, location.Latitude).Scan(&sensor.Id,
		&sensor.Name, &sensor.Location.Longitude, &sensor.Location.Latitude, &distance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, distance, nil
		}
		return nil, distance, err
	}
	tags, err := s.findTags(*sensor.Id)
	if err != nil {
		return nil, distance, err
	}
	sensor.Tags = tags
	return &sensor, distance, nil
}

func (s SensorSqlAdapter) GetAll() ([]domain.Sensor, error) {
	result := make([]domain.Sensor, 0)
	rows, err := s.db.Query("SELECT sensor_uuid, name, location[0], location[1] FROM sensors limit 100")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var sensor domain.Sensor
		var coordinate domain.Coordinate
		err := rows.Scan(&sensor.Id, &sensor.Name, &coordinate.Longitude, &coordinate.Latitude)
		if err != nil {
			return nil, err
		}
		sensor.Location = coordinate
		result = append(result, sensor)
	}
	defer rows.Close()
	return result, nil
}

func (s SensorSqlAdapter) persistTags(tx *sql.Tx, sensorUuid uuid.UUID, tags []domain.Tag) (*[]domain.Tag, error) {
	dbTags := make([]domain.Tag, len(tags))
	for i, tag := range tags {
		dbTag, err := s.createTag(tx, sensorUuid, tag.Name, tag.Value)
		if err != nil {
			return nil, err
		}
		dbTags[i] = *dbTag
	}
	return &dbTags, nil
}

func (s SensorSqlAdapter) findTags(sensorUuid uuid.UUID) (*[]domain.Tag, error) {
	dbTags := make([]domain.Tag, 0)
	rows, err := s.db.Query("SELECT tag_uuid, name, value FROM tags WHERE sensor_uuid = $1", sensorUuid.String())
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.Id, &tag.Name, &tag.Value)
		if err != nil {
			return nil, err
		}
		dbTags = append(dbTags, tag)
	}
	return &dbTags, nil
}

func (s SensorSqlAdapter) createTag(tx *sql.Tx, sensorUuid uuid.UUID, name string, value string) (*domain.Tag, error) {
	var id uuid.UUID
	insertStmt, err := tx.Prepare(`INSERT INTO tags (name, value, sensor_uuid) 
													 VALUES   ($1, $2, $3 ) 
										  RETURNING tag_uuid`)
	if err != nil {
		return nil, err
	}
	defer insertStmt.Close()
	err = insertStmt.QueryRow(name, value, sensorUuid.String()).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &domain.Tag{
		Id:    &id,
		Name:  name,
		Value: value,
	}, nil
}

func (s SensorSqlAdapter) dropOldTags(tx *sql.Tx, id uuid.UUID) error {
	stmt, err := tx.Prepare(`DELETE FROM tags WHERE sensor_uuid = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id.String())
	return err
}
