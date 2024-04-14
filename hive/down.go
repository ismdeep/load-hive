package hive

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ismdeep/load-hive/pkg/cmdutil"
)

func (receiver *Hive) DownMain() error {
	// stop services
	if output, err := cmdutil.Run(receiver.dir, nil, "ssh",
		fmt.Sprintf("root@%v", receiver.config.Main.IP),
		fmt.Sprintf("docker-compose --project-directory '%v' down", receiver.MainDir())); err != nil {

		switch {
		case strings.Contains(output, "no configuration file provided"):
			fmt.Println("[WARN] failed to down docker compose file, err:", err.Error())
		default:
			return errors.Join(errors.New(output), err)
		}
	}
	Log("INFO", receiver.config.Main.IP, "main dir %v docker-compose down success", receiver.MainDir())

	// remove files
	if output, err := cmdutil.Run(receiver.dir, nil, "ssh",
		fmt.Sprintf("root@%v", receiver.config.Main.IP),
		fmt.Sprintf("rm -rf '%v'", receiver.MainDir())); err != nil {
		return errors.Join(errors.New(output), err)
	}
	Log("INFO", receiver.config.Main.IP, "main dir %v removed", receiver.MainDir())

	return nil
}

func (receiver *Hive) DownNode(node NodeConfig) error {
	// stop services
	if output, err := cmdutil.Run(receiver.dir, nil, "ssh",
		fmt.Sprintf("root@%v", node.IP),
		fmt.Sprintf("docker-compose --project-directory '%v' down", receiver.NodeDir())); err != nil {
		switch {
		case strings.Contains(output, "no configuration file provided"):
			fmt.Println("[WARN] failed to down docker compose file, err:", err.Error())
		default:
			return errors.Join(errors.New(output), err)
		}
	}
	Log("INFO", node.IP, "node dir %v docker-compose down success", receiver.MainDir())

	// remove files
	if output, err := cmdutil.Run(receiver.dir, nil, "ssh",
		fmt.Sprintf("root@%v", node.IP),
		fmt.Sprintf("rm -rf '%v'", receiver.NodeDir())); err != nil {
		return errors.Join(errors.New(output), err)
	}
	Log("INFO", node.IP, "node dir %v removed", receiver.MainDir())

	return nil
}

func (receiver *Hive) DownNodes() error {
	var errLst []error
	var wg sync.WaitGroup
	for _, node := range receiver.config.Nodes {
		wg.Add(1)
		go func(node NodeConfig) {
			defer wg.Done()
			if err := receiver.DownNode(node); err != nil {
				errLst = append(errLst, err)
				return
			}
		}(node)
	}
	wg.Wait()
	return nil
}

func (receiver *Hive) Down() error {
	if err := receiver.DownMain(); err != nil {
		return fmt.Errorf("failed to down main, err: %v", err)
	}
	if err := receiver.DownNodes(); err != nil {
		return fmt.Errorf("failed to down nodes, err: %v", err)
	}
	Log("INFO", receiver.config.Hive.UUID, "all nodes are removed")
	return nil
}
