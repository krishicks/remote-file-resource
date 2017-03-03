package types

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version,omitempty"`
}

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

type Response []Version

type Version struct {
	ETag string `json:"etag"`
}
