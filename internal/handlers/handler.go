package handlers

import (
	"context"
	"log"

	db "github.com/grantchen2003/insight/filechunks/internal/database"
	fileStorage "github.com/grantchen2003/insight/filechunks/internal/filestorage"
	pb "github.com/grantchen2003/insight/filechunks/internal/protobufs"
	fcss "github.com/grantchen2003/insight/filechunks/internal/utils/filechunksavesync"
)

type FileChunksServiceHandler struct {
	pb.FileChunksServiceServer
}

func castToPbFileChunkSaveStatuses(fileChunkSaveStatuses []fcss.FileChunkSaveStatus) []*pb.FileChunkSaveStatus {
	var pbFileChunkSaveStatuses []*pb.FileChunkSaveStatus

	for _, fileChunkSaveStatus := range fileChunkSaveStatuses {
		pbFileChunkSaveStatuses = append(pbFileChunkSaveStatuses, &pb.FileChunkSaveStatus{
			FilePath:         fileChunkSaveStatus.FilePath,
			IsLastSavedChunk: fileChunkSaveStatus.IsLastSavedChunk,
		})
	}

	return pbFileChunkSaveStatuses
}

func batchSaveFileChunkPayloadContents(fileChunkPayloads []*pb.FileChunkPayload) ([]string, error) {
	var fileChunkPayloadContents []string

	for _, fileChunkPayload := range fileChunkPayloads {
		fileChunkPayloadContents = append(fileChunkPayloadContents, fileChunkPayload.Content)
	}

	return fileStorage.GetSingletonInstance().BatchSaveFileContents(fileChunkPayloadContents)
}

func getFileChunks(fileChunkPayloads []*pb.FileChunkPayload, fileStorageIds []string) []db.FileChunk {
	if len(fileChunkPayloads) != len(fileStorageIds) {
		panic("lengths do not match")
	}

	var fileChunks []db.FileChunk

	for i, fileChunkPayload := range fileChunkPayloads {
		fileChunks = append(fileChunks, db.FileChunk{
			UserId:         fileChunkPayload.UserId,
			FilePath:       fileChunkPayload.FilePath,
			ChunkIndex:     int(fileChunkPayload.ChunkIndex),
			NumTotalChunks: int(fileChunkPayload.NumTotalChunks),
			FileStorageId:  fileStorageIds[i],
		})
	}

	return fileChunks
}

func (f *FileChunksServiceHandler) SaveFileChunks(ctx context.Context, req *pb.SaveFileChunksRequest) (*pb.SaveFileChunksResponse, error) {
	log.Println("received SaveFileChunks request")

	fileStorageIds, err := batchSaveFileChunkPayloadContents(req.FileChunkPayloads)
	if err != nil {
		return nil, err
	}

	fileChunks := getFileChunks(req.FileChunkPayloads, fileStorageIds)

	err = db.GetSingletonInstance().BatchSaveFileChunks(fileChunks)
	if err != nil {
		return nil, err
	}

	fileChunkSaveStatuses := fcss.GetSingletonInstance().ReportFileChunkSaves(fileChunks)

	resp := &pb.SaveFileChunksResponse{
		FileChunkStatuses: castToPbFileChunkSaveStatuses(fileChunkSaveStatuses),
	}

	return resp, nil
}
