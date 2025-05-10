package remote

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type SSHConfig struct {
	User     string
	Host     string
	Port     int
	Password string
}

type SSHClient struct {
	Config SSHConfig
}

func NewSSHClient(cfg SSHConfig) *SSHClient {
	return &SSHClient{Config: cfg}
}

func (c *SSHClient) StreamFilesPaginated(remoteDir string, pageSize int) (<-chan string, error) {
	outCh := make(chan string)
	cmd := exec.Command("sshpass", "-p", c.Config.Password,
		"ssh", "-p", fmt.Sprint(c.Config.Port),
		fmt.Sprintf("%s@%s", c.Config.User, c.Config.Host),
		fmt.Sprintf("find %s -type f \\( -name '*.csv' -o -name '*.pdf' -o -name '*.zip' \\) | sort", remoteDir),
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		buffer := make([]string, 0, pageSize)
		for scanner.Scan() {
			buffer = append(buffer, scanner.Text())
			if len(buffer) == pageSize {
				for _, file := range buffer {
					outCh <- file
				}
				buffer = buffer[:0]
			}
		}
		if len(buffer) > 0 {
			for _, file := range buffer {
				outCh <- file
			}
		}
		close(outCh)
		cmd.Wait()
	}()

	return outCh, nil
}

func (c *SSHClient) Download(remotePath, localPath string) error {
	os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	cmd := exec.Command("sshpass", "-p", c.Config.Password,
		"rsync", "-avz", "-e", fmt.Sprintf("ssh -p %d", c.Config.Port),
		fmt.Sprintf("%s@%s:%s", c.Config.User, c.Config.Host, remotePath),
		localPath,
	)
	return cmd.Run()
}
