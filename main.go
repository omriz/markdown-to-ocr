package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	vision "cloud.google.com/go/vision/apiv1"
)

func main() {
	fmt.Println("Starting Server")
	http.HandleFunc("/", HandleOCR)
	http.ListenAndServe(":8080", nil)
}

func parseBody(ctx context.Context, b io.ReadCloser) (string, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return "", err
	}

	// We are reading the image from the body
	defer b.Close()
	image, err := vision.NewImageFromReader(b)
	if err != nil {
		return "", err
	}
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return "", err
	}

	if len(annotations) == 0 {
		return "", fmt.Errorf("No text found")
	}
	return annotations[0].Description, nil
}

// HandleOCR handles image requests
func HandleOCR(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Server only support POST requests")
		return
	}
	ctx := r.Context()
	// Getting the body of the request
	ss, err := parseBody(ctx, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error %v", err)))
		return
	}
	fmt.Fprint(w, ss)
}
