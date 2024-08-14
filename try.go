package try

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
)

// OK panics when err is not null.
func OK(err error) {
	checkErr(err)
}

// Check panics when err is not null.
func Check(err error) {
	checkErr(err)
}

// Val returns v or panics when err is not null.
func Val[T any](v T, err error) T {
	checkErr(err)
	return v
}

// Val2 returns v1, v2 or panics when err is not null.
func Val2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	checkErr(err)
	return v1, v2
}

// Val3 returns v1, v2, v3 or panics when err is not null.
func Val3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	checkErr(err)
	return v1, v2, v3
}

// SafeVal returns v and ignores error.
func SafeVal[T any](v T, err error) T {
	// ignore error
	return v
}

// SafeVal2 returns v1, v2 and ignores error.
func SafeVal2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	// ignore error
	return v1, v2
}

// SafeVal3 returns v1, v2, v3 and ignores error.
func SafeVal3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	// ignore error
	return v1, v2, v3
}

// Require panics if statement is false.
func Require(statement bool, err any) {
	if !statement {
		checkErr(toError(err))
	}
}

// Catch recovers and returns error by argument pointer.
func Catch(err *error) {
	if r := recover(); r != nil && err != nil {
		if *err == nil {
			*err = toError(r)
		} else {
			*err = errors.Join(*err, toError(r))
		}
	}
}

// Mute mutes panic-error.
func Mute() {
	recover()
}

// Call runs the function safely, recovers panic-error.
func Call(fn func()) (err error) {
	defer Catch(&err)
	fn()
	return
}

// Go runs the function safely.
func Go(fn func()) {
	go Call(fn)
}

// Async asynchronously runs several functions and waits for them to complete, returns an error in case of panic.
func Async(fn ...func()) (err error) {
	var wg sync.WaitGroup
	wg.Add(len(fn))
	for _, f := range fn {
		go func(fn func()) {
			defer wg.Done()
			defer Catch(&err)
			fn()
		}(f)
	}
	wg.Wait()
	return
}

func toError(err any) error {
	if e, ok := err.(error); ok {
		return e
	}
	return fmt.Errorf("%v", err)
}

func checkErr(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(2)
		panic(fmt.Errorf("%w\n\t%s:%d", err, file, line))
	}
}
