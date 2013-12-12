package golag

type GolagError struct {
  Msg string
  Err error
}

func (g *GolagError) Error() string {
  return g.Msg + "\n\t" + g.Err.Error()
}

func newError(msg string, err error) *GolagError {
  return &GolagError{Msg: msg, Err: err}
}