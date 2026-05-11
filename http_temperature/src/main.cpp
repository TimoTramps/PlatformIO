#include <WiFi.h>
#include <HTTPClient.h>

// =========================
// WiFi
// =========================
const char* ssid = "pomaranca";
const char* password = "limonada";

// Server URL
const char* serverUrl = "http://192.168.1.100:8080/illuminance";

// Analog pin
const int sensorPin = 13;

void setup() {

    Serial.begin(115200);

    analogReadResolution(12);

    WiFi.begin(ssid, password);

    Serial.print("Connecting to WiFi");

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    Serial.println();
    Serial.println("Connected!");
}

void loop() {

    // Read analog value
    int rawValue = analogRead(sensorPin);

    // Convert ADC value to estimated lux
    // Adjust scaling depending on your sensor
    float lux = map(rawValue, 0, 4095, 0, 1000);

    Serial.print("Raw ADC: ");
    Serial.print(rawValue);

    Serial.print("  Illuminance: ");
    Serial.print(lux);
    Serial.println(" lux");

    if (WiFi.status() == WL_CONNECTED) {

        HTTPClient http;

        http.begin(serverUrl);
        http.addHeader("Content-Type", "application/json");

        String json = "{\"illuminance\": " + String(lux, 2) + "}";

        int responseCode = http.POST(json);

        Serial.print("HTTP Response: ");
        Serial.println(responseCode);

        http.end();
    }

    delay(5000);
}