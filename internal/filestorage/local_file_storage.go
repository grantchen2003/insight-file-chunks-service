package filestorage

import (
	"encoding/base64"
	"io/ioutil"
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

func (lfs *LocalFileStorage) GetFileContents(ids []string) ([]string, error) {
	var fileContents []string

	for _, id := range ids {
		storageFilePath := filepath.Join(lfs.storageFolderPath, id)

		file, err := os.Open(storageFilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		fileContent, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		fileContents = append(fileContents, string(fileContent))
	}

	return fileContents, nil
}

func (lfs *LocalFileStorage) BatchSaveFileContents(base64FileContents []string) ([]string, error) {

	ids := make([]string, len(base64FileContents))

	for i, base64FileContent := range base64FileContents {
		id, err := lfs.SaveFileContent(base64FileContent)
		if err != nil {
			return nil, err
		}

		ids[i] = id
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
