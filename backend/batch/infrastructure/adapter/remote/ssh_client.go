package remote

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/huydq/test/internal/pkg/config"
)

type SSHClient struct {
	Config SSHConfig
}

func NewSSHClient(cfg SSHConfig) SSHService {
	return &SSHClient{Config: cfg}
}

func buildTargetPaths(targetDate string) []string {
	appConfig := config.GetConfig()

	return []string{
		fmt.Sprintf("/%s/%s/", targetDate, appConfig.TransactionDetailsNoShippingRelatedPath),
		fmt.Sprintf("/%s/%s/", targetDate, appConfig.TransactionDetailsSummaryPath),
		fmt.Sprintf("/%s/%s/", targetDate, appConfig.TransactionDetailsShippingRelatedPath),
		fmt.Sprintf("/%s/", appConfig.TopUpDetailsPath),
		fmt.Sprintf("/%s/", appConfig.TopUpSummaryDetailsPath),
		fmt.Sprintf("/%s/", appConfig.TopUpReportPath),
		fmt.Sprintf("/%s/", appConfig.ValidInvoicesPath),
		fmt.Sprintf("/%s/", appConfig.ValidInvoicesDuplicatePath),
		fmt.Sprintf("/%s/", appConfig.ValidInvoicesSpreadsheetsPath),
	}
}

func (c *SSHClient) StreamFolderFilesPaginated(remoteDir string, pageSize int, targetDate string) (<-chan FileGroup, error) {
	outCh := make(chan FileGroup)

	targetedPaths := buildTargetPaths(targetDate)

	go func() {
		defer close(outCh)

		for _, folder := range targetedPaths {
			fullPath := filepath.Join(remoteDir, folder)
			log.Printf("[Stream] Searching in folder: %s", fullPath)
			cmdStr := fmt.Sprintf(
				"find %s -type f \\( -name '*.csv' -o -name '*.pdf' -o -name '*.zip' \\) -newermt '%s' | sort",
				fullPath,
				targetDate,
			)

			cmd := exec.Command("/usr/bin/sshpass", "-p", c.Config.Password,
				"ssh", "-o", "StrictHostKeyChecking=no",
				"-o", "UserKnownHostsFile=/dev/null",
				"-p", fmt.Sprint(c.Config.Port),
				fmt.Sprintf("%s@%s", c.Config.User, c.Config.Host),
				cmdStr,
			)

			stdout, _ := cmd.StdoutPipe()
			if err := cmd.Start(); err != nil {
				log.Printf("[Stream] Failed to exec folder %s: %v", folder, err)
				continue
			}

			scanner := bufio.NewScanner(stdout)
			buffer := make([]string, 0, pageSize)
			count := 0
			for scanner.Scan() {
				line := scanner.Text()
				buffer = append(buffer, line)
				count++
			}
			cmd.Wait()

			if count > 0 {
				log.Printf("[Stream] Folder %s: Found %d files", folder, count)
				outCh <- FileGroup{Folder: folder, Files: buffer}
			}
		}
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

	// Capture stderr for debugging
	stderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] Download failed for %s: %s", remotePath, string(stderr))
		return fmt.Errorf("rsync error: %w", err)
	}

	return nil
}

func (c *SSHClient) GetRemoteFileSize(remotePath string) (int64, error) {
	cmd := exec.Command("sshpass", "-p", c.Config.Password,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-p", fmt.Sprint(c.Config.Port),
		fmt.Sprintf("%s@%s", c.Config.User, c.Config.Host),
		fmt.Sprintf("stat -c %%s '%s'", remotePath),
	)
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get remote file size: %w", err)
	}
	var size int64
	_, err = fmt.Sscanf(string(out), "%d", &size)
	if err != nil {
		return 0, fmt.Errorf("failed to parse remote file size: %w", err)
	}
	return size, nil
}

func (c *SSHClient) ReadRemoteFile(remotePath string) ([]byte, error) {
	cmd := exec.Command("sshpass", "-p", c.Config.Password,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-p", fmt.Sprint(c.Config.Port),
		fmt.Sprintf("%s@%s", c.Config.User, c.Config.Host),
		fmt.Sprintf("cat '%s'", remotePath),
	)
	return cmd.Output()
}
