package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/chanced/openapi/yamlutil"
)

// ErrLinkParameterNotFound indicates that a parameter is not found
var ErrLinkParameterNotFound = errors.New("openapi: link parameter not found")

// Link can either be a Link or a Reference
type Link interface {
	Node
	ResolveLink(func(ref string) (*LinkObj, error)) (*LinkObj, error)
}

// Links is a map to hold reusable LinkObjs.
type Links map[string]Link

func (ls *Links) Len() int {
	if ls == nil || *ls == nil {
		return 0
	}
	return len(*ls)
}

func (ls *Links) Get(key string) (Link, bool) {
	if ls == nil || *ls == nil {
		return nil, false
	}
	v, ok := (*ls)[key]
	return v, ok
}

func (ls *Links) Set(key string, val Link) {
	if *ls == nil {
		*ls = Links{
			key: val,
		}
		return
	}
	(*ls)[key] = val
}

func (ls Links) Nodes() Nodes {
	if len(ls) == 0 {
		return nil
	}
	n := make(Nodes, len(ls))
	for k, v := range ls {
		n.maybeAdd(k, v, KindLink)
	}
	return n
}

// Kind returns KindLinks
func (Links) Kind() Kind {
	return KindLinks
}

// UnmarshalJSON unmarshals JSON
func (l *Links) UnmarshalJSON(data []byte) error {
	var dm map[string]json.RawMessage
	if err := json.Unmarshal(data, &dm); err != nil {
		return err
	}
	res := make(Links, len(dm))
	for k, d := range dm {
		if isRefJSON(d) {
			v, err := unmarshalReferenceJSON(d)
			if err != nil {
				return err
			}
			res[k] = v
			continue
		}
		var v link
		if err := unmarshalExtendedJSON(d, &v); err != nil {
			return err
		}
		lv := LinkObj(v)
		res[k] = &lv
	}
	*l = res
	return nil
}

// LinkObj represents a possible design-time link for a response. The presence of a
// link does not guarantee the caller's ability to successfully invoke it,
// rather it provides a known relationship and traversal mechanism between
// responses and other operations.
//
// Unlike dynamic links (i.e. links provided in the response payload), the OAS
// linking mechanism does not require link information in the runtime response.
//
// For computing links, and providing instructions to execute them, a runtime
// expression is used for accessing values in an operation and using them as
// parameters while invoking the linked operation.
type LinkObj struct {
	// A relative or absolute URI reference to an OAS operation. This field is
	// mutually exclusive of the operationId field, and MUST point to an
	// Operation Object. Relative operationRef values MAY be used to locate an
	// existing Operation Object in the OpenAPI definition. See the rules for
	// resolving Relative References.
	OperationRef string `json:"operationRef,omitempty"`
	// The name of an existing, resolvable OAS operation, as defined with a
	// unique operationId. This field is mutually exclusive of the operationRef
	// field.
	OperationID string `json:"operationId,omitempty"`
	// A map representing parameters to pass to an operation as specified with
	// operationId or identified via operationRef. The key is the parameter name
	// to be used, whereas the value can be a constant or an expression to be
	// evaluated and passed to the linked operation. The parameter name can be
	// qualified using the parameter location [{in}.]{name} for operations that
	// use the same parameter name in different locations (e.g. path.id).
	Parameters LinkParameters `json:"parameters,omitempty"`
	// A literal value or {expression} to use as a request body when calling the
	// target operation.
	RequestBody json.RawMessage `json:"requestBody,omitempty"`
	// A description of the link. CommonMark syntax MAY be used for rich text
	// representation.
	Description string `json:"description,omitempty"`
	Extensions  `json:"-"`
}
type link LinkObj

func (LinkObj) Nodes() Nodes { return nil }

// MarshalJSON marshals JSON
func (l LinkObj) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(link(l))
}

// UnmarshalJSON unmarshals JSON
func (l *LinkObj) UnmarshalJSON(data []byte) error {
	var lv link
	if err := unmarshalExtendedJSON(data, &lv); err != nil {
		return err
	}
	*l = LinkObj(lv)
	return nil
}

// Kind returns KindLink
func (*LinkObj) Kind() Kind {
	return KindLink
}

// MarshalYAML marshals YAML
func (l LinkObj) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(l)
}

// UnmarshalYAML unmarshals YAML
func (l *LinkObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, l)
}

// DecodeRequestBody decodes l.RequestBody into dst
//
// dst should be a pointer to a concrete type
func (l *LinkObj) DecodeRequestBody(dst interface{}) error {
	return json.Unmarshal(l.RequestBody, dst)
}

// ResolveLink resolves LinkObj by returning itself. resolve is  not called.
func (l *LinkObj) ResolveLink(func(ref string) (*LinkObj, error)) (*LinkObj, error) {
	return l, nil
}

