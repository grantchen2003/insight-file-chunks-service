package filestorage

import "sync"

type FileStorage interface {
	BatchSaveFileContents([]string) ([]string, error)
	SaveFileContent(string) (string, error)
}

var (
	singletonInstance FileStorage
	once              sync.Once
)

func GetSingletonInstance() FileStorage {
	once.Do(func() {
		singletonInstance = &LocalFileStorage{}
	})

	return singletonInstance
}
