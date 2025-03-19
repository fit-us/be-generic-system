package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
)

const bucketName = "fitus-file-bucket"

var storageClient *storage.Client

func init() {
	ctx := context.Background()
	var err error
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to create storage client: %v", err))
	}
}

func generateFileName() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000000))
}

func FileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error Parsing Form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := generateFileName()
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)

	go func() {
		if err := uploadFileToGCS(file, fileName); err != nil {
			fmt.Printf("Error uploading file to GCS: %v\n", err)
		} else {
			fmt.Println("File uploaded successfully to GCS!")
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File uploaded successfully to GCS!",
		"url":     url,
	})
}

func uploadFileToGCS(file io.Reader, fileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	wc := storageClient.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	defer wc.Close()

	if _, err := io.Copy(wc, file); err != nil {
		return err
	}

	return nil
}
