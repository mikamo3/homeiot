package task

import (
	"context"
	"encoding/json"
	"fmt"
	"homeiot_bluetooth/lib"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/go-ble/ble/linux/hci/cmd"
	"go.uber.org/zap"
)

type BleTask struct {
	context       context.Context
	caracteristic *ble.Characteristic
	client        ble.Client
	SensorCh      chan lib.Sensor
	InRoomCh      chan lib.Sensor
	config        *BleConfig
}

func NewBleTask(ctx context.Context, config *BleConfig) (*BleTask, error) {
	task := &BleTask{
		context:  ctx,
		SensorCh: make(chan lib.Sensor),
		InRoomCh: make(chan lib.Sensor),
		config:   config,
	}
	if err := task.Init(); err != nil {
		return nil, err
	}
	return task, nil
}

func (bt *BleTask) Init() error {
	lib.Logger.Info("ble init")
	options := ble.OptScanParams(cmd.LESetScanParameters{
		LEScanType:           0x01,   // 0x00: passive, 0x01: active
		LEScanInterval:       0x0004, // 0x0004 - 0x4000; N * 0.625msec
		LEScanWindow:         0x0004, // 0x0004 - 0x4000; N * 0.625msec
		OwnAddressType:       0x00,   // 0x00: public, 0x01: random
		ScanningFilterPolicy: 0x00,   // 0x00: accept all, 0x01: ignore non-white-listed.
	})
	d, err := linux.NewDevice(options)
	if err != nil {
		return err
	}
	ble.SetDefaultDevice(d)

	lib.Logger.Info("connect inroom device")
	retryCount := 0
	for retryCount < MAX_RETRIES {
		bt.client, err = ble.Dial(bt.context, ble.NewAddr(bt.config.InroomAddr))
		if err != nil {
			retryCount++
			if retryCount == MAX_RETRIES {
				continue
			}
			lib.Logger.Warn(err.Error())
			time.Sleep(RETRY_WAIT_SEC)
			continue
		} else {
			break
		}
	}
	if err != nil {
		return err
	}

	lib.Logger.Info("find characteristic")
	p, err := bt.client.DiscoverProfile(false)
	if err != nil {
		return err
	}
	targetCID := ble.MustParse(bt.config.CharacteristicUUID)
	bt.caracteristic = p.FindCharacteristic(&ble.Characteristic{UUID: targetCID})
	if bt.caracteristic == nil {
		return fmt.Errorf("characteristic: %v not found", targetCID.String())
	}

	go func() {
		<-bt.client.Disconnected()
		if err := bt.client.CancelConnection(); err != nil {
			lib.Logger.Error(err.Error())
		}
		lib.Logger.Fatal("ble disconnected")
	}()

	return nil
}

func (bt *BleTask) Scan() {
	filter := func(a ble.Advertisement) bool {
		return len(a.ServiceData()) > 0
	}
	handler := func(a ble.Advertisement) {
		for _, v := range a.ServiceData() {
			if sensor := lib.GetDevice(a.Addr().String(), v); sensor != nil {
				bt.SensorCh <- sensor
			}
		}
	}
	lib.Logger.Info("start scan")
	ble.Scan(bt.context, true, handler, filter)
}

func (bt *BleTask) Subscribe() error {
	s := func(addr string, ch chan lib.Sensor) ble.NotificationHandler {
		return func(b []byte) {
			ch <- lib.NewInRoom(addr, b)
		}
	}(bt.client.Addr().String(), bt.InRoomCh)
	if err := bt.client.Subscribe(bt.caracteristic, false, s); err != nil {
		return err
	}

	return nil
}

func (bt *BleTask) SendInroomState(ch <-chan mqtt.Message) {
	for {
		v := <-ch
		var val []byte
		var payload map[string]interface{}
		if err := json.Unmarshal(v.Payload(), &payload); err != nil {
			lib.Logger.Warn(err.Error())
		}
		addr, ok := payload["state"]
		if !ok {
			lib.Logger.Warn("state not found", zap.String("payload", string(v.Payload())))
			continue
		}
		switch addr.(string) {
		case "entrance":
			val = []byte("1")
		default:
			val = []byte("0")
		}
		if err := bt.client.WriteCharacteristic(bt.caracteristic, val, false); err != nil {
			lib.Logger.Warn(err.Error())
		}
	}
}

func (bt *BleTask) Disconnect() error {

	if err := bt.client.CancelConnection(); err != nil {
		return err
	}
	lib.Logger.Info("ble disconnected")
	return nil
}
