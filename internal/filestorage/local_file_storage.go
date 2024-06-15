package filestorage

import (
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalFileStorage struct {
	storageFolderPath string
}

func (lfs *LocalFileStorage) BatchSaveFileContents(base64FileContents []string) ([]string, error) {
	var ids []string

	for _, base64FileContent := range base64FileContents {
		id, err := lfs.SaveFileContent(base64FileContent)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (lfs *LocalFileStorage) SaveFileContent(base64FileContent string) (string, error) {
	if err := lfs.ensureStorageFolderExists(); err != nil {
		panic("could not ensure storage folder exists")
	}

	decodedFileContent, err := base64.StdEncoding.DecodeString(base64FileContent)
	if err != nil {
		panic("could not decode base 64 string")
	}

	id := uuid.New().String()

	storageFilePath := filepath.Join(lfs.storageFolderPath, id)

	file, err := os.Create(storageFilePath)
	if err != nil {
		panic("error creating file")
	}
	defer file.Close()

	if _, err = file.WriteString(string(decodedFileContent)); err != nil {
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
