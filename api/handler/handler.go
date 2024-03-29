package handler

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/photon-storage/go-common/api/pagination"
	"github.com/photon-storage/go-common/log"
)

const reqBodyLabel = "request_body_label"

var (
	errorType        = reflect.TypeOf((*error)(nil)).Elem()
	contextType      = reflect.TypeOf((*gin.Context)(nil))
	paginationType   = reflect.TypeOf((*pagination.Query)(nil))
	paginationResult = reflect.TypeOf((*pagination.Result)(nil))
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

type Handler struct {
	errCodes map[error]int
}

func New(errCodes map[error]int) *Handler {
	return &Handler{errCodes: errCodes}
}

type handleFunc any

func (h *Handler) Handle(fn handleFunc) gin.HandlerFunc {
	if err := validateFunc(fn); err != nil {
		log.Fatal("validate service handle func failed",
			"error", err,
		)
	}

	return func(ctx *gin.Context) {
		ft := reflect.TypeOf(fn)
		args, err := buildInputParams(ft, ctx)
		if err != nil {
			h.errResponse(ctx, err)
			return
		}

		result := callHandleFunc(fn, args...)
		if err := result[len(result)-1]; err != nil {
			h.errResponse(ctx, err.(error))
			return
		}

		// If handle function has already sent response, we could just stop here.
		if ctx.IsAborted() {
			return
		}

		if ft.In(ft.NumIn()-1) == paginationType {
			r := result[0].(*pagination.Result)
			query := args[len(args)-1].(*pagination.Query)
			ctx.AbortWithStatusJSON(http.StatusOK, &pagination.Response{
				Code:   http.StatusOK,
				Result: r,
				Links:  pagination.GetLinks(ctx, r.Total, query),
			},
			)

			return
		}

		ctx.AbortWithStatusJSON(
			http.StatusOK,
			Response{
				Code: http.StatusOK,
				Msg:  "ok",
				Data: result[0],
			},
		)
	}
}

func buildInputParams(ft reflect.Type, ctx *gin.Context) ([]any, error) {
	args := []any{ctx}
	if ft.NumIn() == 1 {
		return args, nil
	}

	if ft.In(1) != paginationType {
		reqArg := reflect.New(ft.In(1).Elem()).Interface()
		if err := ctx.ShouldBindJSON(reqArg); err != nil {
			return nil, err
		}

		reqBytes, err := json.Marshal(reqArg)
		if err != nil {
			return nil, err
		}

		ctx.Set(reqBodyLabel, string(reqBytes))
		if err := validator.New().Struct(reqArg); err != nil {
			return nil, err
		}

		args = append(args, reqArg)
	}

	if ft.In(ft.NumIn()-1) == paginationType {
		query, err := pagination.Parse(ctx)
		if err != nil {
			return nil, err
		}

		args = append(args, query)
	}

	return args, nil
}

func callHandleFunc(fn handleFunc, args ...any) []any {
	params := make([]reflect.Value, len(args))
	for i, arg := range args {
		params[i] = reflect.ValueOf(arg)
	}

	rs := reflect.ValueOf(fn).Call(params)
	result := make([]any, len(rs))
	for i, r := range rs {
		result[i] = r.Interface()
	}
	return result
}

func validateFunc(fn handleFunc) error {
	ft := reflect.TypeOf(fn)
	if ft.Kind() != reflect.Func || ft.IsVariadic() {
		return errors.Errorf("need non variadic func in %s" + ft.String())
	}

	if ft.NumIn() < 1 || ft.NumIn() > 3 {
		return errors.Errorf("the size of input parameters is " +
			"not correct in %s" + ft.String())
	}

	if ft.In(0) != contextType {
		return errors.New("the first parameter must point of context " +
			"in %s" + ft.String())
	}

	if ft.NumIn() == 2 && ft.In(1).Kind() != reflect.Ptr {
		return errors.Errorf("the second parameter must be a "+
			"pointer type in %s", ft.String())
	}

	if ft.NumOut() < 1 || ft.NumOut() > 2 {
		return errors.Errorf("the number of return values must be "+
			"one or two in %s", ft.String())
	}

	if ft.In(ft.NumIn()-1) == paginationType && ft.Out(0) != paginationResult {
		return errors.Errorf("the last of input parameter is "+
			"pagginationQuery type, the first return value must be "+
			"a paginationResult type in %s", ft.String())
	}

	if !ft.Out(ft.NumOut() - 1).Implements(errorType) {
		return errors.Errorf("the last return value must be an " +
			"error type in %s" + ft.String())
	}

	return nil
}

func (h *Handler) errResponse(c *gin.Context, err error) {
	log.Error("Error requesting the api server",
		"url", c.Request.URL,
		"request_body", c.Value(reqBodyLabel),
		"error", err,
	)
	code := getErrCode(err, h.errCodes)
	msg := err.Error()
	c.AbortWithStatusJSON(http.StatusBadRequest, Response{
		Code: code,
		Msg:  msg,
	})
}

func getErrCode(err error, errorCodes map[error]int) int {
	if ok := isComparable(reflect.TypeOf(err)); ok {
		if errCode, ok := errorCodes[err]; ok {
			return errCode
		}
	}

	return -1
}

func isComparable(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Slice, reflect.Func, reflect.Map:
		return false
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			if !isComparable(typ.Field(i).Type) {
				return false
			}
		}
	}

	return true
}
