package jsondiff

import (
	"strconv"
	"strings"
)

type (
	// pointer represents a RFC6901 JSON Pointer.
	pointer string

	jsonNode struct {
		ptr pointer
		val interface{}
	}
)

const (
	emptyPtr = pointer("")
)

var (
	separator = "/"
	// rfc6901Replacer is a replacer used to escape JSON
	// pointer strings in compliance with the JavaScript
	// Object Notation Pointer syntax.
	// https://tools.ietf.org/html/rfc6901
	rfc6901Replacer = strings.NewReplacer("~", "~0", "/", "~1")
)

func UpdateSeparator(newSeparator string) {
	separator = newSeparator
}

// String implements the fmt.Stringer interface.
func (p pointer) String() string {
	return string(p)
}

func (p pointer) appendKey(key string) pointer {
	return pointer(p.String() + separator + rfc6901Replacer.Replace(key))
}

func (p pointer) appendIndex(idx int) pointer {
	return pointer(p.String() + separator + strconv.Itoa(idx))
}

func (p pointer) isRoot() bool {
	return len(p) == 0
}
