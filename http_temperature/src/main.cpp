#include <WiFi.h>
#include <HTTPClient.h>
#include <SmoothThermistor.h>

// =========================
// WiFi configuration
// =========================
const char* ssid = "pomaranca";
const char* password = "limonada";

// Go server address
const char* serverUrl = "http://192.168.1.100:8080/temperature";

// =========================
// Thermistor configuration
// =========================

// ADC pin
const int thermistorPin = 34;

// Thermistor parameters
// Example: 10k thermistor with Beta 3950
const float SERIES_RESISTOR = 10000.0;
const float THERMISTOR_NOMINAL = 10000.0;
const float TEMPERATURE_NOMINAL = 25.0;
const float BETA_COEFFICIENT = 3950.0;

// Create SmoothThermistor object
SmoothThermistor thermistor(
    thermistorPin,
    ADC_12BIT,
    5, // number of samples
    SERIES_RESISTOR,
    THERMISTOR_NOMINAL,
    TEMPERATURE_NOMINAL,
    BETA_COEFFICIENT
);

void setup() {
    Serial.begin(115200);

    WiFi.begin(ssid, password);

    Serial.print("Connecting to WiFi");

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    Serial.println();
    Serial.println("Connected!");
    Serial.print("IP Address: ");
    Serial.println(WiFi.localIP());
}

void loop() {

    // Read temperature
    float temperature = thermistor.temperature();

    Serial.print("Temperature: ");
    Serial.print(temperature);
    Serial.println(" °C");

    // Send data if WiFi connected
    if (WiFi.status() == WL_CONNECTED) {

        HTTPClient http;

        http.begin(serverUrl);
        http.addHeader("Content-Type", "application/json");

        // JSON payload
        String json = "{\"temperature\": " + String(temperature, 2) + "}";

        int httpResponseCode = http.POST(json);

        Serial.print("HTTP Response code: ");
        Serial.println(httpResponseCode);

        http.end();
    }

    // Send every 5 seconds
    delay(5000);
}