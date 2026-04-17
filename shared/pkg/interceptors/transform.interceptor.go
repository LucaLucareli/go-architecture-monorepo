package interceptors

import (
	"net/http"
	"reflect"
	"shared/application/interfaces"

	"github.com/labstack/echo/v4"
)

func TransformInterceptor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			return err
		}

		result, ok := c.Get("response").(interfaces.HttpResponse)
		if !ok {
			return nil
		}

		resultValue := reflect.ValueOf(result)
		resultType := reflect.TypeOf(result)

		if resultValue.Kind() == reflect.Struct {
			var messageField any
			var dataField any

			for i := 0; i < resultValue.NumField(); i++ {
				field := resultType.Field(i)
				value := resultValue.Field(i).Interface()

				switch field.Name {
				case "Message":
					messageField = value
				case "Result":
					dataField = value
				}
			}

			return c.JSON(http.StatusOK, map[string]any{
				"message": messageField,
				"data":    dataField,
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}
