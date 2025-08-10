package contracts

// RequestBuilder builds and executes an HTTP request. It exposes fluent
// builders for headers, body, retries, context and query parameters. The
// final Send call performs the request and returns a Response.
//
// Example:
//
//	resp, err := client.
//	        GET("/users").
//	        Query().AddParam("active", "1").
//	        Send()
type RequestBuilder interface {
	// Header returns a builder for configuring request headers.
	Header() BuilderHeader[RequestBuilder]
	// Body returns a builder for configuring the request body.
	Body() BuilderRequestBody[RequestBuilder]
	// Retry returns a builder for configuring retry logic.
	Retry() BuilderRequestRetry[RequestBuilder]
	// Context returns a builder for setting the request context.
	Context() BuilderRequestContext[RequestBuilder]
	// Query returns a builder for setting query parameters.
	Query() BuilderRequestQuery[RequestBuilder]

	// Send executes the HTTP request.
	Send() (Response, error)
}
