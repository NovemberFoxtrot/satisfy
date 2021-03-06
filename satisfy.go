package satisfy

import (
	"net/http"
)

type Pipe struct {
	Handlers []http.Handler
	Fallback http.HandlerFunc
}

func New(handlers ...http.Handler) *Pipe {
	p := &Pipe{Handlers: handlers}
	p.Fallback = http.NotFound
	return p
}

func (p *Pipe) SetFallback(handler http.Handler) {
	p.Fallback = handler.ServeHTTP
}

func (p *Pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pipewriter := &pipeWriter{false, w}

	for _, handler := range p.Handlers {
		if handler == nil {
			continue
		}

		handler.ServeHTTP(pipewriter, r)

		if pipewriter.written {
			return
		}
	}

	if !pipewriter.written {
		p.Fallback(w, r)
	}
}

type pipeWriter struct {
	written bool
	http.ResponseWriter
}

func (w *pipeWriter) WriteHeader(status int) {
	w.written = true
	w.ResponseWriter.WriteHeader(status)
}
