package log

type Engine interface {
	Write(r *Record) error
}

type Middleware interface {
	Handle(r *Record) (*Record, error)
}
