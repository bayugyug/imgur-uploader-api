package main

// UploadRecords data holder on images uploaded
type UploadRecords struct {
	ID           string        `json:"id"`
	Created      string        `json:"created"`
	Finished     string        `json:"finished"`
	Status       string        `json:"status"`
	UploadedList *UploadedList `json:"uploaded"`
	URLS         []*URLInfo    `json:"-"`
}

// UploadedList data holder on images raw, GET /v1/images/upload/:jobId
type UploadedList struct {
	Complete []string `json:"complete"`
	Pending  []string `json:"pending"`
	Failed   []string `json:"failed"`
}

// ListImage data holder on all successfully pushed to ImGur, GET /v1/images
type ListImage struct {
	Uploaded []string `json:"uploaded"`
}

// ImgurResult response from the imgur api
type ImgurResult struct {
	Success bool      `json:"success,omitempty"`
	Status  int       `json:"status,omitempty"`
	Data    ImgurData `json:"data,omitempty"`
}

// ImgurData mapping of response, just not_all
type ImgurData struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Datetime    int    `json:"datetime,omitempty"`
	Type        string `json:"type,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	Size        int    `json:"size,omitempty"`
	Deletehash  string `json:"deletehash,omitempty"`
	Name        string `json:"name,omitempty"`
	Link        string `json:"link,omitempty"`
	Error       string `json:"error,omitempty"`
	Request     string `json:"request,omitempty"`
	Method      string `json:"method,omitempty"`
}

// URLInfo uri info
type URLInfo struct {
	Status    string `json:"status,omitempty"`
	Link      string `json:"link,omitempty"`
	ImgurLink string `json:"imgurLink,omitempty"`
	ImgurID   string `json:"imgurId,omitempty"`
}

// ParamImageURLS
type ParamImageURLS struct {
	URLS []string `json:"urls,omitempty"`
}

// APIResponse main data holders
type APIResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// UploadJobResponse
type UploadJobResponse struct {
	JobID string `json:"jobId,omitempty"`
}
