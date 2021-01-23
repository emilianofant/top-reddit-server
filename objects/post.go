package objects

// Post entity
type Post struct {
	ID               string `gorm:"primary_key" json:"id,omitempty"`
	Title            string `json:"title,omitempty"`
	Author           string `json:"author,omitempty"`
	EntryDate        string `json:"entrydate,omitempty"`
	Thumbnail        string `json:"thumbnail,omitempty"`
	NumberOfComments int    `json:"numberOfComments,omitempty"`
	Status           bool   `json:"status"`
}
