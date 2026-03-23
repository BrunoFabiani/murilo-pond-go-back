package middlewere

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Telemetry struct {
	DeviceID      string    `json:"device_id"`
	Timestamp     time.Time `json:"timestamp"`
	SensorType    string    `json:"sensor_type"`
	ReadingNature string    `json:"reading_nature"`
	Value         any       `json:"value"`
}

func OpenConn() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/telemetria?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func validateTelemetry(t *Telemetry) error {
	t.DeviceID = strings.TrimSpace(t.DeviceID)
	t.SensorType = strings.TrimSpace(t.SensorType)
	t.ReadingNature = strings.ToLower(strings.TrimSpace(t.ReadingNature))

	if t.DeviceID == "" {
		return errors.New("device_id is required")
	}
	if t.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if t.SensorType == "" {
		return errors.New("sensor_type is required")
	}

	switch t.ReadingNature {
	case "analog":
		if _, ok := t.Value.(float64); !ok {
			return errors.New("analog value must be numeric")
		}
	case "discrete":
		switch t.Value.(type) {
		case bool, string:
		default:
			return errors.New("discrete value must be bool or string")
		}
	default:
		return errors.New("reading_nature must be analog or discrete")
	}
	return nil
}

func insertTelemetry(db *sql.DB, msg Telemetry) error {
	valueJSON, err := json.Marshal(msg.Value)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO telemetry (device_id, timestamp, sensor_type, reading_nature, value)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = db.Exec(query, msg.DeviceID, msg.Timestamp, msg.SensorType, msg.ReadingNature, valueJSON)
	return err
}

func ProcessDelivery(db *sql.DB, d amqp.Delivery) error {
	var msg Telemetry
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		return err
	}

	if err := validateTelemetry(&msg); err != nil {
		return err
	}

	return insertTelemetry(db, msg)
}
