-- +goose Up
CREATE TABLE store_owner (
    id UUID PRIMARY KEY,

    user_id UUID,
    store_id UUID,

    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,

    tile_x INTEGER NOT NULL,
    tile_y INTEGER NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CHECK (lat BETWEEN -90 AND 90),
    CHECK (lng BETWEEN -180 AND 180)
);

CREATE INDEX idx_store_owner_tile_xy_z15
ON store_owner (tile_x, tile_y);

CREATE INDEX idx_store_owner_lat_lng
ON store_owner (lat, lng);


-- +goose Down
DROP TABLE IF EXISTS store_owner;
