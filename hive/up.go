package hive

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ismdeep/load-hive/pkg/cmdutil"
)

func (receiver *Hive) UpMain() error {
	cfg := receiver.config

	// create remote dir
	if output, err := cmdutil.Run(receiver.dir, nil, "ssh",
		fmt.Sprintf("root@%v", cfg.Main.IP),
		fmt.Sprintf("mkdir -p '%v'", receiver.MainDir())); err != nil {
		return errors.Join(errors.New(output), err)
	}
	Log("INFO", cfg.Main.IP, "main dir %v prepared", receiver.MainDir())

	// rsync folders
	if output, err := cmdutil.Run(receiver.dir, nil, "rsync",
		"-avz",
		fmt.Sprintf("%v/", receiver.dir),
		fmt.Sprintf("root@%v:%v/", cfg.Main.IP, receiver.MainDir())); err != nil {
		return errors.Join(errors.New(output), err)
	}
	Log("INFO", cfg.Main.IP, "main dir %v synced", receiver.MainDir())

	// generate docker-compose content
	text, err := GenerateMainDockerCompose(*receiver.config)
	if err != nil {
		return err
	}

	// create temp docker-compose.yaml
	tmpdir, err := os.MkdirTemp("", "load-hive-")
	if err != nil {
		return fmt.Errorf("failed to create temp dir, err: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			fmt.Printf("[WARN] failed to remove temp dir, err: %v", err)
		}
	}()
	if err := os.WriteFile(filepath.Join(tmpdir, "docker-compose.yaml"), []byte(text), 0750); err != nil {
		return fmt.Errorf("failed to write temp docker-compose.yaml, err: %v", err)
	}

	// copy docker-compose.yaml for main
	if output, err := cmdutil.Run(receiver.dir, nil, "rsync", "-avz",
		fmt.Sprintf("%v/docker-compose.yaml", tmpdir),
		fmt.Sprintf("root@%v:%v/docker-compose.yaml", cfg.Main.IP, receiver.MainDir())); err != nil {
		return errors.Join(errors.New(output), err)
	}
	Log("INFO", cfg.Main.IP, "main dir %v docker-compose.yaml prepared", receiver.MainDir())

	// remote docker-compose up -d
	if output, err := cmdutil.Run(receiver.dir, nil, "ssh",
		fmt.Sprintf("root@%v", cfg.Main.IP),
		fmt.Sprintf("docker-compose --project-directory '%v' up -d --force-recreate", receiver.MainDir())); err != nil {
		return errors.Join(
			errors.New("failed to up main service"),
			errors.New(output),
			err)
	}
	Log("INFO", cfg.Main.IP, "main dir %v docker-compose up -d success", receiver.MainDir())

	Log("INFO", receiver.config.Hive.UUID, "main %v is up", receiver.config.Main.IP)

	return nil
}

func (receiver *Hive) UpNode(node NodeConfig) error {
	// create remote dir
	if output, err := cmdutil.Run(receiver.dir, nil,
		"ssh",
		fmt.Sprintf("root@%v", node.IP),
		fmt.Sprintf("mkdir -p '%v'", receiver.NodeDir())); err != nil {
		return errors.Join(
			errors.New(output),
			err)
	}
	Log("INFO", node.IP, "node dir %v prepared", receiver.NodeDir())

	// rsync folders
	if output, err := cmdutil.Run(receiver.dir, nil, "rsync",
		"-avz",
		fmt.Sprintf("%v/", receiver.dir),
		fmt.Sprintf("root@%v:%v/", node.IP, receiver.NodeDir())); err != nil {
		return errors.Join(errors.New(output), err)
	}
	Log("INFO", node.IP, "node dir %v synced", receiver.NodeDir())

	// generate docker-compose content
	text, err := GenerateNodeDockerCompose(*receiver.config)
	if err != nil {
		return err
	}

	// create temp docker-compose.yaml
	tmpdir, err := os.MkdirTemp("", "load-hive-")
	if err != nil {
		return fmt.Errorf("failed to create temp dir, err: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			fmt.Printf("[WARN] failed to remove temp dir, err: %v", err)
		}
	}()
	if err := os.WriteFile(filepath.Join(tmpdir, "docker-compose.yaml"), []byte(text), 0750); err != nil {
		return fmt.Errorf("failed to write temp docker-compose.yaml, err: %v", err)
	}

	// copy docker-compose.yaml for node
	if output, err := cmdutil.Run(receiver.dir, nil, "rsync", "-avz",
		fmt.Sprintf("%v/docker-compose.yaml", tmpdir),
		fmt.Sprintf("root@%v:%v/docker-compose.yaml", node.IP, receiver.NodeDir())); err != nil {
		return errors.Join(errors.New(output), err)
	}
	Log("INFO", node.IP, "node dir %v docker-compose.yaml prepared", receiver.NodeDir())

	// remote docker-compose up -d
	if output, err := cmdutil.Run(receiver.dir, nil, "ssh",
		fmt.Sprintf("root@%v", node.IP),
		fmt.Sprintf("docker-compose --project-directory '%v' up -d --force-recreate --scale load-hive-%v-node=%v",
			receiver.NodeDir(),
			receiver.config.Hive.UUID,
			node.Worker,
		)); err != nil {
		return errors.Join(
			errors.New("failed to up node service"),
			errors.New(output),
			err)
	}
	Log("INFO", node.IP, "node dir %v docker-compose up -d success", receiver.NodeDir())

	return nil
}

func (receiver *Hive) UpNodes() error {
	var errLst []error
	var wg sync.WaitGroup
	for _, node := range receiver.config.Nodes {
		wg.Add(1)
		go func(node NodeConfig) {
			defer wg.Done()
			if err := receiver.UpNode(node); err != nil {
				errLst = append(errLst, err)
				return
			}
			Log("INFO", receiver.config.Hive.UUID, "node %v is up", node.IP)
		}(node)
	}
	wg.Wait()
	return errors.Join(errLst...)
}

func (receiver *Hive) Up() error {
	if err := receiver.UpMain(); err != nil {
		return fmt.Errorf("failed to up main, err: %v", err)
	}
	if err := receiver.UpNodes(); err != nil {
		return fmt.Errorf("failed to up nodes, err: %v", err)
	}
	Log("INFO", receiver.config.Hive.UUID, "all nodes are up")

	Log("INFO", receiver.config.Hive.UUID, "LOAD HIVE IS UP: http://%v:%v", receiver.config.Main.IP, receiver.config.Hive.Port)
	return nil
}
