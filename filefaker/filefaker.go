package filefaker

type FileFaker struct {
	written []byte
	Err     error
	N       int // is this actually needed?
	Closed  bool
}

func (w *FileFaker) Write(p []byte) (n int, err error) {
	w.written = append(w.written, p...)
	n = len(p)
	w.N += n
	return n, err
}

func (w *FileFaker) Observe() string {
	return string(w.written)
}

func (w *FileFaker) Close() error {
	w.Closed = true
	return nil
}

func New() *FileFaker {
	f := new(FileFaker)
	return f
}
