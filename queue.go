package golink

import (

)

type Result struct {
	URL string `json:"url"`
	Determination string `json:"determination,omitempty"`
	Reason string `json:"reason,omitempty"`
	Status string `json:"status"`
}

type Queue struct {
	urls []string
	in chan string
	out chan Result
	done chan Result
	maxWorkers int
}

func (q *Queue) Run() {
	for {
		select {
		case url := <-q.in:
			q.urls = append(q.urls, url)
			q.out <- Result{
				URL: url,
				Determination: "TEST",
				Reason: "CAUSE",
				Status: "checked",
			}
		}
	}
}

func newQueue (in chan string, out chan Result) *Queue {
	q := &Queue{
		urls: make([]string, 0),
		in: in,
		out: out,
		maxWorkers: 5,
		done: make (chan Result),
	}

	go q.Run()

	return q
}
