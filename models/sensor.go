package models

import "time"

type SensorData struct {
	SensorValue int       `json:"sensor_value"`
	ID1         int       `json:"id1"`
	ID2         string    `json:"id2"`
	Timestamp   time.Time `json:"timestamp"`
}
