CREATE TABLE tags
(
    tag_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    sensor_uuid UUID NOT NULL,
    PRIMARY KEY(tag_uuid),
    FOREIGN KEY (sensor_uuid) REFERENCES sensors (sensor_uuid)
);