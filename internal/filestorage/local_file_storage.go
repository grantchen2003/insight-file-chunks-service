package filestorage

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalFileStorage struct {
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
	decodedFileContent, err := base64.StdEncoding.DecodeString(base64FileContent)
	if err != nil {
		panic("could not decode base 64 string")
	}

	id := uuid.New().String()

	storageFolderPath := "./internal/filestorage/localstorage"

	if _, err := os.Stat(storageFolderPath); os.IsNotExist(err) {
		// Create the folder and any necessary parents
		if err := os.MkdirAll(storageFolderPath, 0755); err != nil {
			fmt.Println("Error creating folder:", err)
		}
	}

	storageFilePath := filepath.Join(storageFolderPath, id)

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
