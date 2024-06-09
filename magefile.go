//go:build mage
// +build mage

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"log/slog"

	ilog "github.com/go-zen-chu/ubuntu-k8s-playbook/internal/log"
)

var (
	logger = slog.New(ilog.NewPrettyHandler(os.Stdout, ilog.PrettyHandlerOptions{}))
)

func runCmd(cmd []string) (string, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	logger.Info("Running command", "cmd", c.String())
	out, err := c.CombinedOutput()
	return string(out), err
}

func runLongRunningCmd(cmd []string) (string, string, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	cmdStdoutBuffer := bytes.NewBufferString("")
	cmdStderrBuffer := bytes.NewBufferString("")
	cmdStdoutMultiWriter := io.MultiWriter(os.Stdout, cmdStdoutBuffer)
	cmdStderrMultiWriter := io.MultiWriter(os.Stderr, cmdStderrBuffer)
	logger.Info("Running long running command", "cmd", c.String())
	c.Stdout = cmdStdoutMultiWriter
	c.Stderr = cmdStderrMultiWriter
	err := c.Run()
	return cmdStdoutBuffer.String(), cmdStderrBuffer.String(), err
}

func runCmdWithStdin(cmd []string) (string, string, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	cmdStdoutBuffer := bytes.NewBufferString("")
	cmdStderrBuffer := bytes.NewBufferString("")
	cmdStdoutMultiWriter := io.MultiWriter(os.Stdout, cmdStdoutBuffer)
	cmdStderrMultiWriter := io.MultiWriter(os.Stderr, cmdStderrBuffer)
	logger.Info("Running command with stdin", "cmd", c.String())
	c.Stdout = cmdStdoutMultiWriter
	c.Stderr = cmdStderrMultiWriter
	stdin, err := c.StdinPipe()
	if err != nil {
		return "", "", fmt.Errorf("creating stdin pipe: %w", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// close stdin pipe for proceeding next
		defer stdin.Close()
		for scanner.Scan() {
			// get input when enter pressed
			input := scanner.Text()
			// logger.Info("input", "input", input)
			if _, err := fmt.Fprint(stdin, input); err != nil {
				logger.Error("writing to stdin", "error", err)
				break
			}
			break
		}
		if err := scanner.Err(); err != nil {
			logger.Error("scanning stdin", "error", err)
			return
		}
	}()
	err = c.Run()
	if err != nil {
		return cmdStdoutBuffer.String(), cmdStderrBuffer.String(), fmt.Errorf("running command: %w", err)
	}
	wg.Wait()
	return cmdStdoutBuffer.String(), cmdStderrBuffer.String(), nil
}

// TODO: must be more sophisticated
func splitCmd(cmd string) []string {
	return strings.Split(cmd, " ")
}

/*=======================
setup
=======================*/

// Download secret env file
func DownloadSecretEnvFile() error {
	stdout, stderr, err := runCmdWithStdin(splitCmd("gcloud auth login --no-launch-browser"))
	if err != nil {
		return fmt.Errorf("authorizing gcloud: %s\nerror log: %s\n%w", stdout, stderr, err)
	}
	cout, err := runCmd(splitCmd("gcloud secrets versions access latest --secret HOME_K8S_SECRET_ENV"))
	if err != nil {
		return fmt.Errorf("downloading secret env file: %s\n%w", cout, err)
	}
	err = ioutil.WriteFile(".env", []byte(cout), 0644)
	if err != nil {
		return fmt.Errorf("writing to .env file: %w", err)
	}
	return nil
}

/*=======================
workflow
=======================*/

// Generate ansible hosts.yaml
func GenerateHostsYaml() error {
	// envsubst read env vars and substitute them
	out, err := runCmd([]string{"bash", "-c", "envsubst < hosts-template.yml > hosts.yml"})
	if err != nil {
		return fmt.Errorf("error generating hosts.yaml: %w\nerror log: %s", err, out)
	}
	return nil
}
