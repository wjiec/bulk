package xio

import (
	"bytes"
	"context"
	"io"
)

type ProxyFunc func([]byte) ([]byte, error)

type reader struct {
	src   io.Reader
	ctx   context.Context
	buf   *bytes.Buffer
	proxy ProxyFunc
}

func (r *reader) Context() context.Context {
	return r.ctx
}

func (r *reader) WithContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *reader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		var bs []byte
		if n, err = r.src.Read(p); err == nil {
			if bs, err = r.proxy(p[:n]); err == nil {
				if _, err = r.buf.Write(bs); err == nil {
					n, err = r.buf.Read(p)
				}
			}
		}
	}
	return
}

func ProxyReader(src io.Reader, proxy ProxyFunc) io.Reader {
	return &reader{
		src:   src,
		proxy: proxy,
		buf:   &bytes.Buffer{},
		ctx:   context.Background(),
	}
}
