package logger

import (
	"encoding/json"
	"io"
	"log"
	"sync"
	"time"
)

type Entry struct {
	Time      time.Time      `json:"time"`
	Event     string         `json:"event"`
	RequestID string         `json:"request_id,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
	Error     string         `json:"error,omitempty"`
}

type Logger interface {
	Log(Entry)
	Stop()
}

type Async struct {
	ch      chan Entry
	wg      sync.WaitGroup
	closeMu sync.Mutex
	closed  bool

	out *log.Logger
}

func NewAsync(buffer int, out io.Writer) *Async {
	l := &Async{
		ch:  make(chan Entry, buffer),
		out: log.New(out, "", 0),
	}
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for e := range l.ch {
			b, err := json.Marshal(e)
			if err != nil {
				l.out.Printf(`{"time":"%s","event":"logger_error","error":"%s"}`, time.Now().UTC().Format(time.RFC3339Nano), err.Error())
				continue
			}
			l.out.Println(string(b))
		}
	}()
	return l
}

func (a *Async) Log(e Entry) {
	a.closeMu.Lock()
	closed := a.closed
	a.closeMu.Unlock()
	if closed {
		return
	}
	select {
	case a.ch <- e:
	default:
		a.out.Printf(`{"time":"%s","event":"log_dropped"}`, time.Now().UTC().Format(time.RFC3339Nano))
	}
}

func (a *Async) Stop() {
	a.closeMu.Lock()
	if a.closed {
		a.closeMu.Unlock()
		return
	}
	a.closed = true
	close(a.ch)
	a.closeMu.Unlock()
	a.wg.Wait()
}
