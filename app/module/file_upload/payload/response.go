package payload

import "time"

// FileUploadResp represents a response for a file upload
type FileUploadResp struct {
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}