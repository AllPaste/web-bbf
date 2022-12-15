package dump

type Option func(o *options)

type options struct {
	request      bool
	response     bool
	body         bool
	headers      bool
	cookies      bool
	convertBytes func(dumpStr string)
}

func WithRequest(r bool) Option {
	return func(o *options) {
		o.request = r
	}
}

func WithResponse(r bool) Option {
	return func(o *options) {
		o.response = r
	}
}

func WithBody(b bool) Option {
	return func(o *options) {
		o.body = b
	}
}

func WithHeaders(h bool) Option {
	return func(o *options) {
		o.headers = h
	}
}

func WithCookies(c bool) Option {
	return func(o *options) {
		o.cookies = c
	}
}

func WithCB(f func(dumpStr string)) Option {
	return func(o *options) {
		o.convertBytes = f
	}
}
