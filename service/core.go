package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

func New(client Interface) *Dialer {
	var urls urls

	urls = append(urls, url{
		url: "url_zero",
		key: "api-key",
		ch:  make(chan struct{}),
	}, url{
		url: "url_one",
		key: "api-key",
		ch:  make(chan struct{}),
	}, url{
		url: "url_two",
		key: "api-key",
		ch:  make(chan struct{}),
	}, url{
		url: "url_three",
		key: "api-key",
		ch:  make(chan struct{}),
	})

	return &Dialer{
		current: 0,
		urls:    urls,
		Client:  client,
		close:   make(map[int]struct{}),
	}
}

func (d *Dialer) Middleware(req Req) (Resp, error) {
	var (
		res    Resp
		data   []byte
		err    error
		status int
	)

	for {
		url := d.urls[d.current].url
		key := d.urls[d.current].key

		data, status, err = d.Client.Get(url, key)
		if err != nil || status != http.StatusOK {
			if d.Last() {
				d.urls[len(d.urls)-1].ch <- struct{}{}
				return res, errors.New("internal server error") //all urls are down
			}

			d.Next()

		} else {
			break
		}
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, err

}

func (d *Dialer) Last() bool {
	d.mu.Lock()
	current := d.current
	d.mu.Unlock()
	return current == len(d.urls)-1
}

func (d *Dialer) Next() {
	d.mu.Lock()
	d.urls[d.current].ch <- struct{}{}
	d.current += 1
	d.mu.Unlock()
}

func (d *Dialer) StartHealth() {
	for i := 0; i < len(d.urls); i++ {
		go d.Single(i)
	}
}

func (d *Dialer) Single(index int) {
	for {
		log.Println("waiting on url: ", index)
		<-d.urls[index].ch

		d.close[index] = struct{}{}
		d.urls[index].close = make(chan struct{}) //flush the channel if any value persists

		d.Health(index)
	}
}

func (d *Dialer) Health(index int) {
	for {
		select {
		case <-d.urls[index].close:
			log.Println("closed health check for url no: ", index)
			return
		default:
			log.Println("checking health for url no: ", index)
			time.Sleep(1 * time.Second)

			url := d.urls[index].url
			key := d.urls[index].key

			_, status, err := d.Client.Get(url, key) //health endpoint should be added here
			if err != nil || status != http.StatusOK {
				log.Println("url no ", index, " is down.")
				continue
			}

			d.mu.Lock()
			if d.current > index { //edge case where both 2 and 3 have reached Lock()
				log.Println("index moved from ", d.current, " to ", index, ".")
				d.current = index
				d.Close(index)
			}

			delete(d.close, index)
			d.mu.Unlock()
			return
		}
	}
}

// closes all the channels below the given index to stop health checks
func (d *Dialer) Close(index int) {
	for key := range d.close {
		if key > index {
			delete(d.close, key)
			d.urls[key].close <- struct{}{}
		}
	}
}
