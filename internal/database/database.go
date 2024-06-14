package database

import "sync"

type FileChunk struct {
	UserId         string
	FilePath       string
	ChunkIndex     int
	NumTotalChunks int
}

type Database interface {
	BatchSaveFileChunks([]FileChunk) error
}

var (
	singletonInstance Database
	once              sync.Once
)

func GetSingletonInstance() Database {
	once.Do(func() {
		singletonInstance = &MongoDb{}
	})

	return singletonInstance
}
