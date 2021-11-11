package mbsGo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/0sax/err2"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

func (c *Client) standardRequest(method, endpoint string, body interface{}, response interface{}) (err error) {

	fmt.Printf("\nBaseUrl:\n%v", c.baseUrl)
	fmt.Printf("\nEndpoint:\n%v", endpoint)
	fmt.Printf("\nMethod:\n%v", method)

	var pl io.Reader

	if body != nil {
		pl, err = preparePayload(body)
		if err != nil {
			return
		}
	}

	headers := map[string]string{
		"Client-ID":     c.clientId,
		"Client-Secret": c.clientSecret,
		"Content-Type":  "application/json",
	}

	err = makeRequest(method, c.baseUrl+endpoint, pl, headers, response)
	return

}

func preparePayload(body interface{}) (io.Reader, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nBody:\n%+v", string(b))
	return bytes.NewReader(b), nil
}

func makeRequest(
	method, url string, body io.Reader,
	headers map[string]string, responseTarget interface{}) error {
	if reflect.TypeOf(responseTarget).Kind() != reflect.Ptr {
		return errors.New("mbsGo: responseTarget must be a pointer to a struct for JSON unmarshalling")
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		err2.LogErr1("makeRequest", "NewRequest",
			err)
		return err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		err2.LogErr1("makeRequest", "Do",
			err)
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err2.LogErr1("makeRequest", "ReadAll",
			err)
		return err
	}

	fmt.Printf("\nResponse Body:\n%v", string(b))

	err = json.Unmarshal(b, responseTarget)
	if err != nil {
		err2.LogErr1("makeRequest", "Unmarshal",
			err)
		return err
	}

	return err
}
