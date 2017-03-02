package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/krishicks/remote-file-resource/types"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <destination>", os.Args[0])
	}

	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		log.Fatalf("error creating destination: %s", err)
	}

	var request types.InRequest
	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Fatalf("error reading request from stdin: %s\n", err)
	}

	client := http.Client{}
	resp, err := client.Get(request.Source.URI)
	if err != nil {
		log.Fatalf("error getting artifact at %s: %s\n", request.Source.URI, err)
	}
	defer resp.Body.Close()

	actual := resp.Header.Get("ETag")
	if actual != request.Version.ETag {
		log.Fatalf("error downloading artifact; version %s is no longer available", request.Version.ETag)
	}

	f, err := os.Create(filepath.Join(os.Args[1], request.Params.Filename))
	if err != nil {
		log.Fatalf("error creating file: %s", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatalf("error writing file: %s", err)
	}

	response := types.Response{
		{ETag: actual},
	}

	err = json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		log.Fatalf("error encoding response: %s\n", err)
	}
}
