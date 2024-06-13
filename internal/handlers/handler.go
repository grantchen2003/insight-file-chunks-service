package handlers

import (
	"context"
	"sync"

	"github.com/grantchen2003/insight/filechunks/internal/protobufs"
)

type FileSaveStatus struct {
	isFullySaved     bool
	chunksSaveStatus []bool
}

type FileChunksServiceHandler struct {
	protobufs.FileChunksServiceServer
	fileChunks map[string]map[string]FileSaveStatus
	mutex      sync.Mutex
}

func NewFileChunksServiceHandler() *FileChunksServiceHandler {
	return &FileChunksServiceHandler{
		fileChunks: make(map[string]map[string]FileSaveStatus),
		mutex:      sync.Mutex{},
	}
}

func (f *FileChunksServiceHandler) SaveFileChunks(ctx context.Context, req *protobufs.FileChunks) (*protobufs.SaveFileChunksResponse, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	statuses := make([]*protobufs.FileChunkSaveStatus, len(req.FileChunks))

	for i, fileChunk := range req.FileChunks {
		f.addFileChunk(fileChunk.UserId, fileChunk.FilePath, fileChunk.ChunkIndex, fileChunk.NumTotalChunks)

		isLastSavedChunk := f.isFullySaved(fileChunk.UserId, fileChunk.FilePath)

		statuses[i] = &protobufs.FileChunkSaveStatus{IsLastSavedChunk: isLastSavedChunk, FilePath: fileChunk.FilePath}

		f.cleanupFileChunks(fileChunk.UserId, fileChunk.FilePath)
	}

	return &protobufs.SaveFileChunksResponse{FileChunkStatuses: statuses}, nil
}

func (s *FileChunksServiceHandler) addFileChunk(userId string, filePath string, chunkIndex int32, numTotalChunks int32) {
	if _, userIdExists := s.fileChunks[userId]; !userIdExists {
		s.fileChunks[userId] = make(map[string]FileSaveStatus)
	}

	if _, filePathExists := s.fileChunks[userId][filePath]; !filePathExists {
		s.fileChunks[userId][filePath] = FileSaveStatus{isFullySaved: false, chunksSaveStatus: make([]bool, numTotalChunks)}
	}

	s.markChunkAsSaved(userId, filePath, int(chunkIndex))
	s.updateIsFullySaved(userId, filePath)

}

func (s *FileChunksServiceHandler) markChunkAsSaved(userId string, filePath string, chunkIndex int) {
	s.fileChunks[userId][filePath].chunksSaveStatus[chunkIndex] = true
}

func (s *FileChunksServiceHandler) isFullySaved(userId string, filePath string) bool {
	return s.fileChunks[userId][filePath].isFullySaved
}

func (s *FileChunksServiceHandler) updateIsFullySaved(userId string, filePath string) {
	for _, chunkIsSaved := range s.fileChunks[userId][filePath].chunksSaveStatus {
		if !chunkIsSaved {
			break
		}
	}

	fileSaveStatus := s.fileChunks[userId][filePath]
	fileSaveStatus.isFullySaved = true
	s.fileChunks[userId][filePath] = fileSaveStatus
}

func (s *FileChunksServiceHandler) cleanupFileChunks(userId string, filePath string) {
	if s.isFullySaved(userId, filePath) {
		delete(s.fileChunks[userId], filePath)
	}

	if len(s.fileChunks[userId]) == 0 {
		delete(s.fileChunks, userId)
	}
}
