package types

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version,omitempty"`
}

type CheckResponse []Version

type Source struct {
	URI string `json:"uri"`
}

type InRequest struct {
	Source  Source   `json:"source"`
	Version Version  `json:"version,omitempty"`
	Params  InParams `json:"params,omitempty"`
}

type InParams struct {
	Filename string `json:"filename"`
}

type InResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Filename string `json:"filename"`
	URI      string `json:"uri"`
}

type Version struct {
	ETag string `json:"etag"`
}
