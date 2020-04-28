package hcaptcha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type HCaptcha struct {
	Secret      string
	HCaptchaURL string
}

type Response struct {
	// Success indicates if the challenge was passed
	Success bool `json:"success"`
	// ChallengeTs is the timestamp of the captcha
	ChallengeTs time.Time `json:"challenge_ts"`
	// Hostname is the hostname of the passed captcha
	Hostname string `json:"hostname"`
	// Credit indicates  whether the response will be credited (optional)
	Credit bool `json:"credit"`
	// ErrorCodes contains error codes returned by hCaptcha (optional)
	ErrorCodes []string `json:"error-codes"`
}

func New(secret string) *HCaptcha {
	return &HCaptcha{
		Secret:      secret,
		HCaptchaURL: "https://hcaptcha.com/siteverify",
	}
}

// Verify verifies a "h-captcha-response" data field, with an optional remote IP set.
func (h *HCaptcha) Verify(response, remoteip string) (*Response, error) {
	values := url.Values{"secret": {h.Secret}, "response": {response}}
	if remoteip != "" {
		values.Set("remoteip", remoteip)
	}
	resp, err := http.PostForm(h.HCaptchaURL, values)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP read error: %w", err)
	}

	r := Response{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("JSON error: %w", err)
	}

	return &r, nil
}
