package main

import (
	"context"
	"encoding/json"
	"homeiot_bluetooth/lib"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func(c <-chan os.Signal, can context.CancelFunc) {
		<-c
		log.Println("signal received.")
		can()
	}(sigCh, cancel)

	//prepare scan
	sensorCh := make(chan lib.Sensor)
	d, err := linux.NewDevice()
	if err != nil {
		log.Fatalln(err)
	}
	go ble.Scan(ctx, true, makeAdvHandler(sensorCh), filter())
	go print(sensorCh)
	ble.SetDefaultDevice(d)
	<-ctx.Done()
	log.Println("done")
}

func makeAdvHandler(ch chan lib.Sensor) func(ble.Advertisement) {
	return func(a ble.Advertisement) {
		if len(a.ServiceData()) > 0 {
			for _, v := range a.ServiceData() {
				if sensor := lib.GetDevice(a.Addr().String(), v); sensor != nil {
					ch <- sensor
				}
			}
		}
	}
}
func filter() ble.AdvFilter {
	return func(a ble.Advertisement) bool {
		return len(a.ServiceData()) > 0
	}
}
func print(ch <-chan lib.Sensor) {
	for {
		j, err := json.Marshal(<-ch)
		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Printf("%s", j)
	}
}
