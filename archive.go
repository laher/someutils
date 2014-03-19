package someutils

type ArchiveItem struct {
	//if FileSystemPath is empty, use Data instead
	FileSystemPath string
	ArchivePath    string
	Data           []byte
}