// LinkParameters is a map representing parameters to pass to an operation as
// specified with operationId or identified via operationRef. The key is the
// parameter name to be used, whereas the value can be a constant or an
// expression to be evaluated and passed to the linked operation. The parameter
// name can be qualified using the parameter location [{in}.]{name} for
// operations that use the same parameter name in different locations (e.g.
// path.id).
type LinkParameters map[string]json.RawMessage

// Decode decodes all parameters into dst
func (lp LinkParameters) Decode(dst interface{}) error {
	data, err := json.Marshal(lp)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

// DecodeParameter decodes a given parameter by name.
// It returns an error if the parameter can not be found.
func (lp LinkParameters) DecodeParameter(key string, dst interface{}) error {
	if v, ok := lp[key]; ok {
		if reflect.TypeOf(dst).Kind() == reflect.Ptr {
			return json.Unmarshal(v, dst)
		}
		return json.Unmarshal(v, &dst)
	}
	return fmt.Errorf("%w {%s}", ErrLinkParameterNotFound, key)
}

// Has returns true if key exists in lp
func (lp LinkParameters) Has(key string) bool {
	_, exists := lp[key]
	return exists
}

// Set concrete object to lp. To add JSON, use SetEncoded
func (lp *LinkParameters) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	lp.SetEncoded(key, data)
	return nil
}

// SetEncoded sets the value of key to value.
//
// Value should be a json encoded byte slice
func (lp *LinkParameters) SetEncoded(key string, value []byte) {
	(*lp)[key] = json.RawMessage(value)
}

// ResolvedLinks is a map to hold reusable LinkObjs.
type ResolvedLinks map[string]*ResolvedLink

func (rls *ResolvedLinks) Len() int {
	if rls == nil || *rls == nil {
		return 0
	}
	return len(*rls)
}

func (rls *ResolvedLinks) Get(key string) (*ResolvedLink, bool) {
	if rls == nil || *rls == nil {
		return nil, false
	}
	v, ok := (*rls)[key]
	return v, ok
}

func (rls *ResolvedLinks) Set(key string, val *ResolvedLink) {
	if *rls == nil {
		*rls = ResolvedLinks{
			key: val,
		}
		return
	}
	(*rls)[key] = val
}

func (rls ResolvedLinks) Nodes() Nodes {
	if len(rls) == 0 {
		return nil
	}
	n := make(Nodes, len(rls))
	for k, v := range rls {
		n.maybeAdd(k, v, KindResolvedLink)
	}
	return n
}

// Kind returns KindResolvedLinks
func (ResolvedLinks) Kind() Kind {
	return KindResolvedLinks
}

// ResolvedLink represents a resolved Link Object which is a possible
// design-time link for a response. The presence of a link does not guarantee
// the caller's ability to successfully invoke it, rather it provides a known
// relationship and traversal mechanism between responses and other operations.
//
// Unlike dynamic links (i.e. links provided in the response payload), the OAS
// linking mechanism does not require link information in the runtime response.
//
// For computing links, and providing instructions to execute them, a runtime
// expression is used for accessing values in an operation and using them as
// parameters while invoking the linked operation.
type ResolvedLink struct {

	// TODO: reference

	// A relative or absolute URI reference to an OAS operation. This field is
	// mutually exclusive of the operationId field, and MUST point to an
	// Operation Object. Relative operationRef values MAY be used to locate an
	// existing Operation Object in the OpenAPI definition. See the rules for
	// resolving Relative References.
	OperationRef string `json:"operationRef,omitempty"`
	// The name of an existing, resolvable OAS operation, as defined with a
	// unique operationId. This field is mutually exclusive of the operationRef
	// field.
	OperationID string `json:"operationId,omitempty"`
	// A map representing parameters to pass to an operation as specified with
	// operationId or identified via operationRef. The key is the parameter name
	// to be used, whereas the value can be a constant or an expression to be
	// evaluated and passed to the linked operation. The parameter name can be
	// qualified using the parameter location [{in}.]{name} for operations that
	// use the same parameter name in different locations (e.g. path.id).
	Parameters LinkParameters `json:"parameters,omitempty"`
	// A literal value or {expression} to use as a request body when calling the
	// target operation.
	RequestBody json.RawMessage `json:"requestBody,omitempty"`
	// A description of the link. CommonMark syntax MAY be used for rich text
	// representation.
	Description string `json:"description,omitempty"`
	Extensions  `json:"-"`
}

func (ResolvedLink) Nodes() Nodes { return nil }

// Kind returns KindResolvedLink
func (*ResolvedLink) Kind() Kind {
	return KindResolvedLink
}

var (
	_ Node = (*LinkObj)(nil)
	_ Node = (*ResolvedLink)(nil)
	_ Node = (Links)(nil)
	_ Node = (ResolvedLinks)(nil)
)
