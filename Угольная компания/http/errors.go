package http

import (
	"encoding/json"
	"time"
)

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func NewErrorDTO(message error) ErrorDTO {
	return ErrorDTO{
		Message: message.Error(),
		Time:    time.Now(),
	}
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
