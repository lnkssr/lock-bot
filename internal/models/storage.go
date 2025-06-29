package models

type StorageResponse struct {
	Files   []string    `json:"files"`
	Storage StorageInfo `json:"storage"`
}

type StorageInfo struct {
	Limit int64 `json:"limit"`
	Used  int64 `json:"used"`
}
