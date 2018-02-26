package automerger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func URL(parts ...interface{}) string {
	out := make([]string, len(parts), len(parts))
	for i, part := range parts {
		out[i] = fmt.Sprint(part)
	}
	return strings.Join(out, "/")
}

func GithubRequest(method, url, token string, expectedCode int, body, result interface{}) []error {
	var data io.Reader
	var errs []error
	var err error
	if body != nil {
		var byteData []byte
		switch body.(type) {
		case []byte:
			byteData = body.([]byte)
		default:
			byteData, err = json.Marshal(body)
			if err != nil {
				errs = append(errs, err)
				return errs
			}
		}
		data = bytes.NewBuffer(byteData)
	}
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectedCode {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errs = append(errs, err)
		}
		errs = append(errs, errors.New(string(data)))
		return errs
	}
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			errs = append(errs, err)
			return errs
		}
	}

	return nil
}
