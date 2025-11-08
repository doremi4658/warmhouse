package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Sensor struct {
	ID           string  `json:"id"`
	CurrentTemp  float64 `json:"current_temperature"`
	TargetTemp   float64 `json:"target_temperature"`
	Location     string  `json:"location"`
	LastUpdate   string  `json:"last_update"`
	Status       string  `json:"status"`
}

var sensors = map[string]*Sensor{
	"living_room": {
		ID:          "sensor_001",
		CurrentTemp: 22.5,
		TargetTemp:  22.0,
		Location:    "living_room",
		LastUpdate:  time.Now().Format(time.RFC3339),
		Status:      "online",
	},
	"kitchen": {
		ID:          "sensor_002",
		CurrentTemp: 24.0,
		TargetTemp:  23.0,
		Location:    "kitchen",
		LastUpdate:  time.Now().Format(time.RFC3339),
		Status:      "online",
	},
}

// Упрощенный обработчик для установки температуры
func setTemperatureHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("📨 Received %s request to %s", r.Method, r.URL.String())

    // Разрешаем только POST метод
    if r.Method != http.MethodPost {
        errorMsg := fmt.Sprintf("Method %s not allowed. Please use POST.", r.Method)
        log.Printf("❌ %s", errorMsg)
        http.Error(w, errorMsg, http.StatusMethodNotAllowed)
        return
    }

    // Получаем параметры из query string
    location := r.URL.Query().Get("location")
    tempStr := r.URL.Query().Get("temperature")

    log.Printf("🔧 Parameters - location: %s, temperature: %s", location, tempStr)

    if location == "" || tempStr == "" {
        errorMsg := "Missing location or temperature parameters"
        log.Printf("❌ %s", errorMsg)
        http.Error(w, errorMsg, http.StatusBadRequest)
        return
    }

    temperature, err := strconv.ParseFloat(tempStr, 64)
    if err != nil {
        errorMsg := fmt.Sprintf("Invalid temperature value: %s", tempStr)
        log.Printf("❌ %s", errorMsg)
        http.Error(w, errorMsg, http.StatusBadRequest)
        return
    }

    if temperature < 10 || temperature > 35 {
        errorMsg := "Temperature must be between 10 and 35 degrees"
        log.Printf("❌ %s", errorMsg)
        http.Error(w, errorMsg, http.StatusBadRequest)
        return
    }

    // Находим или создаем сенсор
    sensor, exists := sensors[location]
    if !exists {
        sensor = &Sensor{
            ID:         fmt.Sprintf("sensor_%03d", len(sensors)+1),
            Location:   location,
            Status:     "online",
            LastUpdate: time.Now().Format(time.RFC3339),
        }
        sensors[location] = sensor
        log.Printf("✅ Created new sensor: %s", sensor.ID)
    }

    // Устанавливаем температуру
    sensor.TargetTemp = temperature
    sensor.CurrentTemp = temperature
    sensor.LastUpdate = time.Now().Format(time.RFC3339)

    log.Printf("✅ Temperature set - Location: %s, Temp: %.1f°C, Sensor: %s",
        location, temperature, sensor.ID)

    // Возвращаем успешный ответ
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    response := map[string]interface{}{
        "status":       "success",
        "message":      "Temperature successfully set",
        "location":     location,
        "target_temp":  temperature,
        "current_temp": sensor.CurrentTemp,
        "sensor_id":    sensor.ID,
        "timestamp":    sensor.LastUpdate,
    }

    json.NewEncoder(w).Encode(response)
}

// Обработчик для получения температуры
func getTemperatureHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    location := r.URL.Query().Get("location")

    if location == "" {
        // Возвращаем все сенсоры
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "sensors": sensors,
            "count":   len(sensors),
        })
        return
    }

    sensor, exists := sensors[location]
    if !exists {
        http.Error(w, "Sensor not found for location: "+location, http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "location":     location,
        "current_temp": sensor.CurrentTemp,
        "target_temp":  sensor.TargetTemp,
        "sensor_id":    sensor.ID,
        "last_update":  sensor.LastUpdate,
        "status":       sensor.Status,
    })
}

// Health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":        "healthy",
        "service":       "monolith",
        "timestamp":     time.Now().Format(time.RFC3339),
        "sensors_count": len(sensors),
        "endpoints": map[string]string{
            "POST /temperature/set": "Set temperature",
            "GET /temperature/get":  "Get temperature",
            "GET /health":           "Health check",
        },
    })
}

func main() {
    // Настройка маршрутов
    http.HandleFunc("/temperature/set", setTemperatureHandler)
    http.HandleFunc("/temperature/get", getTemperatureHandler)
    http.HandleFunc("/health", healthHandler)

    port := ":8080"
    log.Printf("🚀 Monolith server started on port %s", port)
    log.Printf("📡 Available endpoints:")
    log.Printf("   POST http://localhost:8080/temperature/set?location=bedroom&temperature=23.5")
    log.Printf("   GET  http://localhost:8080/temperature/get?location=bedroom")
    log.Printf("   GET  http://localhost:8080/health")

    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal("Server error:", err)
    }
}