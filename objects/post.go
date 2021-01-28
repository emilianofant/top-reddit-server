package objects

// Post entity
type Post struct {
	ID               string  `json:"id,omitempty"`
	Title            string  `json:"title,omitempty"`
	Author           string  `json:"author,omitempty"`
	EntryDate        float64 `json:"created_utc,omitempty"`
	Thumbnail        string  `json:"thumbnail"`
	PostURL          string  `json:"permalink,omitempty"`
	NumberOfComments int     `json:"num_comments,omitempty"`
	Viewed           bool    `json:"viewed"`
}
