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
		f.addFileChunk(fileChunk.UserId, fileChunk.FilePath, fileChunk.ChunkIndex, fileChunk.NumTotalChunks)

		isLastSavedChunk := f.isFullySaved(fileChunk.UserId, fileChunk.FilePath)

		statuses[i] = FileChunkSaveStatus{IsLastSavedChunk: isLastSavedChunk, FilePath: fileChunk.FilePath}

		f.cleanupFileChunks(fileChunk.UserId, fileChunk.FilePath)
	}

	return statuses
}

func (f *FileChunkSaveSync) addFileChunk(userId string, filePath string, chunkIndex int, numTotalChunks int) {
	if _, userIdExists := f.data[userId]; !userIdExists {
		f.data[userId] = make(map[string]FileSaveStatus)
	}

	if _, filePathExists := f.data[userId][filePath]; !filePathExists {
		f.data[userId][filePath] = FileSaveStatus{isFullySaved: false, chunksSaveStatus: make([]bool, numTotalChunks)}
	}

	f.markChunkAsSaved(userId, filePath, int(chunkIndex))
	f.updateIsFullySaved(userId, filePath)
}

func (f *FileChunkSaveSync) markChunkAsSaved(userId string, filePath string, chunkIndex int) {
	f.data[userId][filePath].chunksSaveStatus[chunkIndex] = true
}

func (f *FileChunkSaveSync) isFullySaved(userId string, filePath string) bool {
	return f.data[userId][filePath].isFullySaved
}

func (f *FileChunkSaveSync) updateIsFullySaved(userId string, filePath string) {
	for _, chunkIsSaved := range f.data[userId][filePath].chunksSaveStatus {
		if !chunkIsSaved {
			break
		}
	}

	fileSaveStatus := f.data[userId][filePath]
	fileSaveStatus.isFullySaved = true
	f.data[userId][filePath] = fileSaveStatus
}

func (f *FileChunkSaveSync) cleanupFileChunks(userId string, filePath string) {
	if f.isFullySaved(userId, filePath) {
		delete(f.data[userId], filePath)
	}

	if len(f.data[userId]) == 0 {
		delete(f.data, userId)
	}
}
