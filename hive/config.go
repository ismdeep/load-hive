package hive

type Config struct {
	Main       MainConfig        `yaml:"main"`
	Nodes      []NodeConfig      `yaml:"nodes"`
	Hive       HiveConfig        `yaml:"hive"`
	ExtraHosts map[string]string `yaml:"extra-hosts"`
}

type MainConfig struct {
	IP string `yaml:"ip"`
}

type NodeConfig struct {
	IP string `yaml:"ip"`
}

type HiveConfig struct {
	UUID            string `yaml:"uuid"`
	Port            string `yaml:"port"`
	Target          string `yaml:"target"`
	InternalAPIPort string `yaml:"internal-api-port"`
}
