
# 📡 Sensor Data API (Go + Echo + MySQL)

A RESTful backend service for handling sensor data — built using **Golang**, **Echo framework**, and **MySQL**. This server accepts, stores, updates, deletes, and analyzes sensor data with support for filtering, grouping, and basic statistics.

---

## 🧪 Postman Collection
👉 Import the Postman collection to test the API easily:

📎 [Postman Collection Link](https://harshit-6003987.postman.co/workspace/harshit's-Workspace~bb3b2062-7320-454a-8fde-febd3854d040/collection/43825972-a3649ef1-5cd6-455c-9ae9-8c934ad90d58?action=share&creator=43825972)

---

## 🚀 Features

- Insert new sensor data
- Retrieve filtered sensor records
- Fetch latest sensor entries
- Update existing sensor data
- Delete sensor data by ID
- Compute statistics (count, average, min, max)
- Group data by `id1` or `id2`

---

## 🧰 Tech Stack

- **Language**: Go (Golang)
- **Framework**: Echo
- **Database**: MySQL
- **ORM/DB Driver**: `github.com/go-sql-driver/mysql`
- **Environment Handling**: `github.com/joho/godotenv`

---

## 📁 Project Structure

```bash
sensor-server/
├── db/               # Database connection (InitDB)
│   └── db.go
├── models/           # Sensor data model
│   └── sensor.go
├── routes/           # Route handlers
│   └── handlers.go
├── main.go           # Entry point
├── go.mod / go.sum   # Dependencies
├── .env              # Environment variables (NOT COMMITTED)
└── .gitignore        # Files to ignore in Git
```

---

## 🧪 API Endpoints

### 📥 Insert Sensor Data

`POST /insert`

```json
{
  "sensor_value": 85,
  "id1": 1,
  "id2": "A"
}
```

---

### 📤 Get Filtered Sensor Data

`GET /data?ID1=1&ID2=A&start_timestamp=unix&end_timestamp=unix`

Returns records based on optional filters.

---

### 🕓 Get Latest N Records

`GET /latest?limit=5`  
Default limit is `5` if not provided.

---

### ✏️ Update Sensor Data

`PUT /update?id2=A`

```json
{
  "sensor_value": 90,
  "id1": 2
}
```

---

### 🗑️ Delete Sensor Data

`DELETE /delete?id2=A`

---

### 📊 Get Statistics

`GET /stats?start=unix&end=unix`

Returns: count, average, min, max sensor values.

---

## 📦 Environment Variables

Create a `.env` file in the root:

```env
DB_USER=root
DB_PASS=rootpass123
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=sensordb
```

---

## 🛠️ Setup Instructions

### 1. Clone the Repo

```bash
git clone https://github.com/Harshit0741/sensor-server.git
cd sensor-server
```

### 2. Set up `.env`

Create a `.env` file with DB credentials.

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Server

```bash
go run main.go
```

---

## 📝 Example MySQL Table

```sql
CREATE TABLE sensor_data (
  sensor_value INT,
  id1 INT,
  id2 VARCHAR(10),
  timestamp DATETIME
);
```

