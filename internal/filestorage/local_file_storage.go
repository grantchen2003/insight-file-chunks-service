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

func (lfs *LocalFileStorage) GetFileContents(ids []string) ([]FileChunkContent, error) {
	var fileContents []FileChunkContent

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

func (lfs *LocalFileStorage) BatchSaveFileChunksContent(fileChunksContent []FileChunkContent) ([]string, error) {
	ids := make([]string, len(fileChunksContent))

	for i, fileContent := range fileChunksContent {
		id, err := lfs.SaveFileChunkContent(fileContent)
		if err != nil {
			return nil, err
		}

		ids[i] = id
	}

	return ids, nil
}

func (lfs *LocalFileStorage) SaveFileChunkContent(fileChunkContent FileChunkContent) (string, error) {
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

	if _, err = file.Write(fileChunkContent); err != nil {
		panic("Error writing to file:")
	}

	return id, nil
}

func (lfs *LocalFileStorage) BatchDeleteFileChunksContent(ids []string) error {
	fileNames := make(map[string]struct{})
	for _, id := range ids {
		fileNames[id] = struct{}{}
	}

	entries, err := os.ReadDir(lfs.storageFolderPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		_, exists := fileNames[entry.Name()]
		if !exists {
			continue
		}

		filePath := filepath.Join(lfs.storageFolderPath, entry.Name())
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	return nil
}

func (lfs *LocalFileStorage) ensureStorageFolderExists() error {
	if _, err := os.Stat(lfs.storageFolderPath); os.IsExist(err) {
		return nil
	}

	return os.MkdirAll(lfs.storageFolderPath, 0755)
}
