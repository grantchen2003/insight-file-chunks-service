package filestorage

import (
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalFileStorage struct {
	storageFolderPath string
}

func NewLocalFileStorage() *LocalFileStorage {
	return &LocalFileStorage{
		storageFolderPath: "./internal/filestorage/localstorage",
	}
}

func (lfs *LocalFileStorage) GetFileContents(ids []string) ([][]byte, error) {
	var fileContents [][]byte

	for _, id := range ids {
		storageFilePath := filepath.Join(lfs.storageFolderPath, id)

		file, err := os.Open(storageFilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		fileContent, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		fileContents = append(fileContents, fileContent)
	}

	return fileContents, nil
}

func (lfs *LocalFileStorage) BatchSaveFileChunksContent(fileContents [][]byte) ([]string, error) {

	ids := make([]string, len(fileContents))

	for i, fileContent := range fileContents {
		id, err := lfs.SaveFileChunkContent(fileContent)
		if err != nil {
			return nil, err
		}

		ids[i] = id
	}

	return ids, nil
}

func (lfs *LocalFileStorage) SaveFileChunkContent(fileContent []byte) (string, error) {
	if err := lfs.ensureStorageFolderExists(); err != nil {
		panic("could not ensure storage folder exists")
	}

	id := uuid.New().String()

	storageFilePath := filepath.Join(lfs.storageFolderPath, id)

	file, err := os.Create(storageFilePath)
	if err != nil {
		panic("error creating file")
	}
	defer file.Close()

	if _, err = file.Write(fileContent); err != nil {
		panic("Error writing to file:")
	}

	return id, nil
}

func (lfs *LocalFileStorage) ensureStorageFolderExists() error {
	if _, err := os.Stat(lfs.storageFolderPath); os.IsExist(err) {
		return nil
	}

	return os.MkdirAll(lfs.storageFolderPath, 0755)
}
