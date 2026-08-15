package hive

type Config struct {
	Main       MainConfig        `yaml:"main"`
	Nodes      []NodeConfig      `yaml:"nodes"`
	Hive       HiveConfig        `yaml:"hive"`
	ExtraHosts map[string]string `yaml:"extra_hosts"`
}

type MainConfig struct {
	IP string `yaml:"ip"`
}

type NodeConfig struct {
	IP     string `yaml:"ip"`
	Worker int    `yaml:"worker"`
}

type HiveConfig struct {
	UUID            string `yaml:"uuid"`
	Image           string `yaml:"image"`
	Port            string `yaml:"port"`
	Target          string `yaml:"target"`
	InternalAPIPort string `yaml:"internal_api_port"`
	UserCount       int    `yaml:"user_count"`
	SpawnRate       int    `yaml:"spawn_rate"`
}
