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

type GetImplemented interface {
	Init(*http.Request) error
	Get(*http.Request) (int, interface{}, http.Header)
}

type Resource struct {
}

func NewResource() *Resource {
	return &Resource{}
}

func (api *Resource) RequestHandler(resource interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var handler func(*http.Request) (int, interface{}, http.Header)
		var init func(*http.Request) error

		switch request.Method {
		case GET:
			if resource, ok := resource.(GetImplemented); ok {
				init = resource.Init
				handler = resource.Get
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
