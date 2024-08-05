package eval

import "github.com/manuelpepe/interpreter/object"

var builtins = map[string]*object.Builtin{
	"len": {Fn: &LenBuiltin{}},
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
