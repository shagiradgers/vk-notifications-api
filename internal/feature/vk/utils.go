package vk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const apiUrl = "api.vk.ru"

type params map[string]any

func (c *client) makeRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return http.DefaultClient.Do(req)
}

func (c *client) getUrl(method string, params params) string {
	url := fmt.Sprintf("https://%s/method/%s?", apiUrl, method)

	for k, value := range params {
		if v, ok := value.(string); ok {
			value = strings.Join(strings.Split(v, " "), "+")
		}
		url += fmt.Sprintf("%s=%v&", k, value)
	}
	return url
}

func (c *client) getDefaultParams(extraParams params) params {
	p := params{"random_id": 0, "v": 5.131}
	for k, v := range extraParams {
		p[k] = v
	}
	return p
}
