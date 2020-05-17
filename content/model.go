package content

import "time"

type Model struct {
	ID        string
	CreatedAt time.Time
	View      string
	Model     string
	Title     string
	Content   interface{}
}
