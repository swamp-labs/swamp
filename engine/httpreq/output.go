package httpreq

import (
	"bytes"
	"fmt"
	"github.com/swamp-labs/swamp/engine/logger"
	"io"
	"net/http"
	"strconv"
)

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, " { %s=\"%s\" }", key, value)
	}
	return b.String()
}

func (r *Request) displayResult(resp *http.Response, m map[string]string, b bool, s *Sample) {
	body, _ := io.ReadAll(resp.Body)
	logEntry := `--- Request ` + r.Name + " " + "(" + r.Method + ") " + r.URL.Format + "\n" +
		"--- Response status " + strconv.Itoa(resp.StatusCode) + " (valid:" + strconv.FormatBool(b) + ")" + " Returned variables: " + createKeyValuePairs(m) + " " +
		"Response body: " + string(body) + "\n" +
		"--- Traces " +
		fmt.Sprintf("%+v\n", s)
	logger.HttpLogger.Info(logEntry)
}
