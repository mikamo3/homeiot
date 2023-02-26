package task

import (
	"encoding/json"
	"homeiot_bluetooth/lib"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type MQTTTask struct {
	client    mqtt.Client
	MessageCh chan mqtt.Message
	config    *MQTTConfig
}

func NewMQTTTask(config *MQTTConfig) (*MQTTTask, error) {
	task := &MQTTTask{
		MessageCh: make(chan mqtt.Message),
		config:    config,
	}

	if err := task.Init(); err != nil {
		return nil, err
	}
	return task, nil
}

func (mt *MQTTTask) Init() error {
	lib.Logger.Info("mqtt init")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mt.config.Host)
	mt.client = mqtt.NewClient(opts)
	if token := mt.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (mt *MQTTTask) Disconnect() {
	mt.client.Disconnect(1000)
	lib.Logger.Info("mqtt disconnected")

}

func (mt *MQTTTask) Subscribe(topic string) error {
	lib.Logger.Info("add subscribe", zap.String("topic", topic))
	//subscribe
	if sToken := mt.client.Subscribe(topic, 0, func(c mqtt.Client, m mqtt.Message) {
		mt.MessageCh <- m
	}); sToken.Wait() && sToken.Error() != nil {
		return sToken.Error()
	}
	return nil
}

func (mt *MQTTTask) PublishAdvertisement(ch <-chan lib.Sensor) {
	lib.Logger.Info("mqtt publish")
	for {
		s := <-ch
		j, err := json.Marshal(s)
		if err != nil {
			log.Fatalln(err)
			continue
		}
		mt.client.Publish("adv/"+s.DeviceName(), 0, false, j)
	}
}
