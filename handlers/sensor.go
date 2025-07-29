package handlers

import (
	"fmt"
	"net/http"
	"sensor-server/db"
	"sensor-server/models"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func InsertSensor(c echo.Context) error {
	var data models.SensorData

	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	current := time.Now()

	query := `INSERT INTO sensor_data(sensor_value, id1, id2, timestamp)VALUES(?, ?, ?, ?)`
	_, err = db.DB.Exec(query, data.SensorValue, data.ID1, data.ID2, current)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to insert"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":   "Sensor data inserted successfully",
		"inserted":  data,
		"timestamp": current,
	})
}

func GetSensor(c echo.Context) error {
	id1 := c.QueryParam("ID1")
	id2 := c.QueryParam("ID2")
	start := c.QueryParam("start_timestamp")
	end := c.QueryParam("end_timestamp")

	query := "SELECT sensor_value, id1, id2, timestamp FROM sensor_data WHERE 1=1"
	args := []interface{}{}

	if id1 != "" {
		query += " AND id1 = ?"
		args = append(args, id1)
	}
	if id2 != "" {
		query += " AND id2 = ?"
		args = append(args, id2)
	}
	if start != "" && end != "" {
		startTime, err1 := parseUnix(start)
		endTime, err2 := parseUnix(end)
		if err1 != nil || err2 != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "invalid timestamp format.",
			})
		}
		query += " AND timestamp BETWEEN ? AND ?"
		args = append(args, startTime, endTime)
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "falied to retrive data from db"})
	}
	defer rows.Close()

	var res []models.SensorData
	for rows.Next() {
		var s models.SensorData
		err := rows.Scan(&s.SensorValue, &s.ID1, &s.ID2, &s.Timestamp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		res = append(res, s)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "data retrieved successfully",
		"data":    res,
		"count":   len(res),
	})
}

func GetLatest(c echo.Context) error {
	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "5"
	}

	query := "SELECT sensor_value, id1, id2, timestamp FROM sensor_data ORDER BY timestamp DESC LIMIT ?"
	rows, err := db.DB.Query(query, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "falied to retrive data from db"})
	}
	defer rows.Close()

	var res []models.SensorData
	for rows.Next() {
		var g models.SensorData
		err := rows.Scan(&g.SensorValue, &g.ID1, &g.ID2, &g.Timestamp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		res = append(res, g)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "data retrieved successfully",
		"data":    res,
		"count":   len(res),
	})
}

func GetGrouped(c echo.Context) error {
	groupBy := c.QueryParam("group_by")

	allow := map[string]bool{
		"id1": true,
		"id2": true,
	}

	if !allow[groupBy] {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid group_by parameter.",
		})
	}

	query := fmt.Sprintf("SELECT %s AS grouped, COUNT(*) AS total FROM sensor_data GROUP BY %s", groupBy, groupBy)
	rows, err := db.DB.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "falied to retrive data from db"})
	}
	defer rows.Close()

	type GroupDesc struct {
		Group string `json:"grouped"`
		Total string `json:"total"`
	}

	var res []GroupDesc
	for rows.Next() {
		var s GroupDesc
		err := rows.Scan(&s.Group, &s.Total)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		res = append(res, s)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "data retrieved successfully",
		"data":    res,
		"count":   len(res),
	})
}

func GetStats(c echo.Context) error {
	start := c.QueryParam("start")
	end := c.QueryParam("end")

	start_tm, err1 := parseUnix(start)
	end_tm, err2 := parseUnix(end)
	if err1 != nil || err2 != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid timestamp format.",
		})
	}

	query := "SELECT COUNT(*), AVG(sensor_value), MIN(sensor_value), MAX(sensor_value) FROM sensor_data WHERE timestamp BETWEEN ? AND ?"
	row := db.DB.QueryRow(query, start_tm, end_tm)

	type SensorStats struct {
		Count int     `json:"count"`
		Avg   float64 `json:"avg"`
		Min   int     `json:"min"`
		Max   int     `json:"max"`
	}

	var s SensorStats
	err := row.Scan(&s.Count, &s.Avg, &s.Min, &s.Max)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "stats retrieved successfully",
		"data":    s,
		"start":   start_tm,
		"end":     end_tm,
	})
}

func UpdateSensor(c echo.Context) error {
	id2 := c.QueryParam("id2")
	if id2 == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "missing id2",
		})
	}

	var data struct {
		SensorValue int `json:"sensor_value"`
		ID1         int `json:"id1"`
	}

	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid input",
		})
	}

	query := `UPDATE sensor_data SET sensor_value = ?, id1 = ? WHERE id2 = ?`
	result, err := db.DB.Exec(query, data.SensorValue, data.ID1, id2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to update data",
		})
	}

	changed, _ := result.RowsAffected()
	if changed == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "No record found.",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "data updated successfully",
		"id2":     id2,
		"updated": data,
	})
}

func DelSensor(c echo.Context) error {
	id2 := c.QueryParam("id2")
	if id2 == "" {
		return c.JSON(http.StatusBadRequest, "missing id2")
	}

	query := "DELETE FROM sensor_data WHERE id2=?"
	res, err := db.DB.Exec(query, id2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	changed, _ := res.RowsAffected()
	if changed == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "No record found.",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sensor data deleted successfully",
		"id2":     id2,
	})
}

func parseUnix(s string) (time.Time, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}

//------
