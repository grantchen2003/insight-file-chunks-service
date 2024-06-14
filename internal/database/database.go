package database

import "sync"

type FileChunk struct {
	UserId         string
	FilePath       string
	ChunkIndex     int
	NumTotalChunks int
}

type Database interface {
	BatchSave([]FileChunk) error
}

var (
	singletonInstance Database
	once              sync.Once
)

func GetInstance() Database {
	once.Do(func() {
		singletonInstance = &PostgreSql{}
	})

	return singletonInstance
}
