package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"monkey_lang/ast"
)

const (
	// IntegerObj is an string representation for integer type
	IntegerObj = "INTEGER"

	// FLoatObj is an string representation for float type
	FloatObj = "FLOAT"

	// BooleanObj is an string representation for boolean type
	BooleanObj = "BOOLEAN"
	// StringObj is an string representation for string type
	StringObj = "STRING"

	// NullObj is an string representation for null type
	NullObj = "NULL"
	// ReturnValueObj is an string representation for return value type
	ReturnValueObj = "RETURN_VALUE"

	// ErrorObj is an string representation for object type
	ErrorObj = "ERROR"
	// FunctionObj is an string representation for function type
	FunctionObj = "FUNCTION"

	// BuiltinObj is an string representation for built in type
	BuiltinObj = "BUILTIN"

	// ArrayObj is an string representation for built in array type
	ArrayObj = "ARRAY"

	// HashObj is an string representation for built in hash map type
	HashObj = "HASH"
)

type (
	ObjectType      string
	BuiltinFunction func(args ...Object) Object
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

type Null struct {
	Value bool
}

type Integer struct {
	Value int64
}

type Float struct {
	Value float64
}

type Boolean struct {
	Value bool
}

type ReturnValue struct {
	Value Object
}

type Error struct {
	Message string
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

type String struct {
	Value string
}

type Builtin struct {
	Fn BuiltinFunction
}

type Array struct {
	Elements []Object
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return IntegerObj }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%.6f", f.Value) }
func (f *Float) Type() ObjectType { return FloatObj }
func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: uint64(f.Value)}
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BooleanObj }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NullObj }

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return ReturnValueObj }

func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ErrorObj }

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (f *Function) Type() ObjectType { return FunctionObj }

func (s *String) Type() ObjectType { return StringObj }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (b *Builtin) Type() ObjectType { return BuiltinObj }
func (b *Builtin) Inspect() string  { return "builtin function" }

func (ao *Array) Type() ObjectType { return ArrayObj }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (h *Hash) Type() ObjectType { return HashObj }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
