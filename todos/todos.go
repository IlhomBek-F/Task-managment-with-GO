package todos

import (
	"errors"
	"time"
)

type Todo struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t Todo) Validate() error {
	var err error

	if len(t.Title) == 0 {
		err = errors.New("title can't be blank")
	}

	return err
}
