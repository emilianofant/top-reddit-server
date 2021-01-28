package objects

// Post entity
type Post struct {
	ID               string  `gorm:"primary_key" json:"id,omitempty"`
	Title            string  `json:"title,omitempty"`
	Author           string  `json:"author,omitempty"`
	EntryDate        float64 `json:"created_utc,omitempty"`
	Thumbnail        string  `json:"thumbnail"`
	PostURL          string  `json:"permalink,omitempty"`
	NumberOfComments int     `json:"num_comments,omitempty"`
	Status           bool    `json:"status"`
}
