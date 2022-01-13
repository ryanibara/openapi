package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

const (
	// SecuritySchemeTypeAPIKey = "apiKey"
	SecuritySchemeTypeAPIKey SecuritySchemeType = "apiKey"
	// SecuritySchemeTypeHTTP = "http"
	SecuritySchemeTypeHTTP SecuritySchemeType = "http"
	// SecuritySchemeTypeMutualTLS = mutualTLS
	SecuritySchemeTypeMutualTLS SecuritySchemeType = "mutualTLS"
	// SecuritySchemeTypeOAuth2 = oauth2
	SecuritySchemeTypeOAuth2 SecuritySchemeType = "oauth2"
	// SecuritySchemeTypeOpenIDConnect = "openIdConnect"
	SecuritySchemeTypeOpenIDConnect SecuritySchemeType = "openIdConnect"
)

// SecurityRequirements is a list of SecurityRequirement
type SecurityRequirements []SecurityRequirement

// Kind returns KindSecurityRequirements
func (SecurityRequirements) Kind() Kind {
	return KindSecurityRequirements
}

// SecurityRequirement lists the required security schemes to execute this
// operation. The name used for each property MUST correspond to a security
// scheme declared in the Security Schemes under the Components Object.
//
// Security Requirement Objects that contain multiple schemes require that all
// schemes MUST be satisfied for a request to be authorized. This enables
// support for scenarios where multiple query parameters or HTTP headers are
// required to convey security information.
//
// When a list of Security Requirement Objects is defined on the OpenAPI Object
// or Operation Object, only one of the Security Requirement Objects in the list
// needs to be satisfied to authorize the request.
//
// Each name MUST correspond to a security scheme which is declared in the
// Security Schemes under the Components Object. If the security scheme is of
// type "oauth2" or "openIdConnect", then the value is a list of scope names
// required for the execution, and the list MAY be empty if authorization does
// not require a specified scope. For other security scheme types, the array MAY
// contain a list of role names which are required for the execution, but are
// not otherwise defined or exchanged in-band.
type SecurityRequirement map[string][]string

// Kind returns KindSecurityRequirement
func (SecurityRequirement) Kind() Kind {
	return KindSecurityRequirement
}

// SecuritySchemeType represents the type of the security scheme.
type SecuritySchemeType string

func (ss SecuritySchemeType) String() string {
	return string(ss)
}

// SecuritySchemes is a map of SecurityScheme
type SecuritySchemes map[string]SecurityScheme

func (ss *SecuritySchemes) Len() int {
	if ss == nil || *ss == nil {
		return 0
	}
	return len(*ss)
}

func (ss *SecuritySchemes) Get(key string) (SecurityScheme, bool) {
	if ss == nil || *ss == nil {
		return nil, false
	}
	v, ok := (*ss)[key]
	return v, ok
}

func (ss *SecuritySchemes) Set(key string, val SecurityScheme) {
	if *ss == nil {
		*ss = SecuritySchemes{
			key: val,
		}
		return
	}
	(*ss)[key] = val
}

func (ss SecuritySchemes) Nodes() Nodes {
	if len(ss) == 0 {
		return nil
	}
	nodes := make(Nodes, len(ss))
	for k, v := range ss {
		nodes[k] = NodeDetail{
			TargetKind: KindSecurityScheme,
			Node:       v,
		}
	}
	return nodes
}

// Kind returns KindSecuritySchemes
func (SecuritySchemes) Kind() Kind {
	return KindSecuritySchemes
}

// UnmarshalJSON unmarshals json
func (ss *SecuritySchemes) UnmarshalJSON(data []byte) error {
	var dm map[string]json.RawMessage
	if err := json.Unmarshal(data, &dm); err != nil {
		return err
	}
	res := make(SecuritySchemes, len(dm))
	for k, d := range dm {
		if isRefJSON(d) {
			v, err := unmarshalReferenceJSON(d)
			if err != nil {
				return err
			}
			res[k] = v
			continue
		}
		var v SecuritySchemeObj
		if err := json.Unmarshal(d, &v); err != nil {
			return err
		}
		res[k] = &v
	}
	*ss = res
	return nil
}

