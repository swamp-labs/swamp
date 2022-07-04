package assertion

// validateCodeStatus verify if the returned http code
// matchs with at least one value provided by user
// In case user did not provide any code, we check code is 2XX
func (a assertion) validateCodeStatus(statusCode int) bool {

	if a.Code == nil {
		if statusCode > 199 && statusCode < 300 {
			return true
		}
	} else {
		for _, code := range a.Code {
			if code == statusCode {
				return true
			}
		}
	}
	return false
}
