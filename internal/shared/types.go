package shared

type Chapter struct {
	Number string
	Title  string
	Source string
}

type Exercise struct {
	Name        string
	Description string
	Done        bool
}

func (c Chapter) IsValid() bool {
	return c.Number != "" && c.Title != "" && c.Source != ""
}
