package interfaces

import "github.com/labstack/echo/v4"

type ResponseInterface[T any] struct {
	Message string `json:"message,omitempty"`
	Result  T      `json:"result,omitempty"`
}

type HttpResponse interface {
	isHttpResponse()
}

func (ResponseInterface[T]) isHttpResponse() {}

func Set[T any](c echo.Context, resp ResponseInterface[T]) {
	c.Set("response", resp)
}
