package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type TemperatureData struct {
	Temperature float64   `json:"temperature"`
	Time        time.Time `json:"time"`
}

var (
	data []TemperatureData
	mu   sync.Mutex
)

func temperatureHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var temp TemperatureData

	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	temp.Time = time.Now()

	mu.Lock()
	data = append(data, temp)
	mu.Unlock()

	log.Printf("Temperature received: %.2f °C", temp.Temperature)

	w.WriteHeader(http.StatusOK)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>ESP32 Temperature Monitor</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body>

<h2>ESP32 Temperature Monitor</h2>

<canvas id="tempChart" width="800" height="400"></canvas>

<script>

const ctx = document.getElementById('tempChart').getContext('2d');

const chart = new Chart(ctx, {
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            label: 'Temperature °C',
            data: [],
            borderColor: 'red',
            borderWidth: 2,
            fill: false
        }]
    },
    options: {
        responsive: true,
        scales: {
            y: {
                beginAtZero: false
            }
        }
    }
});

async function updateChart() {

    const response = await fetch('/data');
    const values = await response.json();

    chart.data.labels = values.map(v =>
        new Date(v.time).toLocaleTimeString()
    );

    chart.data.datasets[0].data = values.map(v => v.temperature);

    chart.update();
}

// Refresh every 5 seconds
setInterval(updateChart, 5000);

updateChart();

</script>

</body>
</html>
`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/temperature", temperatureHandler)
	http.HandleFunc("/data", dataHandler)

	log.Println("Server started on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}