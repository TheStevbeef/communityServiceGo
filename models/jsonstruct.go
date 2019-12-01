package models

type Post struct {
	Post_ID   string `json:"id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	User      User   `json:"user,omitempty"`
	Message   string `json:"message,omitempty"`
	Media     Media  `json:"media,omitempty"`
}

type User struct {
	User_ID   string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Image_url string `json:"image_url,omitempty"`
}

type Media struct {
	Content_type string `json:"content_type,omitempty"`
	Url          string `json:"url,omitempty"`
}
