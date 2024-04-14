package hive

import _ "embed"

//go:embed main/docker-compose.yaml
var dockerComposeMainContent string

//go:embed node/docker-compose.yaml
var dockerComposeNodeContent string
