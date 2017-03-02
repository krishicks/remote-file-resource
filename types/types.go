package types

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type Source struct {
	URI string `json:"uri"`
}

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
	Params  InParams
}

type InParams struct {
	Filename string `json:"filename"`
}

type Response []Version

type Version struct {
	ETag string `json:"etag"`
}
