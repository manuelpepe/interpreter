package eval

import (
	"github.com/manuelpepe/interpreter/ast"
	"github.com/manuelpepe/interpreter/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatements(node)
	case *ast.ReturnStatement:
		return &object.ReturnValue{Value: Eval(node.ReturnValue)}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return evalBoolean(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node.Condition, node.Consequence, node.Alternative)
	}
	return nil
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

func evalIfExpression(cond ast.Expression, consequence *ast.BlockStatement, alternative *ast.BlockStatement) object.Object {
	res := Eval(cond)
	if isTruthy(res) {
		return Eval(consequence)
	} else if alternative != nil {
		return Eval(alternative)
	}
	return NULL
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return evalBoolean(left == right)
	case operator == "!=":
		return evalBoolean(left != right)
	default:
		return NULL // TODO: Should probably error
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
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return NULL
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
		return NULL // TODO: Should probably error
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalProgram(p *ast.Program) object.Object {
	var result object.Object
	for _, s := range p.Statements {
		result = Eval(s)
		if ret, ok := result.(*object.ReturnValue); ok {
			return ret.Value // unwrap first return value
		}
	}
	return result
}

func evalBlockStatements(b *ast.BlockStatement) object.Object {
	var result object.Object
	for _, s := range b.Statements {
		result = Eval(s)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result // don't unwrap return values in block statements
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
