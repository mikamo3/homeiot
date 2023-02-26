package task

type BleConfig struct {
	CharacteristicUUID string `yaml:"characteristicUUID"`
	InroomAddr         string `yaml:"inroomAddr"`
}
type MQTTConfig struct {
	Host string `yaml:"host"`
}
type DBConfig struct {
	Host   string `yaml:"host"`
	Token  string `yaml:"token"`
	Org    string `yaml:"org"`
	Bucket string `yaml:"bucket"`
}
type Config struct {
	*BleConfig  `yaml:"bluetooth"`
	*MQTTConfig `yaml:"mqtt"`
	*DBConfig   `yaml:"db"`
}
