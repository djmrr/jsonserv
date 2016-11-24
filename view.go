package jsonserv

type View func(ctx interface{}, r *Request, out *Response)
