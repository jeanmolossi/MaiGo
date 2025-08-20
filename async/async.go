// Package async provides an easy way to dispatch a request
// asynchronously and catch the result when needed
package async

import (
	"fmt"
	"sync"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
)

type result struct {
	response contracts.Response
	err      error
}

// DispatchResult wraps the request response and provides a response [.Response()]
// wich waits until receive response and returns with the result
type DispatchResult struct {
	result
	wg sync.WaitGroup
}

// Response waits until the request ends and read the response and error
func (r *DispatchResult) Response() (contracts.Response, error) {
	r.wg.Wait()
	return r.response, r.err
}

// Dispatch the request in a goroutine then provide a easy way
// to retrieve the response of this request results
//
// Example:
//
//	client := maigo.NewClient("http://example.com").Build()
//
//	result, err := async.Dispatch(client.GET("/resource"))
//	if err != nil {
//	    // handle error
//	}
//
//	// Do some stuff that do not need the request response
//
//	// catch the result of request who is dispatched
//	resp, err := result.Response()
//	if err != nil {
//	    // handle response error
//	}
func Dispatch(req contracts.RequestBuilder) (*DispatchResult, error) {
	if req == nil {
		return nil, ErrEmptyRequest
	}

	result := &DispatchResult{
		wg: sync.WaitGroup{},
	}

	result.wg.Add(1)

	go func() {
		result.response, result.err = req.Send()
		result.wg.Done()
	}()

	return result, nil
}

// Group of requests that are dispatched.
type Group struct {
	results []result
	wg      sync.WaitGroup
}

// Results of the requests. The results always be the same size of
// the number of requests passed to [All].
//
// The indexes can be nil if request still in transit.
//
// PAY ATTENTION: that is not a slice copy. Then, if you modify the
// result it affects the original value.
func (g *Group) Results() []result {
	return g.results
}

// Size of results. It have the same size of the number of requests
// called in [All].
//
// Example:
//
//	     // limit concurrent requests to 2
//	     group, err := async.All(2,
//	         client.GET("/users/1"),
//	         client.GET("/users/2"),
//	         client.GET("/users/3"),
//	     )
//		if err != nil {
//		    // handle err
//		}
//
//		group.Size() // 3 event requests still in transit.
func (g *Group) Size() int {
	return len(g.results)
}

// Wait for all requests.
func (g *Group) Wait() {
	g.wg.Wait()
}

// Result read a single result from index [i int]. It panics if
// index greather than group size or is a negative int.
func (g *Group) Result(i int) (contracts.Response, error) {
	length := len(g.results)
	if i > length-1 {
		return nil, fmt.Errorf("index %d out ot range [0-%d]", i, length-1)
	}

	if i < 0 {
		return nil, fmt.Errorf("negative index %d not allowed", i)
	}

	result := g.results[i]

	return result.response, result.err
}

// All dispatches requests in their own goroutines.
// If limit <= 0, all requests run concurrently; otherwise at most limit
// requests are executed simultaneously.
func All(limit int, builders ...contracts.RequestBuilder) (*Group, error) {
	if len(builders) == 0 {
		return nil, ErrNoRequests
	}

	g := &Group{
		results: make([]result, len(builders)),
		wg:      sync.WaitGroup{},
	}
	g.wg.Add(len(builders))

	var sem chan struct{}
	if limit > 0 {
		sem = make(chan struct{}, limit)
	}

	for i, builder := range builders {
		if builder == nil {
			return nil, ErrNilRequestBuilder
		}

		if sem != nil {
			sem <- struct{}{}
		}

		go func(i int, b contracts.RequestBuilder) {
			defer g.wg.Done()
			if sem != nil {
				defer func() { <-sem }()
			}

			resp, err := b.Send()

			g.results[i] = result{
				response: resp,
				err:      err,
			}
		}(i, builder)
	}

	return g, nil
}
