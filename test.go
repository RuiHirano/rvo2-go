package main

import (
	"fmt"
	"time"
)

type fooOptions struct {
	timeout time.Duration
	bar     string
}

type FooOption func(*fooOptions)

func WithTimeout(timeout time.Duration) FooOption {
	return func(ops *fooOptions) {
		ops.timeout = timeout
	}
}

func WithBar(bar string) FooOption {

	return func(ops *fooOptions) {

		fmt.Printf("asdf", bar)
		ops.bar = bar
	}
}

func Foo(arg1 string, options ...FooOption) {
	opt := fooOptions{}
	for _, o := range options {
		o(&opt)
	}

	// use opt
	fmt.Printf("opt: %v", opt.bar == "")
}

func main() {
	Foo("baz", WithTimeout(time.Second))
}
