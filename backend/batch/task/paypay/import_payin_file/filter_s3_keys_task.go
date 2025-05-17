package task

import (
	"fmt"
	"log"
	"strings"
)

type FilterS3KeysTask struct {
	TargetFolders map[string]struct{}
}

// NewFilterS3KeysTask creates a new instance of FilterS3KeysTask
func NewFilterS3KeysTask(targetFolders map[string]struct{}) *FilterS3KeysTask {
	return &FilterS3KeysTask{
		TargetFolders: targetFolders,
	}
}

/**
* Do filters S3 keys based on the target folders and file extension.
* It checks if the key is in one of the target folders and if it has a .zip extension.
* If both conditions are met, it returns true; otherwise, it returns false.
*
* @param s3Key The S3 key to filter.
* @return bool True if the key should be processed, false otherwise.
*/
func (t *FilterS3KeysTask) Do(s3Key string) bool {
	// Check if key is in target folders
	folderMatch := ""
	fmt.Printf("[DEBUG] Checking S3 key: %s\n", s3Key)
	for folder := range t.TargetFolders {
		if strings.HasPrefix(s3Key, folder) {
			folderMatch = folder
			break
		}
	}
	
	if folderMatch == "" {
		return false // Skip files not in target folders
	}
	
	// Check if file has .zip extension
	if !strings.HasSuffix(strings.ToLower(s3Key), ".zip") {
		return false // Skip non-zip files
	}
	
	log.Printf("[DEBUG] S3 file matched for import: %s (in folder: %s)", s3Key, folderMatch)
	return true
}

func (t *FilterS3KeysTask) BuildTargetFolders(remoteDir, topUpSummaryDetailsPath, topUpReportPath string) map[string]struct{} {
	joinRemotePath := func(base, sub string) string {
		sub = strings.TrimPrefix(sub, "/")
		return strings.TrimRight(base, "/") + "/" + sub
	}
	
	return map[string]struct{}{
		strings.TrimLeft(joinRemotePath(remoteDir, topUpSummaryDetailsPath), "/"): {},
		strings.TrimLeft(joinRemotePath(remoteDir, topUpReportPath), "/"):         {},
	}
}
