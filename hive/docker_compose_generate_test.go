package hive

import (
	"fmt"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestGenerateMainDockerCompose(t *testing.T) {
	cfg := Config{
		Main: MainConfig{IP: "10.20.30.40"},
		Hive: HiveConfig{
			UUID:            "hive1",
			Port:            "8080",
			Target:          "http://10.20.30.50:9000",
			InternalAPIPort: "5555",
			UserCount:       200,
			SpawnRate:       20,
		},
		ExtraHosts: map[string]string{
			"api.example.com": "127.0.0.1",
			"db.example.com":  "10.0.0.2",
		},
	}

	got, err := GenerateMainDockerCompose(cfg)
	if err != nil {
		t.Fatalf("GenerateMainDockerCompose() error = %v", err)
	}

	services := mustComposeServices(t, got)
	service := services["lh-hive1-main"]
	if service.ContainerName != "lh-hive1-main" {
		t.Fatalf("container_name = %q, want %q", service.ContainerName, "lh-hive1-main")
	}
	if service.Image != "hub.deepin.com/library/locustio/locust:latest" {
		t.Fatalf("image = %q", service.Image)
	}
	if service.Command != "-f /mnt/locust/locustfile.py --master -u 200 -r 20 -P 8080 --host http://10.20.30.50:9000" {
		t.Fatalf("command = %q", service.Command)
	}
	assertStringSliceContains(t, service.Ports, "8080:8080")
	assertStringSliceContains(t, service.Ports, "5555:5557")
	assertStringSliceContains(t, service.ExtraHosts, "api.example.com:127.0.0.1")
	assertStringSliceContains(t, service.ExtraHosts, "db.example.com:10.0.0.2")
}

func TestGenerateMainDockerComposeWithoutExtraHosts(t *testing.T) {
	got, err := GenerateMainDockerCompose(Config{
		Hive: HiveConfig{
			UUID:            "hive1",
			Port:            "8080",
			Target:          "http://example.com",
			InternalAPIPort: "5555",
			UserCount:       1,
			SpawnRate:       1,
		},
	})
	if err != nil {
		t.Fatalf("GenerateMainDockerCompose() error = %v", err)
	}
	if strings.Contains(got, "extra_hosts:") {
		t.Fatalf("compose should not contain extra_hosts when cfg.ExtraHosts is empty:\n%s", got)
	}
}

func TestGenerateNodeDockerCompose(t *testing.T) {
	cfg := Config{
		Main: MainConfig{IP: "10.20.30.40"},
		Hive: HiveConfig{
			UUID:            "hive1",
			Port:            "8080",
			InternalAPIPort: "5555",
		},
		ExtraHosts: map[string]string{
			"api.example.com": "127.0.0.1",
		},
	}

	got, err := GenerateNodeDockerCompose(cfg)
	if err != nil {
		t.Fatalf("GenerateNodeDockerCompose() error = %v", err)
	}

	services := mustComposeServices(t, got)
	service := services["lh-hive1-node"]
	if service.Image != "hub.deepin.com/library/locustio/locust:latest" {
		t.Fatalf("image = %q", service.Image)
	}
	if service.Command != "-f /mnt/locust/locustfile.py --worker --master-host=10.20.30.40 --master-port=5555" {
		t.Fatalf("command = %q", service.Command)
	}
	assertStringSliceContains(t, service.ExtraHosts, "api.example.com:127.0.0.1")
}

type composeFile struct {
	Services map[string]composeService `yaml:"services"`
}

type composeService struct {
	ContainerName string   `yaml:"container_name"`
	Image         string   `yaml:"image"`
	Command       string   `yaml:"command"`
	Ports         []string `yaml:"ports"`
	ExtraHosts    []string `yaml:"extra_hosts"`
}

func mustComposeServices(t *testing.T, raw string) map[string]composeService {
	t.Helper()

	var cfg composeFile
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("generated compose is not valid yaml: %v\n%s", err, raw)
	}
	if len(cfg.Services) == 0 {
		t.Fatalf("generated compose has no services:\n%s", raw)
	}
	return cfg.Services
}

func assertStringSliceContains(t *testing.T, values []string, want string) {
	t.Helper()

	for _, value := range values {
		if value == want {
			return
		}
	}
	t.Fatalf("%v does not contain %q", values, want)
}

func TestGenerateDockerComposeOutputIsDeterministicEnoughForYAML(t *testing.T) {
	cfg := Config{
		Hive: HiveConfig{
			UUID:            "hive1",
			Port:            "8080",
			Target:          "http://example.com",
			InternalAPIPort: "5555",
		},
	}

	mainCompose, err := GenerateMainDockerCompose(cfg)
	if err != nil {
		t.Fatalf("GenerateMainDockerCompose() error = %v", err)
	}
	nodeCompose, err := GenerateNodeDockerCompose(cfg)
	if err != nil {
		t.Fatalf("GenerateNodeDockerCompose() error = %v", err)
	}

	for name, raw := range map[string]string{"main": mainCompose, "node": nodeCompose} {
		t.Run(name, func(t *testing.T) {
			if _, ok := mustComposeServices(t, raw)[fmt.Sprintf("lh-%s-%s", cfg.Hive.UUID, name)]; !ok {
				t.Fatalf("service lh-%s-%s not found in:\n%s", cfg.Hive.UUID, name, raw)
			}
		})
	}
}
