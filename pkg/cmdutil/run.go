package cmdutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Run(dir string, envMap map[string]string, name string, arg ...string) (string, error) {

	// prepare env
	env := os.Environ()
	for k, v := range envMap {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	var output bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdin = nil
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return output.String(), err
	}
	return output.String(), nil
}
