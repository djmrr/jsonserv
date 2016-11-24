package jsonserv

import (
	"io"
	"net/http"
	"net/url"
)

type mockBody struct {
	pos int
}

func (b *mockBody) Read(buff []byte) (int, error) {
	if b.pos == len(mockRequestBodyString) {
		return 0, io.EOF
	}
	n := copy(buff, []byte(mockRequestBodyString)[b.pos:])
	b.pos += n
	return n, nil
}

func (b *mockBody) Close() error {
	return nil
}

const mockRequestBodyString = `{"foo":"bar"}`

var exampleUrl, _ = url.Parse("http://example.com/foo?query=5")

func mockRequest() *http.Request {
	return &http.Request{
		URL:           exampleUrl,
		Method:        http.MethodGet,
		Body:          &mockBody{},
		ContentLength: int64(len(mockRequestBodyString)),
	}
}
