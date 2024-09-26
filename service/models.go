package service

import (
	"sync"
)

type Dialer struct {
	mu      sync.Mutex
	current int
	urls    urls
	Client  Interface
	close   map[int]struct{}
}

type urls []url

type url struct {
	url   string
	key   string
	ch    chan struct{}
	close chan struct{}
}

/* -------------------------------------------------------------------------- */
/*                                   Request                                  */
/* -------------------------------------------------------------------------- */
type Req struct {
	In string `json:"in"`
}

/* -------------------------------------------------------------------------- */
/*                                  Response                                  */
/* -------------------------------------------------------------------------- */
type Resp struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type Mock struct {
	URL    string `json:"url"`
	Action string `json:"action"`
}
