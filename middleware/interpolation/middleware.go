package interpolation

import (
	"github.com/anticrew/log"
)

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

func (i *Middleware) Handle(rec *log.Record) (*log.Record, error) {
	r := newReplacer(rec.Message, rec.Attrs)
	defer r.dispose()

	msg, err := r.replace()
	if err != nil {
		return rec, err
	}

	rec.Message = msg
	return rec, nil
}
