package handlers

import (
	"context"
	"log"

	db "github.com/grantchen2003/insight/filechunks/internal/database"
	"github.com/grantchen2003/insight/filechunks/internal/filestorage"
	pb "github.com/grantchen2003/insight/filechunks/internal/protobufs"
	fcss "github.com/grantchen2003/insight/filechunks/internal/utils/filechunksavesync"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FileChunksServiceHandler struct {
	pb.FileChunksServiceServer
}

func (f *FileChunksServiceHandler) GetSortedFileChunksContent(ctx context.Context, req *pb.GetSortedFileChunksContentRequest) (*pb.GetSortedFileChunksContentResponse, error) {
	log.Println("received GetSortedFileChunksContent request")

	fileStorageIds, err := db.GetSingletonInstance().GetSortedFileChunksFileStorageIds(req.RepositoryId, req.FilePath)
	if err != nil {
		return nil, err
	}

	fileChunksContent, err := filestorage.GetSingletonInstance().GetFileContents(fileStorageIds)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetSortedFileChunksContentResponse{
		FileChunksContent: castToPbFileChunksContent(fileChunksContent),
	}

	return resp, nil
}

func (f *FileChunksServiceHandler) CreateFileChunks(ctx context.Context, req *pb.CreateFileChunksRequest) (*pb.CreateFileChunksResponse, error) {
	log.Println("received CreateFileChunks request")

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

	resp := &pb.CreateFileChunksResponse{
		FileChunkStatuses: castToPbFileChunkSaveStatuses(fileChunkSaveStatuses),
	}

	return resp, nil
}

func (f *FileChunksServiceHandler) DeleteFileChunksByRepositoryId(ctx context.Context, req *pb.DeleteFileChunksByRepositoryIdRequest) (*emptypb.Empty, error) {
	log.Println("received DeleteFileChunksByRepositoryId request")

	database := db.GetSingletonInstance()

	if err := database.DeleteFileChunksByRepositoryId(req.RepositoryId); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil

}

func castToPbFileChunksContent(fileChunksContent []filestorage.FileChunkContent) []*pb.FileChunkContent {
	var pbFileChunksContent []*pb.FileChunkContent

	for _, fileChunkContent := range fileChunksContent {
		pbFileChunksContent = append(pbFileChunksContent, &pb.FileChunkContent{
			Content: fileChunkContent,
		})
	}

	return pbFileChunksContent
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
	var fileChunkPayloadContents []filestorage.FileChunkContent

	for _, fileChunkPayload := range fileChunkPayloads {
		fileChunkPayloadContents = append(fileChunkPayloadContents, fileChunkPayload.Content)
	}

	return filestorage.GetSingletonInstance().BatchSaveFileChunksContent(fileChunkPayloadContents)
}

func getFileChunks(fileChunkPayloads []*pb.FileChunkPayload, fileStorageIds []string) []db.FileChunk {
	if len(fileChunkPayloads) != len(fileStorageIds) {
		panic("lengths do not match")
	}

	var fileChunks []db.FileChunk

	for i, fileChunkPayload := range fileChunkPayloads {
		fileChunks = append(fileChunks, db.FileChunk{
			RepositoryId:   fileChunkPayload.RepositoryId,
			FilePath:       fileChunkPayload.FilePath,
			ChunkIndex:     int(fileChunkPayload.ChunkIndex),
			NumTotalChunks: int(fileChunkPayload.NumTotalChunks),
			FileStorageId:  fileStorageIds[i],
		})
	}

	return fileChunks
}
