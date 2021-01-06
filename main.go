package main

import (
	"fmt"
	"net/http"

	vision "cloud.google.com/go/vision/apiv1"
)

func main() {
	fmt.Println("Starting Server")
	http.HandleFunc("/", HandleOCR)
	http.ListenAndServe(":8080", nil)
}

// HandleOCR handles image requests
func HandleOCR(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Server only support POST requests")
		return
	}
	ctx := r.Context()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error %v", err)))
		return
	}

	// We are reading the image from the body
	defer r.Body.Close()
	image, err := vision.NewImageFromReader(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error %v", err)))
		return
	}
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error %v", err)))
		return
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No text found.")
		return
	}
	// The first annotation is the raw text
	ss := annotations[0].Description
	fmt.Fprint(w, ss)
}
