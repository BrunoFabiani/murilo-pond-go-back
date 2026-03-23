CREATE TABLE telemetry (
    id BIGSERIAL PRIMARY KEY,
    device_id TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    sensor_type TEXT NOT NULL,
    reading_nature TEXT NOT NULL CHECK (reading_nature IN ('discrete', 'analog')),
    value JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

