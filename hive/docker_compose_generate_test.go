package hive

import "testing"

func TestGenerateMainDockerCompose(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				cfg: Config{
					Main: MainConfig{
						IP: "10.20.30.40",
					},
					Nodes: nil,
					Hive: HiveConfig{
						UUID:            "ad58b7c42aaa9c9a41a4190c6e44fb0119644292ebda688e636546fbbea3e8ce",
						Port:            "8080",
						Target:          "http://10.20.30.50:9000",
						InternalAPIPort: "5555",
					},
					ExtraHosts: map[string]string{
						"www.baidu.com": "127.0.0.1",
					},
				},
			},
		},
		{
			name: "",
			args: args{
				cfg: Config{
					Main: MainConfig{
						IP: "10.20.30.40",
					},
					Nodes: nil,
					Hive: HiveConfig{
						UUID:            "ad58b7c42aaa9c9a41a4190c6e44fb0119644292ebda688e636546fbbea3e8ce",
						Port:            "8080",
						Target:          "http://10.20.30.50:9000",
						InternalAPIPort: "5555",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateMainDockerCompose(tt.args.cfg)
			if err != nil {
				t.Errorf("GenerateMainDockerCompose() error = %v", err)
				t.FailNow()
				return
			}
			t.Logf("got = \n%v", got)
		})
	}
}
