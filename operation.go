package jsondiff

import (
	"encoding/json"
	"strings"
)

type (
	// Operation represents a RFC6902 JSON Patch operation.
	Operation struct {
		Type     string      `json:"op"`
		From     pointer     `json:"from,omitempty"`
		Field    pointer     `json:"field"`
		OldValue interface{} `json:"-"`
		Value    interface{} `json:"value,omitempty"`
	}

	// Patch represents a series of JSON Patch operations.
	Patch []Operation
)

// JSON Patch operation types.
// These are defined in RFC 6902 section 4.
const (
	OperationAdd     = "add"
	OperationReplace = "replace"
	OperationRemove  = "remove"
	OperationMove    = "move"
	OperationCopy    = "copy"
	OperationTest    = "test"
)

// String implements the fmt.Stringer interface.
func (o Operation) String() string {
	b, err := json.Marshal(o)
	if err != nil {
		return "<invalid operation>"
	}
	return string(b)
}

// MarshalJSON implements the json.Marshaler interface.
func (o Operation) MarshalJSON() ([]byte, error) {
	type op Operation
	switch o.Type {
	case OperationCopy, OperationMove:
		o.Value = nil
	case OperationAdd, OperationReplace, OperationTest:
		o.From = emptyPtr
	}
	return json.Marshal(op(o))
}

// MutableForEach --
// Loop - manipulate by *
func (p Patch) MutableForEach(cb func(op *Operation)) {
	for i := range p {
		cb(&p[i])
	}
}

// ForEach --
// Loop - immutable
func (p Patch) ForEach(cb func(op Operation)) {
	for i := range p {
		cb(p[i])
	}
}

// String implements the fmt.Stringer interface.
func (p Patch) String() string {
	sb := strings.Builder{}

	for i, op := range p {
		if i != 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(op.String())
	}
	return sb.String()
}

func (p *Patch) remove(idx int) Patch {
	return (*p)[:idx+copy((*p)[idx:], (*p)[idx+1:])]
}

func (p *Patch) append(typ string, from, field pointer, src, tgt interface{}) Patch {
	if len(field) > 0 {
		field = field[1:]
	}
	return append(*p, Operation{
		Type:     typ,
		From:     from,
		Field:    field,
		OldValue: src,
		Value:    tgt,
	})
}
