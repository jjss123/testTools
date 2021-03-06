package http_helpers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testTools/src/utils/clog"
	"testTools/src/utils/strings"
	"testing"
)

// Creates a new file upload http request with optional extra params
func NewFileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}

func GetIDFromRequest(req *http.Request) string {
	var requestId string
	requestId = req.Header.Get("requestId")
	if strings.IsEmptyString(requestId) {
		requestId = uuid.NewV1().String()
		req.Header.Set("requestId", requestId)
	}

	return requestId
}

func PrintRequestInfo(req *http.Request) {
	GetIDFromRequest(req)
	clog.Hlog.SetRequest(req).Debugf("request detail: Header: [%+v]", req.Header)
}

func GetVarFromRequest(key string, req *http.Request) string {
	strVar := ""
	context := mux.Vars(req)
	if context != nil {
		strVar, _ = context[key]
	}
	return strVar
}

func MakeQueryString(queries map[string]string) string {
	count := 0
	str := ""
	for k, v := range queries {
		if count == 0 {
			str = k + "=" + v
		} else {
			str = "&" + k + "=" + v
		}
		count++
	}
	return str
}

func MockRequest(
	t *testing.T,
	method string,
	path string,
	queries map[string]string,
	header map[string][]string,
	bodyItem interface{}) *http.Request {

	//setup body
	data, err := json.Marshal(bodyItem)
	assert.Nil(t, err)
	var body = &bytes.Buffer{}
	body.Write(data)

	req, err := http.NewRequest(method, path+"?"+MakeQueryString(queries), body)
	assert.Nil(t, err)
	req.Header = header

	return req
}


func MockBodyRequest(t *testing.T, bodyItem interface{}) *http.Request {
	return MockRequest(t, http.MethodPost, "dummy", nil, nil, bodyItem)
}

func MockPathRequest(t *testing.T, path string) *http.Request {
	return MockRequest(t, http.MethodGet, path, nil, nil, nil)
}

func MockQueryRequest(t *testing.T, queries map[string]string) *http.Request {
	return MockRequest(t, http.MethodGet, "dummy", queries, nil, nil)
}
