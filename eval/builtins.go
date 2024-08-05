package eval

import (
	"fmt"

	"github.com/manuelpepe/interpreter/object"
)

var builtins = map[string]*object.Builtin{
	"len":     {Fn: &LenBuiltin{}},
	"first":   {Fn: &FirstBuiltin{}},
	"last":    {Fn: &LastBuiltin{}},
	"rest":    {Fn: &RestBuiltin{}},
	"push":    {Fn: &PushBuiltin{}},
	"inspect": {Fn: &InspectBuiltin{}},
}

func checkArgs(n int, args []object.Object) (bool, *object.Error) {
	if len(args) != n {
		return false, newError("wrong number of arguments. got=%d, want=%d", len(args), n)
	}
	return true, nil
}

type LenBuiltin struct{}

func (lb *LenBuiltin) Do(args ...object.Object) object.Object {
	if ok, err := checkArgs(1, args); !ok {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Items))}
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to `len` not supported, got %s",
			args[0].Type())
	}
}
func (lb *LenBuiltin) Name() string {
	return "len"
}

type FirstBuiltin struct{}

func (fb *FirstBuiltin) Do(args ...object.Object) object.Object {
	if ok, err := checkArgs(1, args); !ok {
		return err
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `first` not suported, got %s", args[0].Type())
	}

	arg := args[0].(*object.Array)
	if len(arg.Items) > 0 {
		return arg.Items[0]
	}

	return NULL
}
func (fb *FirstBuiltin) Name() string {
	return "first"
}

type LastBuiltin struct{}

func (fb *LastBuiltin) Do(args ...object.Object) object.Object {
	if ok, err := checkArgs(1, args); !ok {
		return err
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `last` not suported, got %s", args[0].Type())
	}

	arg := args[0].(*object.Array)
	if countItems := len(arg.Items); countItems > 0 {
		return arg.Items[countItems-1]
	}

	return NULL
}
func (fb *LastBuiltin) Name() string {
	return "last"
}

type RestBuiltin struct{}

func (rb *RestBuiltin) Do(args ...object.Object) object.Object {
	if ok, err := checkArgs(1, args); !ok {
		return err
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `rest` not suported, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	if length := len(arr.Items); length > 0 {
		newItems := make([]object.Object, length-1)
		copy(newItems, arr.Items[1:length])
		return &object.Array{Items: newItems}
	}

	return NULL
}
func (rb *RestBuiltin) Name() string {
	return "rest"
}

type PushBuiltin struct{}

func (pb *PushBuiltin) Do(args ...object.Object) object.Object {
	if ok, err := checkArgs(2, args); !ok {
		return err
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `push` not suported, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Items)
	newItems := make([]object.Object, length+1)
	copy(newItems, arr.Items)
	newItems[length] = args[1]

	return &object.Array{Items: newItems}
}
func (pb *PushBuiltin) Name() string {
	return "push"
}

type InspectBuiltin struct{}

func (ib *InspectBuiltin) Do(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Printf("%s\n", arg.Inspect())
	}
	return NULL
}

func (ib *InspectBuiltin) Name() string {
	return "inspect"
}
