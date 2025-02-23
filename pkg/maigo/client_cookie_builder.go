package maigo

import (
	"net/http"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
)

var _ contracts.BuilderCookie[contracts.ClientBuilder] = (*ClientCookieBuilder)(nil)

type ClientCookieBuilder struct {
	parent *ClientBuilder
}

func (b *ClientBuilder) Cookie() contracts.BuilderCookie[contracts.ClientBuilder] {
	return &ClientCookieBuilder{parent: b}
}

// Add implements contracts.BuilderCookie.
func (c *ClientCookieBuilder) Add(cookie *http.Cookie) contracts.ClientBuilder {
	c.parent.client.Cookies().Add(cookie)
	return c.parent
}
