package main

import (
	"fmt"
	"net/http"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error %v", err)))
		return
	}

	//f, err := os.Open("/resources/text.png")
	f, err := os.Open("/resources/pre_ocr.jpg")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error %v", err)))
		return
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
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
	ss := ""
	for _, annotation := range annotations {
		ss += annotation.Description + "\n"
	}
	fmt.Fprint(w, ss)
}
