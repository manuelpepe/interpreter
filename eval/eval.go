package eval

import (
	"fmt"

	"github.com/manuelpepe/interpreter/ast"
	"github.com/manuelpepe/interpreter/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.String(), val)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return evalBoolean(node.Value)
	case *ast.ArrayLiteral:
		items := evalExpressionList(node.Items, env)
		if len(items) == 1 && isError(items[0]) {
			return items[0]
		}
		return &object.Array{Items: items}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.IndexExpression:
		arr := Eval(node.Left, env)
		if isError(arr) {
			return arr
		}
		ix := Eval(node.Index, env)
		if isError(ix) {
			return ix
		}
		return evalIndexExpression(arr, ix)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node.Condition, node.Consequence, node.Alternative, env)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		args := evalExpressionList(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		return applyFunction(fn, args)

	}
	return nil
}

func evalIndexExpression(arr object.Object, ix object.Object) object.Object {
	switch cArr := arr.(type) {
	case *object.Array:
		cIx, ok := ix.(*object.Integer)
		if !ok {
			return newError("expected integer, got %s", ix.Type())
		}
		if cIx.Value < 0 || int(cIx.Value) >= len(cArr.Items) {
			return NULL
		}
		return cArr.Items[cIx.Value]
	case *object.Hash:
		cIx, ok := ix.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", ix.Type())
		}
		val, found := cArr.Pairs[cIx.HashKey()]
		if !found {
			return NULL
		}
		return val.Value
	default:
		return newError("index operator not supported: %s", arr.Type())

	}

}

func evalExpressionList(lst []ast.Expression, env *object.Environment) []object.Object {
	out := make([]object.Object, len(lst))
	for ix, item := range lst {
		argVal := Eval(item, env)
		if isError(argVal) {
			return []object.Object{argVal}
		}
		out[ix] = argVal
	}
	return out
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	hash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
	for keyExpr, valExpr := range node.Items {
		key := Eval(keyExpr, env)
		if isError(key) {
			return key
		}
		hashedKey, ok := key.(object.Hashable)
		if !ok {
			return newError("object is not hashable, got=%T", hash)
		}
		val := Eval(valExpr, env)
		if isError(val) {
			return val
		}
		hash.Pairs[hashedKey.HashKey()] = object.HashPair{
			Key:   key,
			Value: val,
		}
	}
	return hash

}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	funcEnv := fn.Env.Enclose()
	for ix, arg := range args {
		funcEnv.Set(fn.Parameters[ix].Value, arg)
	}
	return funcEnv
}

func unwrapReturnValue(val object.Object) object.Object {
	if returnValue, ok := val.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return val
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Builtin:
		return fn.Fn.Do(args...)
	case *object.Function:
		if len(fn.Parameters) != len(args) {
			return newError("expected %d arguments, got %d", len(fn.Parameters), len(args))
		}
		newEnv := extendFunctionEnv(fn, args)
		ret := Eval(fn.Body, newEnv)
		return unwrapReturnValue(ret)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", node.Value)
}

func evalIfExpression(cond ast.Expression, consequence *ast.BlockStatement, alternative *ast.BlockStatement, env *object.Environment) object.Object {
	res := Eval(cond, env)
	if isError(res) {
		return res
	}
	if isTruthy(res) {
		return Eval(consequence, env)
	} else if alternative != nil {
		return Eval(alternative, env)
	}
	return NULL
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return evalBoolean(left == right)
	case operator == "!=":
		return evalBoolean(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return evalBoolean(leftVal == rightVal)
	case "!=":
		return evalBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "<":
		return evalBoolean(leftVal < rightVal)
	case ">":
		return evalBoolean(leftVal > rightVal)
	case "==":
		return evalBoolean(leftVal == rightVal)
	case "!=":
		return evalBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalProgram(p *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, s := range p.Statements {
		result = Eval(s, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value // unwrap first return value
		case *object.Error:
			return result // stop executing at first error
		}
	}
	return result
}

func evalBlockStatements(b *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, s := range b.Statements {
		result = Eval(s, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalBoolean(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
