package handlers

import (
	"context"
	"log"
	"sync"

	db "github.com/grantchen2003/insight/filechunks/internal/database"
	fs "github.com/grantchen2003/insight/filechunks/internal/filestorage"
	pb "github.com/grantchen2003/insight/filechunks/internal/protobufs"
)

type FileSaveStatus struct {
	isFullySaved     bool
	numTotalChunks   int
	chunksSaveStatus []bool
}

type FileChunkSaveStatus struct {
	filePath         string
	isLastSavedChunk bool
}

type FileChunksServiceHandler struct {
	pb.FileChunksServiceServer
	fileChunks map[string]map[string]FileSaveStatus
	mutex      sync.Mutex
}

func NewFileChunksServiceHandler() *FileChunksServiceHandler {
	return &FileChunksServiceHandler{
		fileChunks: make(map[string]map[string]FileSaveStatus),
		mutex:      sync.Mutex{},
	}
}

func (f *FileChunksServiceHandler) SaveFileChunks(ctx context.Context, req *pb.SaveFileChunksRequest) (*pb.SaveFileChunksResponse, error) {
	log.Println("received SaveFileChunks request")

	var fileChunks []db.FileChunk

	for _, fileChunk := range req.FileChunks {
		fileStorageId, err := fs.GetSingletonInstance().SaveFile(fileChunk.Content)
		if err != nil {
			panic("failed to save file")
		}

		fileChunks = append(fileChunks, db.FileChunk{
			UserId:         fileChunk.UserId,
			FilePath:       fileChunk.FilePath,
			ChunkIndex:     int(fileChunk.ChunkIndex),
			NumTotalChunks: int(fileChunk.NumTotalChunks),
			FileStorageId:  fileStorageId,
		})
	}

	if err := db.GetSingletonInstance().BatchSaveFileChunks(fileChunks); err != nil {
		panic("error batch-saving file chunks")
	}

	fileChunkSaveStatuses := f.getFileChunkSaveStatuses(fileChunks)

	var pbFileChunkSaveStatuses []*pb.FileChunkSaveStatus

	for _, fileChunkSaveStatus := range fileChunkSaveStatuses {
		pbFileChunkSaveStatuses = append(pbFileChunkSaveStatuses, &pb.FileChunkSaveStatus{
			FilePath:         fileChunkSaveStatus.filePath,
			IsLastSavedChunk: fileChunkSaveStatus.isLastSavedChunk,
		})
	}

	return &pb.SaveFileChunksResponse{FileChunkStatuses: pbFileChunkSaveStatuses}, nil
}

func (f *FileChunksServiceHandler) getFileChunkSaveStatuses(fileChunks []db.FileChunk) []FileChunkSaveStatus {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	statuses := make([]FileChunkSaveStatus, len(fileChunks))

	for i, fileChunk := range fileChunks {
		f.addFileChunk(fileChunk.UserId, fileChunk.FilePath, fileChunk.ChunkIndex, fileChunk.NumTotalChunks)

		isLastSavedChunk := f.isFullySaved(fileChunk.UserId, fileChunk.FilePath)

		statuses[i] = FileChunkSaveStatus{isLastSavedChunk: isLastSavedChunk, filePath: fileChunk.FilePath}

		f.cleanupFileChunks(fileChunk.UserId, fileChunk.FilePath)
	}

	return statuses
}

func (s *FileChunksServiceHandler) addFileChunk(userId string, filePath string, chunkIndex int, numTotalChunks int) {
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
