package main

import (
	"context"
	"flag"
	"fmt"
	"homeiot_bluetooth/lib"
	"homeiot_bluetooth/task"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"
)

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
func loadConfig(configPath string) (*task.Config, error) {
	if err := validateConfigPath(configPath); err != nil {
		return nil, err
	}
	config := &task.Config{}
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
	mqttTask, err := task.NewMQTTTask(config.MQTTConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer mqttTask.Disconnect()
	if err := mqttTask.Subscribe("adv/+"); err != nil {
		log.Fatalln(err)
	}
	if err := mqttTask.Subscribe("state/+"); err != nil {
		log.Fatalln(err)
	}
	dbTask := task.NewDBTask(config.DBConfig)
	defer dbTask.Disconnect()
	go dbTask.WriteSensorData(mqttTask.MessageCh)

	<-ctx.Done()
	log.Print("done")
}
