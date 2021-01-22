package objects

// Post entity
type Post struct {
	ID               string
	title            string
	author           string
	entryDate        string
	thumbnail        string
	numberOfComments int
	status           bool
}
