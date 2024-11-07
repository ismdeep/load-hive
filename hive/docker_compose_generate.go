package hive

import (
	"fmt"

	"github.com/ismdeep/load-hive/pkg/templateutil"
)

func GenerateMainDockerCompose(cfg Config) (string, error) {
	var extraHostsList []string
	for host, ip := range cfg.ExtraHosts {
		extraHostsList = append(extraHostsList, fmt.Sprintf("%v:%v", host, ip))
	}

	text, err := templateutil.Generate(dockerComposeMainContent, map[string]any{
		"HiveID":          cfg.Hive.UUID,
		"UserCount":       cfg.Hive.UserCount,
		"SpawnRate":       cfg.Hive.SpawnRate,
		"WebPort":         cfg.Hive.Port,
		"Target":          cfg.Hive.Target,
		"InternalAPIPort": cfg.Hive.InternalAPIPort,
		"ExtraHosts":      extraHostsList,
	})
	if err != nil {
		return "", err
	}

	return text, nil
}

func GenerateNodeDockerCompose(cfg Config) (string, error) {
	var extraHostsList []string
	for host, ip := range cfg.ExtraHosts {
		extraHostsList = append(extraHostsList, fmt.Sprintf("%v:%v", host, ip))
	}

	text, err := templateutil.Generate(dockerComposeNodeContent, map[string]any{
		"HiveID":          cfg.Hive.UUID,
		"WebPort":         cfg.Hive.Port,
		"MainIP":          cfg.Main.IP,
		"InternalAPIPort": cfg.Hive.InternalAPIPort,
		"ExtraHosts":      extraHostsList,
	})
	if err != nil {
		return "", err
	}
	return text, nil
}
