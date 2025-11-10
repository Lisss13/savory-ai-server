package payload

// QRCodeResp represents a response for a QR code
type QRCodeResp struct {
	URL       string `json:"url"`
	ImageURL  string `json:"image_url"`
	TargetURL string `json:"target_url"`
}