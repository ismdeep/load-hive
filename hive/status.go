package hive

import (
	"fmt"
	"strings"

	"github.com/ismdeep/load-hive/pkg/cmdutil"
)

func (receiver *Hive) StatusMain() {
	cfg := receiver.config
	// remote docker-compose ps 'table{{.Name}}\t{{.Status}}'
	output, err := cmdutil.Run(
		receiver.dir,
		nil,
		"ssh",
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("root@%v", cfg.Main.IP),
		fmt.Sprintf("docker ps --format 'table{{.Names}}\\t{{.Status}}'"))
	switch {
	case err != nil:
		fmt.Printf("========== UUID: %v    MAIN: %v FAILED ==========\n", receiver.config.Hive.UUID, cfg.Main.IP)
		fmt.Println(output)
	default:
		fmt.Printf("========== UUID: %v    MAIN: %v ==========\n", receiver.config.Hive.UUID, cfg.Main.IP)
		lines := strings.Split(output, "\n")
		fmt.Println(lines[0])
		for _, s := range lines[1:] {
			if strings.Contains(s, fmt.Sprintf("lh-%v-main", receiver.config.Hive.UUID)) {
				fmt.Println(s)
			}
		}
	}
	fmt.Println()
}

func (receiver *Hive) StatusNode(node NodeConfig) {
	// remote docker-compose ps 'table{{.Name}}\t{{.Status}}'
	output, err := cmdutil.Run(
		receiver.dir,
		nil,
		"ssh",
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("root@%v", node.IP),
		fmt.Sprintf("docker ps --format 'table{{.Names}}\\t{{.Status}}'"))
	switch {
	case err != nil:
		fmt.Printf("========== UUID: %v    NODE: %v FAILED ==========\n", receiver.config.Hive.UUID, node.IP)
		fmt.Println(output)
	default:
		fmt.Printf("========== UUID: %v    NODE: %v ==========\n", receiver.config.Hive.UUID, node.IP)
		lines := strings.Split(output, "\n")
		fmt.Println(lines[0])
		for _, s := range lines[1:] {
			if strings.Contains(s, fmt.Sprintf("lh-%v-node", receiver.config.Hive.UUID)) {
				fmt.Println(s)
			}
		}
	}
	fmt.Println()
}

func (receiver *Hive) StatusNodes() {
	for _, node := range receiver.config.Nodes {
		receiver.StatusNode(node)
	}
}

func (receiver *Hive) Status() {
	receiver.StatusMain()
	receiver.StatusNodes()
}
