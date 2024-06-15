package filestorage

import (
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalFileStorage struct {
}

func (lfs *LocalFileStorage) SaveFile(base64FileContent string) (string, error) {
	decodedFileContent, err := base64.StdEncoding.DecodeString(base64FileContent)
	if err != nil {
		panic("could not decode base 64 string")
	}

	id := uuid.New().String()

	storageFilePath := filepath.Join("./internal/filestorage/localstorage", id)

	file, err := os.Create(storageFilePath)
	if err != nil {
		panic("error creating file")
	}
	defer file.Close()

	_, err = file.WriteString(string(decodedFileContent))
	if err != nil {
		panic("Error writing to file:")
	}

	return id, nil
}
