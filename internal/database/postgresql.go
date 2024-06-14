package database

import "fmt"

type PostgreSql struct {
}

func (psql *PostgreSql) BatchSave(fileChunks []FileChunk) error {
	fmt.Println("saving the following files")
	for _, fileChunk := range fileChunks {
		fmt.Println(fileChunk.FilePath)
	}
	return nil
}
