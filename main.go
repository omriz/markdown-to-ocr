package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogap/config"
	"github.com/gogap/go-pandoc/pandoc"
	_ "github.com/gogap/go-pandoc/pandoc/fetcher/data"
	"github.com/otiai10/gosseract/v2"
)

var pandocConf *config.Config

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	pandocConf = config.NewConfig(
		config.ConfigFile(filepath.Join(dir, "pandoc.conf")),
	)
	srcDir := flag.String("source", "", "Source Directory to scan")
	destDir := flag.String("dest", "", "Destination Directory to scan")
	flag.Parse()
	if err := convertFiles(*srcDir, *destDir); err != nil {
		log.Fatal(err)
	}
}

func parseBody(b io.ReadCloser) (string, error) {
	client := gosseract.NewClient()
	// We are reading the image from the body
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return "", err
	}
	client.SetImageFromBytes(body)
	return client.Text()
}

func parseMarkdown(md string) ([]byte, error) {
	pdoc, err := pandoc.New(pandocConf)
	if err != nil {
		return nil, err
	}
	fetcherOpts := pandoc.FetcherOptions{
		Name:   "data",
		Params: json.RawMessage(`{"data": "` + string(base64.StdEncoding.EncodeToString([]byte(md))+`"}`)),
	}
	convertOpts := pandoc.ConvertOptions{
		From: "markdown",
		To:   "docx",
	}
	return pdoc.Convert(fetcherOpts, convertOpts)
}

func convertFile(s, t string) error {
	ss, err := os.Open(s)
	if err != nil {
		return err
	}
	text, err := parseBody(ss)
	if err != nil {
		return err
	}
	tt, err := parseMarkdown(text)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(t, tt, 0644)
}

func convertFiles(src, dest string) error {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		err := os.Mkdir(dest, 0700)
		if err != nil {
			return err
		}
	}
	srcReadDir, err := os.Open(src)
	if err != nil {
		return err
	}
	srcFiles, err := srcReadDir.Readdir(0)
	if err != nil {
		return err
	}

	for _, f := range srcFiles {
		s := filepath.Join(src, f.Name())
		splitted := strings.Split(f.Name(), ".")
		t := filepath.Join(dest, strings.Join(splitted[:len(splitted)-1], ".")+".docx")
		log.Printf("Converting %s -> %s", s, t)
		err := convertFile(s, t)
		if err != nil {
			log.Printf("Failed to convert: %v", err)
		} else {
			_ = os.Remove(s)
		}
	}
	return nil
}
