package main

import (
	"fmt"
	"net/http"

	function "github.com/fit-us/be-generic-system"
)

func main() {
	http.HandleFunc("/file-upload", function.FileUpload)
	fmt.Println("starting server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
