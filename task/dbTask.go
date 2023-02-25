package task

import (
	"context"
	"encoding/json"
	"homeiot_bluetooth/lib"
	"regexp"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go.uber.org/zap"
)

type DBTask struct {
	client          influxdb2.Client
	writeAPI        api.WriteAPIBlocking
	sensorDataCache map[string]*sensorDataCache
	config          *DBConfig
}
type sensorDataCache struct {
	date       time.Time
	macaddress string
}

func NewDBTask(config *DBConfig) *DBTask {
	task := &DBTask{
		sensorDataCache: make(map[string]*sensorDataCache),
		config:          config,
	}
	task.Init()
	return task
}
func (dt *DBTask) Init() {
	lib.Logger.Info("db init")
	dt.client = influxdb2.NewClient("http://localhost:8086", "eqZDJyLTi0trnt09RxtkrfinbxeHMofrTvWfzSka1TMSGxrXJjNDlazqAlDIJkXwnmkD7ltx03UiyQmuM2j5Wg==")
	dt.writeAPI = dt.client.WriteAPIBlocking("homeiot", "homeiot3")
}
func (dt *DBTask) Disconnect() {
	lib.Logger.Info("db disconnected")
	dt.client.Close()
}

func (dt *DBTask) WriteSensorData(ch <-chan mqtt.Message) {
	for {
		msg := <-ch
		now := time.Now()
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			lib.Logger.Warn(err.Error())
		}
		regAdv := regexp.MustCompile("^adv/")
		regState := regexp.MustCompile("^state/")
		var p *write.Point
		var topic string
		var tags map[string]string
		if regAdv.Match([]byte(msg.Topic())) {
			topic = regAdv.ReplaceAllString(msg.Topic(), "")
			addr, ok := payload["macaddress"]
			if !ok {
				lib.Logger.Warn("macaddress not found ", zap.String("payload", string(msg.Payload())))
				continue
			}
			if sc, ok := dt.sensorDataCache[topic]; ok && sc != nil {
				if sc.date.Add(1 * time.Second).After(now) {
					continue
				}
			}

			dt.sensorDataCache[topic] = &sensorDataCache{
				macaddress: addr.(string),
				date:       now,
			}
			tags = map[string]string{"macaddress": addr.(string)}
		} else if regState.Match([]byte(msg.Topic())) {
			topic = regState.ReplaceAllString(msg.Topic(), "")
			tags = nil
		} else {
			lib.Logger.Warn("invalid data", zap.String("payload", string(msg.Payload())))
			continue
		}

		p = influxdb2.NewPoint(topic,
			tags,
			payload,
			now)
		// write point immediately
		err := dt.writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			lib.Logger.Fatal(err.Error())
			continue
		}
	}
}
