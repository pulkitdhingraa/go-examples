package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileMetadata struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	UserID     string    `json:"user_id"`
	Size       float64   `json:"size"`
	Uploaded   time.Time `json:"uploaded"`
	Version    int       `json:"version"`
	SharedWith []string  `json:"shared_with"`
}

var (
	storageDir      = "./storage"
	metadataDir     = "./metadata"
	fileMetadataMap = make(map[string][]FileMetadata) // key string ->  value []FileMetadata
	mu              sync.Mutex
)

func init() {
	os.MkdirAll(storageDir, os.ModePerm)
	os.MkdirAll(metadataDir, os.ModePerm)
}

func saveMetadataToFile(metadata FileMetadata) error {
	filePath := filepath.Join(metadataDir, metadata.ID+".json")
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func loadMetadataFromFile(fileID string) (FileMetadata, error) {
	filePath := filepath.Join(metadataDir, fileID+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return FileMetadata{}, err
	}

	var metadata FileMetadata
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return FileMetadata{}, err
	}

	return metadata, nil
}

func saveFile(w http.ResponseWriter, r *http.Request) {
	// Get the file from header 
	userID := r.Header.Get("User-ID")
	if userID == "" {
		http.Error(w, "User-ID header missing", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Write the file to storage dir
	filename := fileHeader.Filename
	key := userID + ":" + filename

	mu.Lock()
	existingVersions := fileMetadataMap[key]
	currentVersion := len(existingVersions) + 1
	mu.Unlock()

	fileID := fmt.Sprintf("%s_v%d", filename, currentVersion) // file_v1
	filePath := filepath.Join(storageDir, fileID)
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	size, err := io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Save metadata to metadata dir
	metadata := FileMetadata{
		ID: fileID,
		Name: filename,
		UserID: userID,
		Size: float64(size),
		Uploaded: time.Now(),
		Version: currentVersion,
		SharedWith: []string{},
	}
	mu.Lock()
	fileMetadataMap[key] = append(fileMetadataMap[key], metadata)
	mu.Unlock()

	err = saveMetadataToFile(metadata)
	if err != nil {
		http.Error(w, "Error saving metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("File uploaded successfully with version %d", currentVersion)))
}

func getFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID")
	filename := r.URL.Query().Get("filename")
	if filename == "" || userID == "" {
		http.Error(w, "Filename or UserID query parameter missing", http.StatusBadRequest)
		return
	}

	key := userID + ":" + filename

	mu.Lock()
	defer mu.Unlock()

	versions, ok := fileMetadataMap[key]
	if !ok || len(versions) == 0 {
		http.Error(w, "File Not Found", http.StatusNotFound)
		return
	}

	latestMetadata := versions[len(versions)-1]
	filePath := filepath.Join(storageDir, latestMetadata.ID)

	if _,err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found on disk", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}

func shareFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID")
	filename := r.URL.Query().Get("filename")
	sharedWith := r.URL.Query().Get("sharedWith")
	if filename == "" || userID == ""  || sharedWith == "" {
		http.Error(w, "Filename or UserID or SharedWith query parameter missing", http.StatusBadRequest)
		return
	}

	key := userID + ":" + filename

	mu.Lock()
	defer mu.Unlock()

	versions, ok := fileMetadataMap[key]
	if !ok || len(versions) == 0 {
		http.Error(w, "File Not Found", http.StatusNotFound)
		return
	}

	latestMetadata := versions[len(versions)-1]
	latestMetadata.SharedWith = append(latestMetadata.SharedWith, sharedWith)

	err := saveMetadataToFile(latestMetadata)
	if err != nil {
		http.Error(w, "Error saving metadata", http.StatusInternalServerError)
		return
	}

	sharedKey := sharedWith + ":" + filename
	fileMetadataMap[sharedKey] = append(fileMetadataMap[sharedKey], latestMetadata)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File Shared Successfully"))
}

func main() {
	http.HandleFunc("/upload", saveFile)
	http.HandleFunc("/download", getFile)
	http.HandleFunc("/share", shareFile)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}