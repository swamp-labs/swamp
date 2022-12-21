package session

import (
	"fmt"
	"github.com/google/uuid"
)

type DataSource[T interface{}] interface {
	GetNextSessionData() T
}
type ExecutionTraces struct{}

type Session struct {
	Id           uuid.UUID
	Data         map[string]DataSource[any]
	Requirements []string
}

func (s *Session) New() {
	s.Id = uuid.New()
	fmt.Println(s.Id.String())
}

func (s *Session) GetContext() {
	//To be implemented to get context from previous dependency task
}
