package resp

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/domain/model/common"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/amahdian/cliplab-be/pkg/msg"
	"github.com/amahdian/cliplab-be/server/binding"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

func AbortWithError(ctx *gin.Context, err error) {
	var ve validator.ValidationErrors
	var customErr *errs.Error
	switch {
	case errors.Is(err, io.EOF):
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(fmt.Errorf("Request payload cannot be empty.")))
		return
	case errors.As(err, &ve):
		mc := binding.MapValidationErrorsToMessageContainer(ve)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(mc))
		return
	case errors.As(err, &errs.EntryNotFoundErr{}), errors.As(err, &errs.UserSettingFoundErr{}):
		ctx.AbortWithStatusJSON(http.StatusNotFound, NewErrorResponse(err))
		return
	case errors.As(err, &customErr):
		ctx.AbortWithStatusJSON(customErr.Code.HttpStatus(), NewErrorResponse(err))
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}
}

func AbortWithStatus(ctx *gin.Context, code int, data interface{}) {
	switch v := data.(type) {
	case error:
		{
			var ve validator.ValidationErrors
			switch {
			case errors.As(v, &ve):
				mc := binding.MapValidationErrorsToMessageContainer(ve)
				ctx.AbortWithStatusJSON(code, NewErrorResponse(mc))
				return
			default:
				ctx.AbortWithStatusJSON(code, NewErrorResponse(v))
				return
			}
		}
	case *msg.MessageContainer:
		{
			ctx.AbortWithStatusJSON(code, NewErrorResponse(v))
			return
		}
	default:
		logger.Errorf("Could not recognize the error type for proper error handling: %v", data)
	}
}

func Success(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, NewResponse(true))
}

func Ok(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, NewResponse(data))
}

func OkWithMessage(ctx *gin.Context, data interface{}, messages *msg.MessageContainer) {
	ctx.JSON(http.StatusOK, NewResponseWithMessage(data, messages))
}

func Stream(ctx *gin.Context, streamChan <-chan *model.StreamedMessage) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")

	ctx.Stream(func(w io.Writer) bool {
		if message, ok := <-streamChan; ok {
			// Stream assistant reply
			ctx.SSEvent(string(message.Type), gin.H{
				"content": message.Content,
			})
			return true
		}

		return false
	})
}

func ConflictWithMessage(ctx *gin.Context, data interface{}, messages *msg.MessageContainer) {
	ctx.JSON(http.StatusConflict, NewResponseWithMessage(data, messages))
}

func PaginatedOk[T any](ctx *gin.Context, data []T, pagination *common.Pagination) {
	ctx.JSON(http.StatusOK, NewPaginatedResponse(data, pagination))
}

func Created(ctx *gin.Context, data ...interface{}) {
	if len(data) > 0 {
		ctx.JSON(http.StatusCreated, NewResponse(data[0]))
	} else {
		ctx.JSON(http.StatusCreated, NewResponse(true))
	}
}

func NoContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func CreatedWithMessage(ctx *gin.Context, data interface{}, messages *msg.MessageContainer) {
	ctx.JSON(http.StatusCreated, NewResponseWithMessage(data, messages))
}

func WriteJsonFileBytes(ctx *gin.Context, data []byte, fileName string) {
	fileName = strings.TrimSuffix(fileName, ".json") // remove suffix ".json" if present
	fileName = fmt.Sprintf("%s.json", fileName)      // add suffix ".json"
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q; filename*=utf-8''%q", fileName, fileName))
	_, _ = ctx.Writer.Write(data) // ignore error
}

func WriteZipFileBytes(ctx *gin.Context, data []byte, fileName string) {
	fileName = strings.TrimSuffix(fileName, ".zip") // remove suffix ".zip" if present
	fileName = fmt.Sprintf("%s.zip", fileName)      // add suffix ".zip"
	ctx.Writer.Header().Set("Content-Type", "application/zip")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q; filename*=utf-8''%q", fileName, fileName))
	_, _ = ctx.Writer.Write(data) // ignore error
}

func Redirect(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusSeeOther, location)
}
