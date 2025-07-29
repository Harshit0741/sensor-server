package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sensor-server/models"
	"time"
)

// type SensorData struct {
// 	SensorValue int       `json:"sensor_value"`
// 	ID1         int       `json:"id1"`
// 	ID2         string    `json:"id2"`
// 	Timestamp   time.Time `json:"timestamp"`
// }

var ids2 = []string{"A", "B"}

func main() {
	for {
		data := models.SensorData{
			SensorValue: rand.Intn(100),
			ID1:         rand.Intn(3) + 1,
			ID2:         ids2[rand.Intn(2)],
			Timestamp:   time.Now(),
		}

		payload, _ := json.Marshal(data)

		res, err := http.Post("http://localhost:8080/data", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("Error sending:", err)
		} else {
			fmt.Println("Sent:", data)
			res.Body.Close()
		}

		time.Sleep(1 * time.Second)
	}
}
