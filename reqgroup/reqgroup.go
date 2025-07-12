package reqgroup

import (
	"fmt"
	"sync"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
)

type Result struct {
	response contracts.Response
	err      error
}

type Group struct {
	results []Result
	wg      sync.WaitGroup
}

func (g *Group) Results() []Result {
	return g.results
}

func (g *Group) Result(i int) (contracts.Response, error) {
	length := len(g.results)
	if i > length-1 {
		panic(fmt.Sprintf("can not access result index: %d", i))
	}

	if i < 0 {
		panic("can not access result with negative index")
	}

	result := g.results[i]

	return result.response, result.err
}

func (g *Group) Wait() {
	g.wg.Wait()
}

func All(builders ...contracts.RequestBuilder) (*Group, error) {
	if len(builders) == 0 {
		return nil, ErrNoRequests
	}

	g := &Group{
		results: make([]Result, len(builders)),
		wg:      sync.WaitGroup{},
	}
	g.wg.Add(len(builders))

	for i, builder := range builders {
		if builder == nil {
			return nil, ErrNilRequestBuilder
		}

		go func(i int, b contracts.RequestBuilder) {
			defer g.wg.Done()

			resp, err := b.Send()

			g.results[i] = Result{
				response: resp,
				err:      err,
			}
		}(i, builder)
	}

	return g, nil
}
