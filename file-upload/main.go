package function

import (
	"encoding/json"
	"net/http"
)
type Hello struct{
	Message string `json:"message"`
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	response:= Hello{Message: "Hello, World!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func FileUpload(w http.ResponseWriter, r *http.Request){
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	fmt.Fprint(w, "File uploaded successfully!")
// }