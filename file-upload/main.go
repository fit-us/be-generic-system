package function

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"context"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const bucketName = "my-uploaded-files" // Cloud Storage 버킷 이름

func init() {
	// HTTP 트리거로 UploadFile 함수를 연결
	functions.HTTP("UploadFile", UploadFile)
}

// FileResponse는 JSON 응답 형식을 정의하는 구조체입니다.
type FileResponse struct {
	Message string `json:"message"`
	FileURL string `json:"fileUrl"`
}

// UploadFile 함수
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// POST 요청만 허용
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 파일 업로드를 위한 멀티파트 폼 데이터 파싱
	err := r.ParseMultipartForm(10 << 20) // 10MB 크기 제한
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// 'file' 필드에서 업로드된 파일 가져오기
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Cloud Storage 클라이언트 초기화
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "Failed to create storage client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// 파일을 GCS에 업로드
	objectName := "uploaded-file-" + filepath.Base(r.URL.Path) // 파일 이름 설정
	object := client.Bucket(bucketName).Object(objectName)
	writer := object.NewWriter(ctx)
	_, err = io.Copy(writer, file)
	if err != nil {
		http.Error(w, "Failed to write file to GCS", http.StatusInternalServerError)
		return
	}
	if err := writer.Close(); err != nil {
		http.Error(w, "Failed to close file writer", http.StatusInternalServerError)
		return
	}

	// 파일을 공개적으로 접근 가능하도록 설정
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		http.Error(w, "Failed to set ACL for file", http.StatusInternalServerError)
		return
	}

	// 파일의 공개 URL 생성
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	// JSON 응답 구조체 생성
	response := FileResponse{
		Message: "File uploaded successfully",
		FileURL: publicURL,
	}

	// JSON으로 응답 반환
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
