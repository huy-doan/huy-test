package remote

type SSHConfig struct {
	User     string
	Host     string
	Port     int
	Password string
}

type FileGroup struct {
	Folder string
	Files  []string
}

type SSHService interface {
	StreamFolderFilesPaginated(remoteDir string, pageSize int, targetDate string) (<-chan FileGroup, error)

	Download(remotePath, localPath string) error

	GetRemoteFileSize(remotePath string) (int64, error)

	ReadRemoteFile(remotePath string) ([]byte, error)
}
