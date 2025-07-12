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
