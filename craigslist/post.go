package craigslist

type Post struct {
	emails  []string
	numbers []string
}

func NewPost(body string) *Post { return &Post{} }

func extractEmails(body string)  {}
func extractNumbers(body string) {}
func showContactInfo()           {}