// MarshalYAML marshals YAML
func (ss SecuritySchemes) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(ss)
}

// UnmarshalYAML unmarshals YAML
func (ss *SecuritySchemes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, ss)
}

// SecuritySchemeObj defines a security scheme that can be used by the operations.
type SecuritySchemeObj struct {
	// The type of the security scheme.
	//
	// *required
	Type SecuritySchemeType `json:"type,omitempty"`
	// Any description for security scheme. CommonMark syntax MAY be used for
	// rich text representation.
	Description string `json:"description,omitempty"`
	// The name of the header, query or cookie parameter to be used.
	//
	// Applies to: API Key
	//
	// 	*required*
	Name string `json:"name,omitempty"`
	// The location of the API key. Valid values are "query", "header" or "cookie".
	//
	// Applies to: APIKey
	//
	// 	*required*
	In In `json:"in,omitempty"`
	// The name of the HTTP Authorization scheme to be used in the Authorization
	// header as defined in RFC7235. The values used SHOULD be registered in the
	// IANA Authentication Scheme registry.
	//
	// 	*required*
	Scheme string `json:"scheme,omitempty"`

	// http ("bearer")  A hint to the client to identify how the bearer token is
	// formatted. Bearer tokens are usually generated by an authorization
	// server, so this information is primarily for documentation purposes.
	BearerFormat string `json:"bearerFormat,omitempty"`

	// An object containing configuration information for the flow types supported.
	//
	// 	*required*
	Flows *OAuthFlows `json:"flows,omitempty"`

	// OpenId Connect URL to discover OAuth2 configuration values. This MUST be
	// in the form of a URL. The OpenID Connect standard requires the use of
	// TLS.
	//
	// 	*required*
	OpenIDConnectURL string `json:"openIdConnect,omitempty"`
	Extensions       `json:"-"`
}

func (ss *SecuritySchemeObj) Nodes() Nodes {
	return makeNodes(nodes{
		{"flows", ss.Flows, KindOAuthFlows},
	})
}

type securityscheme SecuritySchemeObj

// UnmarshalJSON unmarshals JSON
func (sso *SecuritySchemeObj) UnmarshalJSON(data []byte) error {
	var v securityscheme
	err := unmarshalExtendedJSON(data, &v)
	*sso = SecuritySchemeObj(v)
	return err
}

// MarshalJSON marshals JSON
func (sso SecuritySchemeObj) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(securityscheme(sso))
}

// MarshalYAML marshals YAML
func (sso SecuritySchemeObj) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(sso)
}

// UnmarshalYAML unmarshals YAML
func (sso *SecuritySchemeObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, sso)
}

// ResolveSecurityScheme resolves SecuritySchemeObj by returning itself. resolve is  not called.
func (sso *SecuritySchemeObj) ResolveSecurityScheme(func(ref string) (*SecuritySchemeObj, error)) (*SecuritySchemeObj, error) {
	return sso, nil
}

// Kind returns KindSecurityScheme
func (*SecuritySchemeObj) Kind() Kind {
	return KindSecurityScheme
}

// SecurityScheme can either be a ScecuritySchemeObj or a Reference
type SecurityScheme interface {
	ResolveSecurityScheme(func(ref string) (*SecuritySchemeObj, error)) (*SecuritySchemeObj, error)
	Kind() Kind
}

// ResolvedSecuritySchemes is a map of *ResolvedSecurityScheme
type ResolvedSecuritySchemes map[string]*ResolvedSecurityScheme

// Kind returns KindResolvedSecuritySchemes
func (ResolvedSecuritySchemes) Kind() Kind {
	return KindResolvedSecuritySchemes
}

func (rss *ResolvedSecuritySchemes) Len() int {
	if rss == nil || *rss == nil {
		return 0
	}
	return len(*rss)
}

