package zlp

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	RequestVariableMissing = "REQUEST_VARIABLE_MISSING"
	BadRequest             = "BAD_REQUEST"
	UserDeactivated        = "USER_DEACTIVATED"
	RealmDeactivated       = "REALM_DEACTIVATED"
	RateLimitHit           = "RATE_LIMIT_HIT"
)

const (
	Success = "success"
	Error   = "error"
)

const defaultUserAgent = "go-zlp"

/**
 * Common response data type
 *
 * Check the docs here: https://zulip.com/api/rest-error-handling
 */
type Response struct {
	// Always present
	Msg string `json:"msg"`
	// Will be either "success" or "error"
	Result string `json:"result"`

	// Sometimes present for common errors
	Code string `json:"code"`

	// Present if Code == RateLimitHit
	RetryAfter float32 `json:"retry-after"`

	// Present if Code == RequestVariableMissing
	VarName string `json:"var_name"`
}

func (r *Response) IsError() bool {
	return r.Result == Error
}

func (b *Bot) constructRequest(method, endpoint string, body *url.Values) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/%s/%s", b.Config.APIUrl, b.ApiVersion, endpoint)
	var bodyReader io.Reader = nil
	if body != nil {
		bodyReader = strings.NewReader(body.Encode())
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    userAgent := defaultUserAgent
    if len(b.Config.UserAgent) > 0 {
        userAgent = b.Config.UserAgent
    }
    req.Header.Set("User-Agent", userAgent)

	req.SetBasicAuth(b.Email, b.Config.APIKey)

	return req, nil
}

func (b *Bot) makeRequest(method, endpoint string, body *url.Values) (*http.Response, error) {
	req, err := b.constructRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	return b.Client.Do(req)
}

func (b *Bot) getResponseData(method, endpoint string, body *url.Values) ([]byte, error) {
	resp, err := b.makeRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (b *Bot) respToError(resp *http.Response) error {
	var codeType = resp.StatusCode / 100

	if codeType == 2 {
		// 2xx
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	// TODO: handle error from body
	fmt.Println(string(body))

	return fmt.Errorf(resp.Status)
}
