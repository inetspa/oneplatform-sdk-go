package oneid

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/inetspa/oneplatform-sdk-go/requests"
)

func checkError(resp requests.Response, err error) error {
	// connection error
	if err != nil {
		return errors.New("service unavailable")
	}

	// request ok
	if 200 <= resp.Code && resp.Code < 300 {
		return nil
	}

	// request error
	var r interface{}
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return errors.New("internal server error")
	}
	errStr := fmt.Sprintf("error(%d) %v", resp.Code, r.(map[string]interface{})["errorMessage"])
	return errors.New(errStr)
}
