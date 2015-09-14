package miggle

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func init() {
	http.Handle("/", Router())
}

func Router() *mux.Router {
	var res *Resource = NewResource()
	r := mux.NewRouter()

	r.HandleFunc("/test", res.RequestHandler(new(EndPoint)))
	r.HandleFunc("/status/{status:[0-9]+}", res.RequestHandler(new(EndPointAdvance)))
	return r
}

type EndPoint struct{}

func (this *EndPoint) Init(request *http.Request) error {
	return nil
}

func (this *EndPoint) Get(request *http.Request) (int, interface{}, http.Header) {
	return 200, "data", nil
}

func TestBasicGet(t *testing.T) {
	r, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	Router().ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Error("Response status code does not equal")
	}

	body, _ := ioutil.ReadAll(w.Body)

	if string(body) != "data" {
		t.Error("Response body does not equal")
	}

	if header, ok := w.HeaderMap["Content-Type"]; ok {
		if header[0] != "plain/text" {
			t.Error("Response header Content-Type: %s does not exists", header[0])
		}
	} else {
		t.Error("Response header Content-Type: does not exists", header[0])
	}
}

type EndPointAdvance struct {
	Status int
}

func (this *EndPointAdvance) Init(request *http.Request) error {
	vars := mux.Vars(request)
	if id, ok := vars["status"]; ok {
		status, err := strconv.Atoi(id)

		if err != nil {
			return err
		}
		this.Status = status
	}
	return nil
}

func (this *EndPointAdvance) Get(request *http.Request) (int, interface{}, http.Header) {
	return this.Status, "data", http.Header{"Content-Type": {"application/json; charset=utf-8"}}
}

func (this *EndPointAdvance) Post(request *http.Request) (int, interface{}, http.Header) {
	data, _ := ioutil.ReadAll(request.Body)
	return this.Status, data, http.Header{"Content-Type": {"application/json; charset=utf-8"}}
}

func TestAdvanceGet(t *testing.T) {
	statuses := []int{100, 200, 300, 400, 404, 500}
	for _, status := range statuses {
		r, err := http.NewRequest("GET", fmt.Sprintf("/status/%d", status), nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		Router().ServeHTTP(w, r)
		if w.Code != status {
			t.Error("Response status code does not equal")
		}
		body, _ := ioutil.ReadAll(w.Body)

		if string(body) != "data" {
			t.Error("Response body does not equal")
		}

		if header, ok := w.HeaderMap["Content-Type"]; ok {
			if header[0] != "application/json; charset=utf-8" {
				t.Error("Response header Content-Type: %s does not exists", header[0])
			}
		} else {
			t.Error("Response header Content-Type' does not exists", header)
		}
	}

	r, err := http.NewRequest("GET", "/status/not_a_number", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	Router().ServeHTTP(w, r)
	if w.Code != 404 {
		t.Error("Response status code should be 404")
	}

}

func TestAdvancePost(t *testing.T) {
	status := 201
	data := "data data"
	r, err := http.NewRequest("POST", fmt.Sprintf("/status/%d", status), strings.NewReader(data))
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	Router().ServeHTTP(w, r)
	if w.Code != status {
		t.Error("Response status code does not equal")
	}
	body, _ := ioutil.ReadAll(w.Body)

	if string(body) != data {
		t.Error("Response body does not equal")
	}

	if header, ok := w.HeaderMap["Content-Type"]; ok {
		if header[0] != "application/json; charset=utf-8" {
			t.Error("Response header Content-Type: %s does not exists", header[0])
		}
	} else {
		t.Error("Response header Content-Type' does not exists", header)
	}

}
