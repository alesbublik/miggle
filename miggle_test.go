package miggle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type EndPoint struct{}

func (this *EndPoint) Init(request *http.Request) error {
	return nil
}

func (this *EndPoint) Get(request *http.Request) (int, interface{}, http.Header) {
	return 200, "data", nil
}

func TestBasicGet(t *testing.T) {
	var res *Resource = NewResource()
	ep := new(EndPoint)

	mux := http.NewServeMux()
	mux.HandleFunc("/test", res.RequestHandler(ep))
	portString := fmt.Sprintf(":%d", 1234)
	go http.ListenAndServe(portString, mux)

	resp, err := http.Get("http://localhost:1234/test")
	if err != nil {
		t.Error(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "{}" {
		t.Error("Not Equal")
	}

}
