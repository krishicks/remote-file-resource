package types

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type Source struct {
	URI string `json:"uri"`
}

type Response []Version

type Version struct {
	ETag string `json:"etag"`
}
