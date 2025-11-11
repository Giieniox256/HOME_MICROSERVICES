package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/MichaelS11/go-dht"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

var (
	pinLed gpio.PinIO
	// dhtSensor gpio.PinIO
	mu       sync.Mutex
	muDht    sync.Mutex
	ledState bool
	edht     *dht.DHT
)

func initTools() {
	// initialize periph
	if _, err := host.Init(); err != nil {
		fmt.Println("Initialize error:", err)
		return
	}

	err := dht.HostInit()
	if err != nil {
		log.Fatal("HostInit error:", err)
	}
}

func setupGPIO() {
	// Get pin GPIO (ex "GPIO17" on Raspberry Pi)
	pinLed = gpioreg.ByName("GPIO26")
	if pinLed == nil {
		fmt.Println("No gpio fo led found")
		return
	}

	//set pin for dht
	// dhtSensor = gpioreg.ByName("GPIO06")
	// if dhtSensor == nil {
	// 	fmt.Println("No gpio for dht found")
	// }

	const dhtSensor = "GPIO17"
	var err error
	edht, err = dht.NewDHT(dhtSensor, dht.Celsius, "")
	if err != nil {
		log.Fatal("NewDHT error:", err)
	}

	//set pin as output
	if err := pinLed.Out(gpio.Low); err != nil {
		fmt.Println("Error with configuration:", err)
		return
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home Assistant")
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if ledState {
		ledState = false
		pinLed.Out(gpio.Low)
		fmt.Fprint(w, "Led switched off")
	} else {
		ledState = true
		pinLed.Out(gpio.High)
		fmt.Fprint(w, "Led switched on")
	}

}

func tempHumiHandler(w http.ResponseWriter, r *http.Request) {
	muDht.Lock()
	defer mu.Unlock()

	humidity, temperature, err := edht.Read()
	if err != nil {
		http.Error(w, "Error reading DHT11: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Temperature: %.1f°C, Humidity: %.1f%%", temperature, humidity)
}

func main() {
	initTools()
	setupGPIO()

	http.HandleFunc("/my", mainPageHandler)
	http.HandleFunc("/toggle", toggleHandler)
	http.HandleFunc("/temp", tempHumiHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
