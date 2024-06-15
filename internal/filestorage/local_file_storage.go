package filestorage

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type LocalFileStorage struct {
	storageFolderPath string
}

func (lfs *LocalFileStorage) BatchSaveFileContents(base64FileContents []string) ([]string, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	ids := make([]string, len(base64FileContents))
	errs := make(chan error, len(base64FileContents))

	for i, base64FileContent := range base64FileContents {
		wg.Add(1)
		go func(index int, content string) {
			defer wg.Done()
			id, err := lfs.SaveFileContent(content)
			if err != nil {
				errs <- err
				return
			}
			mu.Lock()
			ids[index] = id
			mu.Unlock()
		}(i, base64FileContent)
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return nil, <-errs
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
