package miggle

import (
	"encoding/json"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
	PATCH  = "PATCH"
)

type HttpContextHandler interface {
	Init(*http.Request) error
	NewRequest() HttpContextHandler
}

type GetImplemented interface {
	HttpContextHandler
	Get(*http.Request) (int, interface{}, http.Header)
}

type PostImplemented interface {
	HttpContextHandler
	Post(*http.Request) (int, interface{}, http.Header)
}

type PutImplemented interface {
	HttpContextHandler
	Put(*http.Request) (int, interface{}, http.Header)
}

type DeleteImplemented interface {
	HttpContextHandler
	Delete(*http.Request) (int, interface{}, http.Header)
}

type HeadImplemented interface {
	HttpContextHandler
	Head(*http.Request) (int, interface{}, http.Header)
}

type PatchImplemented interface {
	HttpContextHandler
	Patch(*http.Request) (int, interface{}, http.Header)
}

type Resource struct {
}

func NewResource() *Resource {
	return &Resource{}
}

func (api *Resource) RequestHandler(resource HttpContextHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var handler func(*http.Request) (int, interface{}, http.Header)
		var init func(*http.Request) error
		var context HttpContextHandler = resource.NewRequest()

		switch request.Method {
		case GET:
			if context, ok := context.(GetImplemented); ok {
				init = context.Init
				handler = context.Get
			}
		case POST:
			if context, ok := context.(PostImplemented); ok {
				init = context.Init
				handler = context.Post
			}
		case PUT:
			if context, ok := context.(PutImplemented); ok {
				init = context.Init
				handler = context.Put
			}
		case DELETE:
			if context, ok := context.(DeleteImplemented); ok {
				init = context.Init
				handler = context.Delete
			}
		case HEAD:
			if context, ok := context.(HeadImplemented); ok {
				init = context.Init
				handler = context.Head
			}
		case PATCH:
			if context, ok := context.(PatchImplemented); ok {
				init = context.Init
				handler = context.Patch
			}
		}

		if handler == nil {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := init(request)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return

		}
		code, data, headers := handler(request)

		if headers == nil {
			headers = http.Header{"Content-Type": {"plain/text"}}
		}

		var content []byte

		switch data.(type) {
		case string:
			content = []byte(data.(string))
		case []byte:
			content = data.([]byte)
		case error:
			rw.WriteHeader(http.StatusInternalServerError)
			return
		default:
			headers.Set("Content-Type", "application/json; charset=utf-8")

			var marshal_err error
			content, marshal_err = json.MarshalIndent(data, "", " ")
			if marshal_err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		for header, values := range headers {
			for _, value := range values {
				rw.Header().Add(header, value)
			}
		}
		rw.WriteHeader(code)
		rw.Write(content)
	}
}