func (rss *ResolvedSecuritySchemes) Get(key string) (*ResolvedSecurityScheme, bool) {
	if rss == nil || *rss == nil {
		return nil, false
	}
	v, ok := (*rss)[key]
	return v, ok
}

func (rss *ResolvedSecuritySchemes) Set(key string, val *ResolvedSecurityScheme) {
	if *rss == nil {
		*rss = ResolvedSecuritySchemes{
			key: val,
		}
		return
	}
	(*rss)[key] = val
}

func (rss ResolvedSecuritySchemes) Nodes() Nodes {
	if len(rss) == 0 {
		return nil
	}
	nodes := make(Nodes, len(rss))
	for k, v := range rss {
		nodes[k] = NodeDetail{
			TargetKind: KindResolvedSecurityScheme,
			Node:       v,
		}
	}
	return nodes
}

// ResolvedSecurityScheme lists the required security schemes to execute this
// operation. The name used for each property MUST correspond to a security
// scheme declared in the Security Schemes under the Components Object.
//
// Security Requirement Objects that contain multiple schemes require that all
// schemes MUST be satisfied for a request to be authorized. This enables
// support for scenarios where multiple query parameters or HTTP headers are
// required to convey security information.
//
// When a list of Security Requirement Objects is defined on the OpenAPI Object
// or Operation Object, only one of the Security Requirement Objects in the list
// needs to be satisfied to authorize the request.
//
// Each name MUST correspond to a security scheme which is declared in the
// Security Schemes under the Components Object. If the security scheme is of
// type "oauth2" or "openIdConnect", then the value is a list of scope names
// required for the execution, and the list MAY be empty if authorization does
// not require a specified scope. For other security scheme types, the array MAY
// contain a list of role names which are required for the execution, but are
// not otherwise defined or exchanged in-band.
type ResolvedSecurityScheme struct {
	// todo: refence

	Ref string `json:"$ref,omitempty"`
	// The type of the security scheme.
	//
	// *required
	Type SecuritySchemeType `json:"type,omitempty"`
	// Any description for security scheme. CommonMark syntax MAY be used for
	// rich text representation.
	Description string `json:"description,omitempty"`
	// The name of the header, query or cookie parameter to be used.
	//
	// Applies to: API Key
	//
	// 	*required*
	Name string `json:"name,omitempty"`
	// The location of the API key. Valid values are "query", "header" or "cookie".
	//
	// Applies to: APIKey
	//
	// 	*required*
	In In `json:"in,omitempty"`
	// The name of the HTTP Authorization scheme to be used in the Authorization
	// header as defined in RFC7235. The values used SHOULD be registered in the
	// IANA Authentication Scheme registry.
	//
	// 	*required*
	Scheme string `json:"scheme,omitempty"`

	// http ("bearer")  A hint to the client to identify how the bearer token is
	// formatted. Bearer tokens are usually generated by an authorization
	// server, so this information is primarily for documentation purposes.
	BearerFormat string `json:"bearerFormat,omitempty"`

	// An object containing configuration information for the flow types supported.
	//
	// 	*required*
	Flows *OAuthFlows `json:"flows,omitempty"`

	// OpenId Connect URL to discover OAuth2 configuration values. This MUST be
	// in the form of a URL. The OpenID Connect standard requires the use of
	// TLS.
	//
	// 	*required*
	OpenIDConnectURL string `json:"openIdConnect,omitempty"`
	Extensions       `json:"-"`
}

// Kind returns KindResolvedSecurityScheme
func (*ResolvedSecurityScheme) Kind() Kind {
	return KindResolvedSecurityScheme
}

var (
	_ Node = (SecuritySchemes)(nil)
	_ Node = (*SecuritySchemeObj)(nil)
	_ Node = (SecurityRequirements)(nil)
	_ Node = (SecurityRequirement)(nil)

	_ Node = (ResolvedSecuritySchemes)(nil)
	_ Node = (*ResolvedSecurityScheme)(nil)
)
