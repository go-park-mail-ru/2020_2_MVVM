package general

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func PerformRequest(r http.Handler, method, path string, bodyIn interface{}) (*httptest.ResponseRecorder, error) {
	var bodyReader io.Reader = nil
	if bodyIn != nil {
		bodyReader = strings.NewReader(getStrFromJson(bodyIn))
	}
	req, err := http.NewRequest(method, path, bodyReader)
	if err != nil {
		return nil, err
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, nil
}

func ResponseComparator(w httptest.ResponseRecorder, expectedCode int, expectedBodyIn interface{}) error {
	if w.Code != expectedCode {
		return errors.New(fmt.Sprintf("Expected to get status %d but instead got %d\n", expectedCode, w.Code))
	}
	expectedBody := getStrFromJson(expectedBodyIn)
	actualBody := w.Body.String()
	if expectedBody != actualBody {
		return errors.New(fmt.Sprintf("Expected to get response %s but instead got %s\n", actualBody, expectedBody))
	}
	return nil
}

func getStrFromJson(in interface{}) string {
	bArr, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	return string(bArr)
}

