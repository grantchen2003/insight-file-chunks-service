package database

import "sync"

type FileChunk struct {
	RepositoryId   string
	FilePath       string
	ChunkIndex     int
	NumTotalChunks int
	FileStorageId  string
}

type FileChunkSaveStatus struct {
	FilePath         string
	IsLastSavedChunk bool
}

type Database interface {
	Connect() error
	Close() error
	GetSortedFileChunksFileStorageIds(string, string) ([]string, error)
	GetFileChunksFileStorageIdsByRepositoryId(string) ([]string, error)
	GetFileChunksFileStorageIdsByRepositoryIdAndFilePaths(string, []string) ([]string, error)
	SaveFileChunk(FileChunk) error
	BatchSaveFileChunks([]FileChunk) error
	DeleteFileChunksByRepositoryId(string) error
	ReportFileChunkSaves([]FileChunk) ([]FileChunkSaveStatus, error)
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
