package assertion

import "net/http"

// validateHeaders verify for each key []values inserted by user if
// the all the values exists for the associated key.
// Example : - Access-Control-Allow-Origin: ["*"]
// The function will find if the header exists with the associated value.
func (a *Assertion) validateHeaders(headers *http.Header) bool {
	valid := true
	// loop over table of maps in Assertion.Headers
	for _, m := range a.Headers {
		// loop over each map
		for headerName, headerValues := range m {
			// loop over the assertions values to execute contains function
			for _, headerValue := range headerValues {
				// We use the AND operation (&&). This way if a value define in assertion
				// is not present, the request is considered not valid
				valid = valid && contains(headers.Values(headerName), headerValue)
			}
		}
	}
	return valid
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
