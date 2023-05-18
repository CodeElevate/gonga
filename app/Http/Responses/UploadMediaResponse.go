package responses


type UploadMediaResponse struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Filename string `json:"filename"`
	Size     int  `json:"size"`
	ID       uint   `json:"id"`
	// Add more fields as needed
}