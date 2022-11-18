package code

import (
	"errors"
	"fmt"
	"testing"
)

type wrapped struct {
	msg string
	err error
}

func (e wrapped) Error() string { return e.msg }

func (e wrapped) Unwrap() error { return e.err }

//采用组合形式相当于实现了Unwrap()
type inheritImplUnwrap struct {
	wrapped
}

func (u inheritImplUnwrap) Error() string {
	return u.err.Error()
}

//内嵌模式当成一个成员所以不能调用里面的Unwrap()方法
type unImplUnwrap struct {
	w wrapped
}

func (u unImplUnwrap) Error() string {
	return u.w.err.Error()
}

func TestCode_Code(t *testing.T) {
	err := func() error {
		return OK
	}
	if err().Error() != "OK" {
		t.Fatalf("want 'OK' real %s\n", err().Error())
	}
}

func TestCause(t *testing.T) {
	var (
		err1 = inheritImplUnwrap{wrapped{"wrap 1", OK}}
		err2 = unImplUnwrap{wrapped{"wrap 1.1", OK}}
		err3 = fmt.Errorf("wrap2:%w", OK)
		err4 = fmt.Errorf("wrap3:%s", OK)
	)

	if _, ok := Cause(err1).(Codes); ok == false {
		t.Error("err1 is not Codes type")
	}

	if _, ok := Cause(err2).(Codes); ok == true {
		t.Error("err2 is not Codes type")
	}

	if _, ok := Cause(err3).(Codes); ok == false {
		t.Error("err3 is not Codes type")
	}

	if _, ok := Cause(err4).(Codes); ok == true {
		t.Error("err4 is not Codes type")
	}

	if _, ok := Cause(nil).(Codes); ok == true {
		t.Error("nil is not Codes type")
	}
}

func TestEqual(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		if Equal(OK, OK) != true {
			t.Fatal("want true but is not equal")
		}
	})

	t.Run("Equal nil", func(t *testing.T) {
		if Equal(nil, nil) != true {
			t.Fatal("want true but is not equal")
		}
	})
}

func TestEqualErr(t *testing.T) {
	var (
		err    = OK
		wrap1  = wrapped{"wrap 1", err}
		tables = []struct {
			seq  int
			err  error
			code Codes
			want bool
		}{
			{1, nil, err, false},
			{2, wrapped{"wrapped", nil}, err, false},
			{3, err, ServeErr, false},
			{4, wrap1, err, true},
			{5, wrapped{"wrap 2", wrap1}, err, true},
			{6, err, err, true},
			{7, errors.New("new err"), err, false},
			{8, fmt.Errorf("fmt err"), err, false},
			{9, fmt.Errorf("fmt err:%w", err), err, true},
			{10, fmt.Errorf("fmt err:%s", err), err, false},
			{11, nil, nil, true},
		}
	)
	for _, tc := range tables {
		if got := EqualErr(tc.err, tc.code); got != tc.want {
			t.Errorf("Seq(%d) Unwrap(%v) = %v, want %v\n", tc.seq, tc.err, got, tc.want)
		}
	}
}
