package miggle

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func init() {
	http.Handle("/", Router())
}

func Router() *mux.Router {
	var res *Resource = NewResource()
	r := mux.NewRouter()

	r.HandleFunc("/test", res.RequestHandler(new(EndPoint)))
	r.HandleFunc("/get/status/{status:[0-9]+}", res.RequestHandler(new(EndPointStatus)))

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

}

type EndPointStatus struct{}

func (this *EndPointStatus) Init(request *http.Request) error {
	return nil
}

func (this *EndPointStatus) Get(request *http.Request) (int, interface{}, http.Header) {
	fmt.Println(request.RequestURI)
	vars := mux.Vars(request)
	status, err := strconv.Atoi(vars["status"])

	if err != nil {
		return 400, err, nil

	}
	return status, "data", nil
}

func TestAdvanceGet(t *testing.T) {
	statuses := []int{100, 200, 300, 400, 404, 500}
	for _, status := range statuses {
		r, err := http.NewRequest("GET", fmt.Sprintf("/get/status/%d", status), nil)
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
	}

	r, err := http.NewRequest("GET", "/get/status/string", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	Router().ServeHTTP(w, r)
	if w.Code != 404 {
		t.Error("Response status code does not equal")
	}

}
