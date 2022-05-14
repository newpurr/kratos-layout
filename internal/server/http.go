package server

import (
	"context"
	"github.com/newpurr/stock/api/user"
	netHttp "net/http"

	"github.com/go-kratos/grpc-gateway/v2/protoc-gen-openapiv2/generator"
	"github.com/go-kratos/swagger-api/openapiv2"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/newpurr/stock/internal/conf"
	"github.com/newpurr/stock/internal/data"
	"github.com/newpurr/stock/internal/service"
)

// Validator is a validator middleware.
func Validator() middleware.Middleware {
	type validator interface {
		Validate() error
		ValidateAll() error
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(validator); ok {
				if err := v.ValidateAll(); err != nil {
					//return nil, errors.BadRequest("VALIDATOR", err.Error())
					return nil, err
				}
			}
			return handler(ctx, req)
		}
	}
}

type HttpResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    struct {
		Trace struct {
			TraceId   string `json:"trace_id"`
			RequestId string `json:"request_id"`
		} `json:"trace"`
	} `json:"meta"`
	Errors     []interface{} `json:"errors"`
	StatusCode string        `json:"status_code"`
}

type MultiError interface {
	Error() string
	AllErrors() []error
}
type ValidationError interface {
	Field() string
	Reason() string
	Cause() error
	Key() bool
	ErrorName() string
	Error() string
}

// ErrorEncoder encodes the error to the HTTP response.
func ErrorEncoder(w netHttp.ResponseWriter, r *netHttp.Request, err error) {
	se := HttpResponse{
		Message: err.Error(),
	}
	multiError, ok := err.(MultiError)
	if ok {
		se.Message = "参数验证失败"
		for _, err := range multiError.AllErrors() {
			validationError, ok := err.(ValidationError)
			if ok {
				se.Errors = append(se.Errors, struct {
					Field  string `json:"field"`
					Reason string `json:"reason"`
				}{
					validationError.Field(),
					validationError.Reason(),
				})
				continue
			}

			se.Errors = append(se.Errors, err.Error())
		}
	}

	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(netHttp.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	_, _ = w.Write(body)
}

// NewHTTPServer new an HTTP server. at the start
func NewHTTPServer(
	c *conf.Server,
	logger log.Logger,
	source data.DefaultDataSource,
	us *service.UserService,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			Validator(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	opts = append(opts, http.ErrorEncoder(ErrorEncoder))
	srv := http.NewServer(opts...)

	openAPIhandler := openapiv2.NewHandler(openapiv2.WithGeneratorOptions(generator.UseJSONNamesForFields(false), generator.EnumsAsInts(true)))
	srv.HandlePrefix("/q/", openAPIhandler)

	// register
	user.RegisterUserHTTPServer(srv, us)

	return srv
}
