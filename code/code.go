package code

import (
	"errors"
	"fmt"
)

var (
	_codes = map[int]string{}
)

func add(e int, msg string) Code {
	if e <= 0 {
		panic("business code must greater than zero")
	}
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = msg
	return Code(e)
}

type Codes interface {
	Code() int
	Error() string
	Detail() string
}

type Code int

// 返回错误码
func (c Code) Code() int { return int(c) }

// 错误描述
func (c Code) Error() string {
	if msg, ok := _codes[c.Code()]; ok {
		return msg
	}
	return "服务器异常错误"
}

func (c Code) Detail() string {
	return ""
}

// 错误拆包
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = cause.Unwrap()
	}
	return err
}

// 根据code码来判断是否是同一个错误
func Equal(a, b Codes) bool {
	if a == nil {
		a = OK
	}
	if b == nil {
		b = OK
	}
	return a.Code() == b.Code()
}

// 错误比较
func EqualErr(err error, target Codes) bool {
	return errors.Is(err, target)
}

// 错误包装
func Wrap(msg string, err error) error {
	return fmt.Errorf(msg+":%w", err)
}