package lib

import (
	"testing"

	"github.com/go-ble/ble"
	"github.com/stretchr/testify/assert"
)

func TestGetDevice(t *testing.T) {
	rawdata := []ble.ServiceData{
		{
			UUID: ble.UUID16(0xfd3d),
			Data: []byte{0x48, 0xc0, 0x5b, 0x80},
		},
		{
			UUID: ble.UUID16(0x0d00),
			Data: []byte{0x63, 0xc0, 0x36, 0x49, 0x11, 0x02},
		},
		{
			UUID: ble.UUID16(0xfd3d),
			Data: []byte{0x54, 0x00, 0x64, 0x03, 0x95, 0x2d},
		},
		{
			UUID: ble.UUID16(0x0d00),
			Data: []byte{0x73, 0xc0, 0x64, 0x00, 0x17, 0x0a},
		},
		{
			UUID: ble.UUID16(0xfd3d),
			Data: []byte{0x64, 0x00, 0x64, 0x01, 0x07, 0xab, 0x08, 0x75, 0x14},
		},
		{
			UUID: ble.UUID16(0x2301),
			Data: []byte{0x5a, 0x00, 0xcd},
		},
	}
	tests := []struct {
		name        string
		addr        string
		serviceData ble.ServiceData
		want        interface{}
	}{
		{name: "Bot",
			addr:        "00:11:22:33:44:55",
			serviceData: rawdata[0],
			want: &Bot{SwitchbotDevice: SwitchbotDevice{Device: Device{
				Rawdata: "48c05b80",
				Addr:    "00:11:22:33:44:55"},
				DeviceType: "H",
				Battery:    91},
				Mode:  1,
				State: 1}},
		{name: "Curtain",
			addr:        "00:11:22:33:44:55",
			serviceData: rawdata[1],
			want: &Curtain{SwitchbotDevice: SwitchbotDevice{
				Device: Device{
					Rawdata: "63c036491102",
					Addr:    "00:11:22:33:44:55",
				},
				DeviceType: "c",
				Battery:    54,
			},
				CalibrationSituation:    1,
				DeviceChain:             1,
				LightLevel:              1,
				MotionState:             0,
				Position:                73,
				WheterToAllowConnection: 1,
			}},

		{name: "Meter",
			addr:        "00:11:22:33:44:55",
			serviceData: rawdata[2],
			want: &Meter{SwitchbotDevice: SwitchbotDevice{
				Device: Device{
					Addr:    "00:11:22:33:44:55",
					Rawdata: "54006403952d",
				},
				DeviceType: "T",
				Battery:    100,
			},
				TempertureHighAlert: 0,
				TempertureLowAlert:  0,
				HumidityHighAlert:   0,
				HumidityLowAlert:    0,
				Temperature:         21.3,
				Humidity:            45}},
		{name: "Motion",
			addr:        "00:11:22:33:44:55",
			serviceData: rawdata[3],
			want: &Motion{
				SwitchbotDevice: SwitchbotDevice{
					Device: Device{
						Addr:    "00:11:22:33:44:55",
						Rawdata: "73c06400170a",
					},
					DeviceType: "s",
					Battery:    100,
				},
				IotState:                0,
				LedState:                0,
				LightIntensity:          2,
				RipState:                1,
				SensingDistance:         2,
				SinceLastTriggerRipTime: 23,
				ScopeTested:             1}},
		{name: "Contact",
			addr:        "00:11:22:33:44:55",
			serviceData: rawdata[4],
			want: &Contact{
				SwitchbotDevice: SwitchbotDevice{
					Device: Device{
						Addr:    "00:11:22:33:44:55",
						Rawdata: "6400640107ab087514",
					},
					DeviceType: "d",
					Battery:    100,
				},
				ButtonPushCount: 4,
				EntranceCount:   0,
				GooutCount:      1,
				HalState:        0,
				HalUTC:          2165,
				LightLevel:      1,
				RipState:        0,
				RipUTC:          1963,
				ScopeTested:     0}},
		{name: "M5StackEnv",

			addr:        "00:11:22:33:44:55",
			serviceData: rawdata[5],
			want: &M5StackEnv{
				Device: Device{
					Addr:    "00:11:22:33:44:55",
					Rawdata: "5a00cd",
				},
				DeviceType: "Z",
				Range:      0xcd,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDevice(tt.addr, tt.serviceData)
			assert.Equal(t, tt.want, got)
		})
	}
}
