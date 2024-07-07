package utils

import (
	"sync"

	db "github.com/grantchen2003/insight/filechunks/internal/database"
)

type FileChunkSaveStatus struct {
	FilePath         string
	IsLastSavedChunk bool
}

type FileSaveStatus struct {
	isFullySaved     bool
	chunksSaveStatus []bool
}

type FileChunkSaveSync struct {
	mutex sync.Mutex
	data  map[string]map[string]FileSaveStatus
}

var (
	singletonInstance *FileChunkSaveSync
	once              sync.Once
)

func GetSingletonInstance() *FileChunkSaveSync {
	once.Do(func() {
		singletonInstance = &FileChunkSaveSync{
			data:  make(map[string]map[string]FileSaveStatus),
			mutex: sync.Mutex{},
		}
	})

	return singletonInstance
}

func (f *FileChunkSaveSync) ReportFileChunkSaves(fileChunks []db.FileChunk) []FileChunkSaveStatus {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	statuses := make([]FileChunkSaveStatus, len(fileChunks))

	for i, fileChunk := range fileChunks {
		f.addFileChunk(fileChunk.RepositoryId, fileChunk.FilePath, fileChunk.ChunkIndex, fileChunk.NumTotalChunks)

		isLastSavedChunk := f.isFullySaved(fileChunk.RepositoryId, fileChunk.FilePath)

		statuses[i] = FileChunkSaveStatus{IsLastSavedChunk: isLastSavedChunk, FilePath: fileChunk.FilePath}

		f.cleanupFileChunks(fileChunk.RepositoryId, fileChunk.FilePath)
	}

	return statuses
}

func (f *FileChunkSaveSync) addFileChunk(repositoryId string, filePath string, chunkIndex int, numTotalChunks int) {
	if _, repositoryIdExists := f.data[repositoryId]; !repositoryIdExists {
		f.data[repositoryId] = make(map[string]FileSaveStatus)
	}

	if _, filePathExists := f.data[repositoryId][filePath]; !filePathExists {
		f.data[repositoryId][filePath] = FileSaveStatus{isFullySaved: false, chunksSaveStatus: make([]bool, numTotalChunks)}
	}

	f.markChunkAsSaved(repositoryId, filePath, int(chunkIndex))
	f.updateIsFullySaved(repositoryId, filePath)
}

func (f *FileChunkSaveSync) markChunkAsSaved(repositoryId string, filePath string, chunkIndex int) {
	f.data[repositoryId][filePath].chunksSaveStatus[chunkIndex] = true
}

func (f *FileChunkSaveSync) isFullySaved(repositoryId string, filePath string) bool {
	return f.data[repositoryId][filePath].isFullySaved
}

func (f *FileChunkSaveSync) updateIsFullySaved(repositoryId string, filePath string) {
	for _, chunkIsSaved := range f.data[repositoryId][filePath].chunksSaveStatus {
		if !chunkIsSaved {
			break
		}
	}

	fileSaveStatus := f.data[repositoryId][filePath]
	fileSaveStatus.isFullySaved = true
	f.data[repositoryId][filePath] = fileSaveStatus
}

func (f *FileChunkSaveSync) cleanupFileChunks(repositoryId string, filePath string) {
	if f.isFullySaved(repositoryId, filePath) {
		delete(f.data[repositoryId], filePath)
	}

	if len(f.data[repositoryId]) == 0 {
		delete(f.data, repositoryId)
	}
}
