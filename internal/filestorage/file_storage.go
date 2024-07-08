package filestorage

import "sync"

type FileChunkContent []byte

type FileStorage interface {
	GetFileContents([]string) ([]FileChunkContent, error)
	BatchSaveFileChunksContent([]FileChunkContent) ([]string, error)
	SaveFileChunkContent(FileChunkContent) (string, error)
	BatchDeleteFileChunksContent([]string) error
}

var (
	singletonInstance FileStorage
	once              sync.Once
)

func GetSingletonInstance() FileStorage {
	once.Do(func() {
		singletonInstance = NewLocalFileStorage()
	})

	return singletonInstance
}
