package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySql struct {
	db                          *sql.DB
	isConnected                 bool
	isFileChunksTableCreated    bool
	isFileChunksDatabaseCreated bool
	isFileChunksDatabaseUsed    bool
}

func (mySql *MySql) connect() error {
	var (
		username     = os.Getenv("MYSQL_USERNAME")
		password     = os.Getenv("MYSQL_PASSWORD")
		hostname     = os.Getenv("MYSQL_HOSTNAME")
		port         = os.Getenv("MYSQL_PORT")
		databaseName = os.Getenv("MYSQL_DATABASE_NAME")
	)

	dbSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, databaseName)

	db, err := sql.Open("mysql", dbSourceName)

	if err != nil {
		log.Panicln("Failed to connect to MySQL:", err)
		return err
	}

	fmt.Println("Connected to MySQL")

	mySql.db = db
	mySql.isConnected = true

	return nil
}

func (mySql *MySql) createFileChunksTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS file_chunks (
			id INT PRIMARY KEY AUTO_INCREMENT,
			repository_id VARCHAR(255),
			file_path VARCHAR(255),
			chunk_index INT,
			num_total_chunks INT
		);
	`

	_, err := mySql.db.Exec(query)

	mySql.isFileChunksTableCreated = true

	return err
}

func (mySql *MySql) createFileChunksDatabase() error {
	mySql.isFileChunksDatabaseCreated = true
	return nil
}

func (mySql *MySql) useFileChunksDatabase() error {
	mySql.isFileChunksDatabaseUsed = true
	return nil
}

func (mySql *MySql) BatchSaveFileChunks(fileChunks []FileChunk) error {
	if !mySql.isConnected {
		mySql.connect()
	}

	if !mySql.isFileChunksDatabaseCreated {
		mySql.createFileChunksDatabase()
	}

	if !mySql.isFileChunksDatabaseUsed {
		mySql.useFileChunksDatabase()
	}

	if !mySql.isFileChunksTableCreated {
		mySql.createFileChunksTable()
	}

	log.Println("here3")

	// Prepare the SQL query.
	query := "INSERT INTO file_chunks (repository_id, file_path, chunk_index, num_total_chunks) VALUES "
	placeholders := []string{}
	values := []interface{}{}

	for _, fileChunk := range fileChunks {
		placeholders = append(placeholders, "(?, ?, ?, ?)")
		values = append(values, fileChunk.RepositoryId, fileChunk.FilePath, fileChunk.ChunkIndex, fileChunk.NumTotalChunks)
	}

	query += fmt.Sprintf("%s", strings.Join(placeholders, ","))

	log.Println(query)

	// Execute the query.
	_, err := mySql.db.Exec(query, values...)

	return err
}
