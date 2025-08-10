// Package maigo provides a library for making HTTP requests with built-in retry logic and error handling.
//
// This package simplifies making HTTP requests in Go applications, offering features like:
//
//   - **Retry Logic:** Automatic retries with configurable backoff.
//   - **Error Handling:** Comprehensive error reporting and management.
//   - **Type Safety:**  Strongly typed parameters and response data for enhanced reliability.
//
// The package defines the following key types:
//
//   - `Client`: Represents a client for making HTTP requests.  Provides methods for sending requests and handling responses.
//   - `Body`: Represents the request body, including headers and parameters.
//   - `Resource`: Represents the structure of HTTP responses.
//
// The package also provides the following functions:
//
//   - `CreateClient()`: Creates a new `Client` instance.
//   - `Get(url string) (*Body, error)`: Sends a GET request to the specified URL.
//   - `Post(url string, body interface{}) (*Body, error)`: Sends a POST request with the provided data.
//   - `Patch(url string, field string, value interface{}) (*Body, error)`: Sends a PATCH request to update a resource.
//   - `Delete(url string) (*Body, error)`: Sends a DELETE request to remove a resource.
package maigo
