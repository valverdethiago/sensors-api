CREATE TABLE  sensors
(
    sensor_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    location point NOT NULL,
    PRIMARY KEY(sensor_uuid),
    UNIQUE(name)
);