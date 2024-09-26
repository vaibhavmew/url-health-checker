package service

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

/* -------------------------------------------------------------------------- */
/*                                  Interface                                 */
/* -------------------------------------------------------------------------- */
type Interface interface {
	Get(url string, key string) ([]byte, int, error)
	Up(url string)   //only used by mock
	Down(url string) //only used by mock
}

type client struct {
	client *http.Client
}

func NewClient() *client {
	return &client{
		client: http.DefaultClient,
	}
}

func (c *client) Get(url string, key string) ([]byte, int, error) {
	var (
		data []byte
		code int
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return data, code, err
	}

	req.Header.Add("api-key", key)

	res, err := c.client.Do(req)
	if err != nil {
		return data, code, err
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return data, code, err
	}

	return data, res.StatusCode, err
}

func (c *client) Up(url string) {
}

func (c *client) Down(url string) {
}

type mockclient struct {
	mu     sync.Mutex
	client map[string]struct{}
}

func NewMockClient() *mockclient {
	client := make(map[string]struct{})

	client["url_zero"] = struct{}{}
	client["url_one"] = struct{}{}
	client["url_two"] = struct{}{}
	client["url_three"] = struct{}{}

	return &mockclient{
		client: client,
	}
}

func (m *mockclient) Get(url string, key string) ([]byte, int, error) {
	var (
		data []byte
		err  error
	)

	//logic for mocking
	_, ok := m.client[url]
	if !ok {
		return data, http.StatusInternalServerError, err
	}

	data, _ = json.Marshal(&Resp{
		URL:  url,
		Name: "mock",
	})

	return data, http.StatusOK, err
}

func (m *mockclient) Up(url string) {
	m.mu.Lock()
	m.client[url] = struct{}{}
	m.mu.Unlock()
}

func (m *mockclient) Down(url string) {
	m.mu.Lock()
	delete(m.client, url)
	m.mu.Unlock()
}
