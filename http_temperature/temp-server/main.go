package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type IlluminanceData struct {
	Illuminance float64   `json:"illuminance"`
	Time        time.Time `json:"time"`
}

var (
	data []IlluminanceData
	mu   sync.Mutex
)

func illuminanceHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var value IlluminanceData

	err := json.NewDecoder(r.Body).Decode(&value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	value.Time = time.Now()

	mu.Lock()
	data = append(data, value)
	mu.Unlock()

	log.Printf("Illuminance: %.2f lux", value.Illuminance)

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
    <title>ESP32 Illuminance Monitor</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body>

<h2>ESP32 Illuminance Monitor</h2>

<canvas id="chart" width="800" height="400"></canvas>

<script>

const ctx = document.getElementById('chart').getContext('2d');

const chart = new Chart(ctx, {
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            label: 'Illuminance (lux)',
            data: [],
            borderColor: 'orange',
            borderWidth: 2,
            fill: false
        }]
    },
    options: {
        responsive: true
    }
});

async function updateChart() {

    const response = await fetch('/data');
    const values = await response.json();

    chart.data.labels = values.map(v =>
        new Date(v.time).toLocaleTimeString()
    );

    chart.data.datasets[0].data =
        values.map(v => v.illuminance);

    chart.update();
}

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
	http.HandleFunc("/illuminance", illuminanceHandler)
	http.HandleFunc("/data", dataHandler)

	log.Println("Server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}