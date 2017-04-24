package endpoint

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/aeppert/swag/swagger"
)

// Builder uses the builder pattern to generate swagger endpoint definitions
type Builder struct {
	Endpoint *swagger.Endpoint
}

// Option represents a functional option to customize the swagger endpoint
type Option func(builder *Builder)

// Handler allows an instance of the web handler to be associated with the endpoint.  This can be especially useful when
// using swag to bind the endpoints to the web router.  See the examples package for how the Handler can be used in
// conjunction with Walk to simplify binding endpoints to a router
func Handler(handler interface{}) Option {
	return func(b *Builder) {
		if v, ok := handler.(func(w http.ResponseWriter, r *http.Request)); ok {
			handler = http.HandlerFunc(v)
		}
		b.Endpoint.Handler = handler
	}
}

// Description sets the endpoint's description
func Description(v string) Option {
	return func(b *Builder) {
		b.Endpoint.Description = v
	}
}

// OperationID sets the endpoint's operationId
func OperationID(v string) Option {
	return func(b *Builder) {
		b.Endpoint.OperationID = v
	}
}

// Produces sets the endpoint's produces; by default this will be set to application/json
func Produces(v ...string) Option {
	return func(b *Builder) {
		b.Endpoint.Produces = v
	}
}

// Consumes sets the endpoint's produces; by default this will be set to application/json
func Consumes(v ...string) Option {
	return func(b *Builder) {
		b.Endpoint.Consumes = v
	}
}

func parameter(p swagger.Parameter) Option {
	return func(b *Builder) {
		if b.Endpoint.Parameters == nil {
			b.Endpoint.Parameters = []swagger.Parameter{}
		}

		b.Endpoint.Parameters = append(b.Endpoint.Parameters, p)
	}
}

// Path defines a path parameter for the endpoint; name, typ, description, and required correspond to the matching
// swagger fields
func Path(name, typ, description string, required bool) Option {
	p := swagger.Parameter{
		Name:        name,
		In:          "path",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

// Query defines a query parameter for the endpoint; name, typ, description, and required correspond to the matching
// swagger fields
func Query(name, typ, description string, required bool) Option {
	p := swagger.Parameter{
		Name:        name,
		In:          "query",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

// FormData defines a formData parameter for the endpoint; name, typ, description, and required correspond to the matching
// swagger fields
func FormData(name, typ, description string, required bool) Option {
	p := swagger.Parameter{
		Name:        name,
		In:          "formData",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

// Body defines a body parameter for the swagger endpoint as would commonly be used for the POST, PUT, and PATCH methods
// prototype should be a struct or a pointer to struct that swag can use to reflect upon the return type
func Body(prototype interface{}, description string, required bool) Option {
	p := swagger.Parameter{
		In:          "body",
		Name:        "body",
		Description: description,
		Schema:      swagger.MakeSchema(prototype),
		Required:    required,
	}
	return parameter(p)
}

// Tags allows one or more tags to be associated with the endpoint
func Tags(tags ...string) Option {
	return func(b *Builder) {
		if b.Endpoint.Tags == nil {
			b.Endpoint.Tags = []string{}
		}

		b.Endpoint.Tags = append(b.Endpoint.Tags, tags...)
	}
}

// ResponseOption allows for additional configurations on responses like header information
type ResponseOption func(response *swagger.Response)

// Header adds header definitions to swagger responses
func Header(name, typ, format, description string) ResponseOption {
	return func(response *swagger.Response) {
		if response.Headers == nil {
			response.Headers = map[string]swagger.Header{}
		}

		response.Headers[name] = swagger.Header{
			Type:        typ,
			Format:      format,
			Description: description,
		}
	}
}

// Response sets the endpoint response for the specified code; may be used multiple times with different status codes
func Response(code int, prototype interface{}, description string, opts ...ResponseOption) Option {
	return func(b *Builder) {
		if b.Endpoint.Responses == nil {
			b.Endpoint.Responses = map[string]swagger.Response{}
		}

		r := swagger.Response{
			Description: description,
			Schema:      swagger.MakeSchema(prototype),
		}

		for _, opt := range opts {
			opt(&r)
		}

		b.Endpoint.Responses[strconv.Itoa(code)] = r
	}
}

// New constructs a new swagger endpoint using the fields and functional options provided
func New(method, path, summary string, options ...Option) *swagger.Endpoint {
	method = strings.ToUpper(method)
	e := &Builder{
		Endpoint: &swagger.Endpoint{
			Method:      method,
			Path:        path,
			Summary:     summary,
			OperationID: strings.ToLower(method) + camel(path),
			Produces:    []string{"application/json"},
			Consumes:    []string{"application/json"},
			Tags:        []string{},
		},
	}

	for _, opt := range options {
		opt(e)
	}

	return e.Endpoint
}
