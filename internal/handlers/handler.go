package handlers

import (
	"context"
	"log"

	db "github.com/grantchen2003/insight/filechunks/internal/database"
	"github.com/grantchen2003/insight/filechunks/internal/filestorage"
	pb "github.com/grantchen2003/insight/filechunks/internal/protobufs"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FileChunksServiceHandler struct {
	pb.FileChunksServiceServer
}

func (f *FileChunksServiceHandler) GetSortedFileChunksContent(req *pb.GetSortedFileChunksContentRequest, stream pb.FileChunksService_GetSortedFileChunksContentServer) error {
	log.Println("received GetSortedFileChunksContent request")

	fileStorageIds, err := db.GetSingletonInstance().GetSortedFileChunksFileStorageIds(req.RepositoryId, req.FilePath)
	if err != nil {
		return err
	}

	fileChunksContent, err := filestorage.GetSingletonInstance().GetFileContents(fileStorageIds)
	if err != nil {
		return err
	}

	for _, fileChunkContent := range fileChunksContent {
		response := &pb.FileChunkContent{Content: fileChunkContent}
		if err := stream.Send(response); err != nil {
			return err
		}
	}

	return nil
}

func (f *FileChunksServiceHandler) CreateFileChunks(ctx context.Context, req *pb.CreateFileChunksRequest) (*pb.CreateFileChunksResponse, error) {
	log.Println("received CreateFileChunks request")

	fileStorageIds, err := batchSaveFileChunkPayloadContents(req.FileChunkPayloads)
	if err != nil {
		return nil, err
	}

	fileChunks := getFileChunks(req.FileChunkPayloads, fileStorageIds)

	database := db.GetSingletonInstance()

	err = database.BatchSaveFileChunks(fileChunks)
	if err != nil {
		return nil, err
	}

	fileChunkSaveStatuses, err := database.ReportFileChunkSaves(fileChunks)
	if err != nil {
		return nil, err
	}

	resp := &pb.CreateFileChunksResponse{
		FileChunkStatuses: castToPbFileChunkSaveStatuses(fileChunkSaveStatuses),
	}

	return resp, nil
}

func (f *FileChunksServiceHandler) DeleteFileChunksByRepositoryId(ctx context.Context, req *pb.DeleteFileChunksByRepositoryIdRequest) (*emptypb.Empty, error) {
	log.Println("received DeleteFileChunksByRepositoryId request")

	database := db.GetSingletonInstance()

	fileStorageIds, err := database.GetFileChunksFileStorageIdsByRepositoryId(req.RepositoryId)
	if err != nil {
		return nil, err
	}

	if err := filestorage.GetSingletonInstance().BatchDeleteFileChunksContent(fileStorageIds); err != nil {
		return nil, err
	}

	if err := database.DeleteFileChunksByRepositoryId(req.RepositoryId); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil

}

func (f *FileChunksServiceHandler) DeleteFileChunksByRepositoryIdAndFilePaths(ctx context.Context, req *pb.DeleteFileChunksByRepositoryIdAndFilePathsRequest) (*emptypb.Empty, error) {
	log.Println("received DeleteFileChunksByRepositoryIdAndFilePaths request")

	database := db.GetSingletonInstance()

	fileStorageIds, err := database.GetFileChunksFileStorageIdsByRepositoryIdAndFilePaths(req.RepositoryId, req.FilePaths)
	if err != nil {
		return nil, err
	}

	if err := filestorage.GetSingletonInstance().BatchDeleteFileChunksContent(fileStorageIds); err != nil {
		return nil, err
	}

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

func castToPbFileChunkSaveStatuses(fileChunkSaveStatuses []db.FileChunkSaveStatus) []*pb.FileChunkSaveStatus {
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
