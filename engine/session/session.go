package session

import (
	"fmt"
	"github.com/google/uuid"
)

type DataSource[T interface{}] interface {
	GetNextSessionData() T
}
type SessionExecutionTrace struct {
	//	Traces map[string]httpreq.Trace
}

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

}
