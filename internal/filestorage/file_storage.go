package filestorage

import "sync"

type FileStorage interface {
	GetFileContents([]string) ([][]byte, error)
	BatchSaveFileChunksContent([][]byte) ([]string, error)
	SaveFileChunkContent([]byte) (string, error)
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
