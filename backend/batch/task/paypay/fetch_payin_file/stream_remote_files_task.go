package task

import (
	"log"
	"strings"

	remoteAdapter "github.com/huydq/test/batch/infrastructure/adapter/remote"
)

// RemoteFileInfo contains information about a remote file task
type RemoteFileInfo struct {
	RemotePath string
	LocalPath  string
	Folder     string
}

// StreamRemoteFilesTask handles streaming files from remote folders
type StreamRemoteFilesTask struct {
	SSHClient remoteAdapter.SSHService
	TargetDate string
}

// NewStreamRemoteFilesTask creates a new instance of StreamRemoteFilesTask
func NewStreamRemoteFilesTask(
	sshClient remoteAdapter.SSHService,
	targetDate string,
) *StreamRemoteFilesTask {
	return &StreamRemoteFilesTask{
		SSHClient: sshClient,
		TargetDate: targetDate,
	}
}

func (t *StreamRemoteFilesTask) Do(remoteDir string, pageSize int) (<-chan RemoteFileInfo, error) {
	// Create channel for remote file info
	fileInfoCh := make(chan RemoteFileInfo)
	
	// Get stream of file groups from SSH client
	stream, err := t.SSHClient.StreamFolderFilesPaginated(remoteDir, pageSize, t.TargetDate)
	if err != nil {
		close(fileInfoCh)
		return nil, err
	}
	
	// Process file groups in a goroutine
	go func() {
		defer close(fileInfoCh)
		
		// Process each folder
		for remoteFolder := range stream {
			if len(remoteFolder.Files) == 0 {
				continue
			}
			
			log.Printf("===== Starting to process folder: %s (%d files) =====", remoteFolder.Folder, len(remoteFolder.Files))
			
			// Send each file to the channel
			for _, file := range remoteFolder.Files {
				fileInfoCh <- RemoteFileInfo{
					RemotePath: file,
					LocalPath:  remoteDir,
					Folder:     remoteFolder.Folder,
				}
			}
			
			log.Printf("===== Finished processing folder: %s =====", remoteFolder.Folder)
		}
	}()
	
	return fileInfoCh, nil
}

// FilterFilesByExtension filters files by their extensions
func (t *StreamRemoteFilesTask) FilterFilesByExtension(path string, extensions []string) bool {
	lowerPath := strings.ToLower(path)
	for _, ext := range extensions {
		if strings.HasSuffix(lowerPath, ext) {
			return true
		}
	}
	return false
}
