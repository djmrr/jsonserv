package jsonserv

type View func(app interface{}, r *Request, out *Response)
