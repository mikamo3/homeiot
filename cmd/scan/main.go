package main

import (
	"context"
	"flag"
	"fmt"
	"homeiot_bluetooth/lib"
	"homeiot_bluetooth/task"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"
)

type scanConfig struct {
	*task.BleConfig  `yaml:"bluetooth"`
	*task.MQTTConfig `yaml:"mqtt"`
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a Directory", path)
	}
	return nil
}
func parseFlags() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yml", "config path")
	flag.Parse()
	if err := validateConfigPath(configPath); err != nil {
		return "", err
	}
	return configPath, nil
}
func loadConfig(configPath string) (*scanConfig, error) {
	if err := validateConfigPath(configPath); err != nil {
		return nil, err
	}
	config := &scanConfig{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
func main() {
	lib.Logger.Info("start")
	defer lib.Logger.Info("end")
	defer lib.Logger.Sync()
	configPath, err := parseFlags()
	if err != nil {
		lib.Logger.Fatal(err.Error())
	}
	config, err := loadConfig(configPath)
	if err != nil {
		lib.Logger.Fatal(err.Error())
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	//mqtt
	mqttTask, err := task.NewMQTTTask(config.MQTTConfig)
	if err != nil {
		lib.Logger.Fatal(err.Error())
	}
	defer mqttTask.Disconnect()
	if err := mqttTask.Subscribe("state/state_inroom"); err != nil {
		lib.Logger.Fatal(err.Error())
	}

	//ble
	bleTask, err := task.NewBleTask(ctx, config.BleConfig)
	if err != nil {
		lib.Logger.Fatal(err.Error())
	}
	defer func() {
		if err := bleTask.Disconnect(); err != nil {
			lib.Logger.Fatal(err.Error())
		}
	}()
	if err := bleTask.Subscribe(); err != nil {
		lib.Logger.Fatal(err.Error())
	}
	go mqttTask.PublishAdvertisement(bleTask.SensorCh)
	go mqttTask.PublishAdvertisement(bleTask.InRoomCh)
	go bleTask.Scan()
	go bleTask.SendInroomState(mqttTask.MessageCh)
	<-ctx.Done()
}
