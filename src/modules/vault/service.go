package vault

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"

	_ "modernc.org/sqlite"
)

type Service struct {
	db       *sql.DB
	httpAddr string
}

func NewService(dbPath string, httpAddr string) (*Service, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		data BLOB
	)`)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return &Service{db: db, httpAddr: httpAddr}, nil
}

func (s *Service) ProcessCommand(command string) string {
	switch strings.ToUpper(command) {
	case "STORE":
		return fmt.Sprintf("Please upload your file at: http://%s/upload", s.httpAddr)
	case "RETRIEVE":
		return s.getFileList()
	default:
		return "Unknown command. Use 'STORE' to upload a file or 'RETRIEVE' to list files."
	}
}

func (s *Service) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = s.db.Exec("INSERT INTO files (name, data) VALUES (?, ?)", header.Filename, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error storing file: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully", header.Filename)
}

func (s *Service) getFileList() string {
	rows, err := s.db.Query("SELECT name FROM files")
	if err != nil {
		return fmt.Sprintf("Error retrieving files: %v", err)
	}
	defer rows.Close()

	var files []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Sprintf("Error scanning row: %v", err)
		}
		files = append(files, name)
	}

	if len(files) == 0 {
		return "No files stored yet."
	}

	return fmt.Sprintf("Stored files:\n%s", strings.Join(files, "\n"))
}
