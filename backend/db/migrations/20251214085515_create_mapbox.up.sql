CREATE TABLE IF NOT EXISTS map_locs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    address TEXT,
    description TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    location_type VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE rooms
ADD CONSTRAINT fk_rooms_map_locs
FOREIGN KEY (map_loc_id)
REFERENCES map_locs(id)
ON DELETE SET NULL;

CREATE INDEX idx_map_locs_lat_lng ON map_locs (latitude, longitude);
