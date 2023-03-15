package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	zerror "github.com/zmicro-team/zmicro/core/errors"
)

type ErrorReply struct {
	Code   int
	Body   []byte
	Header http.Header
}

func (e *ErrorReply) Error() string {
	return fmt.Sprintf("Invoke: Status Code: %d, Status Text: %s", e.Code, http.StatusText(e.Code))
}

func IntoErrno(err error) error {
	if err == nil {
		return nil
	}

	e := &ErrorReply{}
	ok := errors.As(err, &e)
	if !ok {
		return err
	}
	e1 := new(zerror.Error)
	if e2 := json.Unmarshal([]byte(e.Body), e1); e2 != nil {
		return err
	}
	return e1
}
