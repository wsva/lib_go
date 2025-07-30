package stack

import (
	"errors"
	"fmt"
)

type Stack struct {
	Element []interface{}
}

func (stack *Stack) Push(value interface{}) {
	stack.Element = append(stack.Element, value)
}

func (stack *Stack) Top() interface{} {
	if stack.Size() > 0 {
		return stack.Element[stack.Size()-1]
	}
	return nil
}

func (stack *Stack) Pop() error {
	if stack.Size() > 0 {
		stack.Element = stack.Element[:stack.Size()-1]
		return nil
	}
	return errors.New("empty stack")
}

func (stack *Stack) Swap(other *Stack) {
	switch {
	case stack.Size() == 0 && other.Size() == 0:
		return
	case other.Size() == 0:
		other.Element = stack.Element[:stack.Size()]
		stack.Element = nil
	case stack.Size() == 0:
		stack.Element = other.Element
		other.Element = nil
	default:
		stack.Element, other.Element = other.Element, stack.Element
	}
}

func (stack *Stack) Set(idx int, value interface{}) error {
	if idx >= 0 && stack.Size() > 0 && stack.Size() > idx {
		stack.Element[idx] = value
		return nil
	}
	return errors.New("Set failed")
}

func (stack *Stack) Get(idx int) interface{} {
	if idx >= 0 && stack.Size() > 0 && stack.Size() > idx {
		return stack.Element[idx]
	}
	return nil
}

func (stack *Stack) Size() int {
	return len(stack.Element)
}

func (stack *Stack) Empty() bool {
	if stack.Element == nil || stack.Size() == 0 {
		return true
	}
	return false
}

func (stack *Stack) Print() {
	for i := len(stack.Element) - 1; i >= 0; i-- {
		fmt.Println(i, "=>", stack.Element[i])
	}
}
