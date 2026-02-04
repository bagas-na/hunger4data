package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapGRPCError(c echo.Context, err error) error {
	status, ok := status.FromError(err)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	switch status.Code() {
	case codes.InvalidArgument: // 400
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": status.Message(),
		})
	case codes.Unauthenticated: // 401
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"message": status.Message(),
		})
	case codes.PermissionDenied: // 403
		return c.JSON(http.StatusForbidden, map[string]any{
			"message": status.Message(),
		})
	case codes.NotFound: // 404
		return c.JSON(http.StatusNotFound, map[string]any{
			"message": status.Message(),
		})
	case codes.AlreadyExists: // 409
		return c.JSON(http.StatusConflict, map[string]any{
			"message": status.Message(),
		})
	case codes.FailedPrecondition: // 412
		return c.JSON(http.StatusPreconditionFailed, map[string]any{
			"message": status.Message(),
		})
	default: // 500
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "internal server error",
			"detail":  err.Error(),
		})
	}
}
