package filestorage

import "sync"

type FileStorage interface {
	GetFileContents([]string) ([]string, error)
	BatchSaveFileContents([]string) ([]string, error)
	SaveFileContent(string) (string, error)
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
