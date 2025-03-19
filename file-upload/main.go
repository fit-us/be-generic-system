package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
)
type Hello struct{
	Message string `json:"message"`
}

type UploadResponse struct {
    Message string `json:"message"`
    URL     string `json:"url"`
}

const (
	bucketName = "fitus-file-bucket"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse multipart form (최대 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error Parsing Form", http.StatusBadRequest)
		return
	}
	// "file" 필드에서 파일 추출
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// GCS 클라이언트 생성
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "Error creating storage client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// 업로드할 객체 생성 (파일명 사용)
	object := client.Bucket(bucketName).Object(header.Filename)
	wc := object.NewWriter(ctx)
	// 파일 데이터를 객체에 복사
	if _, err := io.Copy(wc, file); err != nil {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}
	if err := wc.Close(); err != nil {
		http.Error(w, "Error closing writer", http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, header.Filename)

    w.Header().Set("Content-Type", "application/json")
    response := UploadResponse{
        Message: "File uploaded successfully to GCS!",
        URL:     url,
    }
    json.NewEncoder(w).Encode(response)
}

// func FileUpload(w http.ResponseWriter, r *http.Request){
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	// Parse multipart form (최대 10MB)
// 	if err := r.ParseMultipartForm(10 << 20); err != nil {
// 		http.Error(w, "Error Parsing Form", http.StatusBadRequest)
// 		return
// 	}
// 	// "file" 필드에서 파일 추출
// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// GCS 클라이언트 생성
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
// 	defer cancel()
// 	client, err := storage.NewClient(ctx, option.WithCredentialsFile("fitus-file-bucket.json"))
// 	if err != nil {
// 		http.Error(w, "Error creating storage client", http.StatusInternalServerError)
// 		return
// 	}
// 	defer client.Close()

// 	// 업로드할 객체 생성 (파일명 사용)
// 	object := client.Bucket(bucketName).Object(header.Filename)
// 	wc := object.NewWriter(ctx)
// 	// 파일 데이터를 객체에 복사
// 	if _, err := io.Copy(wc, file); err != nil {
// 		http.Error(w, "Error uploading file", http.StatusInternalServerError)
// 		return
// 	}
// 	if err := wc.Close(); err != nil {
// 		http.Error(w, "Error closing writer", http.StatusInternalServerError)
// 		return
// 	}
// 	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, header.Filename)

//     w.Header().Set("Content-Type", "application/json")
//     response := UploadResponse{
//         Message: "File uploaded successfully to GCS!",
//         URL:     url,
//     }
//     json.NewEncoder(w).Encode(response)
// }