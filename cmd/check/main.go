package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/krishicks/remote-file-resource/types"
)

func main() {
	var request types.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Fatalf("error reading request from stdin: %s\n", err)
	}

	resp, err := http.Head(request.Source.URI)
	if err != nil {
		log.Fatalf("error getting artifact at %s: %s\n", request.Source.URI, err)
	}

	actual := resp.Header.Get("ETag")

	var response types.CheckResponse
	if actual == request.Version.ETag {
		response = types.CheckResponse{
			{ETag: request.Version.ETag},
		}
	} else {
		response = types.CheckResponse{
			{ETag: request.Version.ETag},
			{ETag: actual},
		}
	}

	err = json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		log.Fatalf("error encoding response: %s\n", err)
	}
}
