package require

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/d4l3k/messagediff"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/protobuf/proto"

	"github.com/photon-storage/go-common/encoding/ssz"
)

type TestingTB interface {
	Fatalf(format string, args ...interface{})
}

// Fail unconditionally.
func Fail(tb TestingTB, msgs ...interface{}) {
	errMsg := parseMsg("Unexpected failure", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s", filepath.Base(file), line, errMsg)
}

// Panic asserts value is true.
func Panic(tb TestingTB, f func(), msgs ...interface{}) {
	didPanic, panicValue := didPanic(f)
	if !didPanic {
		tb.Fatalf(fmt.Sprintf("func should panic\n\tPanic value:\t%#v", panicValue), msgs...)
	}
}

// True asserts value is true.
func True(tb TestingTB, actual bool, msgs ...interface{}) {
	if actual == true {
		return
	}
	errMsg := parseMsg("Values are not equal", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s, want: %[4]v (%[4]T), got: %[5]v (%[5]T)",
		filepath.Base(file), line, errMsg, true, actual)
}

// False asserts value is false.
func False(tb TestingTB, actual bool, msgs ...interface{}) {
	if actual == false {
		return
	}
	errMsg := parseMsg("Values are not equal", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s, want: %[4]v (%[4]T), got: %[5]v (%[5]T)",
		filepath.Base(file), line, errMsg, false, actual)
}

// Equal compares values using comparison operator.
func Equal(
	tb TestingTB,
	expected interface{},
	actual interface{},
	msgs ...interface{},
) {
	if expected == actual {
		return
	}
	errMsg := parseMsg("Values are not equal", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s, want: %[4]v (%[4]T), got: %[5]v (%[5]T)",
		filepath.Base(file), line, errMsg, expected, actual)
}

// NotEqual compares values using comparison operator.
func NotEqual(
	tb TestingTB,
	expected interface{},
	actual interface{},
	msgs ...interface{},
) {
	if expected != actual {
		return
	}
	errMsg := parseMsg("Values are equal", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s, both values are equal: %[4]v (%[4]T)",
		filepath.Base(file), line, errMsg, expected)
}

// DeepEqual compares values using DeepEqual.
// NOTE: this function does not work for checking arrays/slices or maps of
// protobuf messages. For arrays/slices, please use DeepSSZEqual.
// For maps, please iterate through and compare the individual keys and values.
func DeepEqual(
	tb TestingTB,
	expected interface{},
	actual interface{},
	msgs ...interface{},
) {
	if isDeepEqual(expected, actual) {
		return
	}

	errMsg := parseMsg("Values are not equal", msgs...)
	_, file, line, _ := runtime.Caller(1)
	diff, _ := messagediff.PrettyDiff(expected, actual)
	tb.Fatalf("%s:%d %s, want: %#v, got: %#v, diff: %s",
		filepath.Base(file), line, errMsg, expected, actual, diff)
}

// DeepNotEqual compares values using DeepEqual.
// NOTE: this function does not work for checking arrays/slices or maps of
// protobuf messages. For arrays/slices, please use DeepNotSSZEqual.
// For maps, please iterate through and compare the individual keys and values.
func DeepNotEqual(
	tb TestingTB,
	expected interface{},
	actual interface{},
	msgs ...interface{},
) {
	if !isDeepEqual(expected, actual) {
		return
	}

	errMsg := parseMsg("Values are equal", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s, want: %#v, got: %#v",
		filepath.Base(file), line, errMsg, expected, actual)
}

// DeepSSZEqual compares values using DeepEqual.
func DeepSSZEqual(
	tb TestingTB,
	expected interface{},
	actual interface{},
	msgs ...interface{},
) {
	if ssz.DeepEqual(expected, actual) {
		return
	}

	errMsg := parseMsg("Values are not equal", msgs...)
	_, file, line, _ := runtime.Caller(1)
	diff, _ := messagediff.PrettyDiff(expected, actual)
	tb.Fatalf("%s:%d %s, want: %#v, got: %#v, diff: %s",
		filepath.Base(file), line, errMsg, expected, actual, diff)
}

// DeepNotSSZEqual compares values using DeepEqual.
func DeepNotSSZEqual(
	tb TestingTB,
	expected interface{},
	actual interface{},
	msgs ...interface{},
) {
	if !ssz.DeepEqual(expected, actual) {
		return
	}

	errMsg := parseMsg("Values are equal", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s, want: %#v, got: %#v",
		filepath.Base(file), line, errMsg, expected, actual)
}

// NoError asserts that error is nil.
func NoError(tb TestingTB, err error, msgs ...interface{}) {
	if err == nil {
		return
	}

	errMsg := parseMsg("Unexpected error", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s: %v", filepath.Base(file), line, errMsg, err)
}

// ErrorIs uses Errors.Is to recursively unwrap err looking for target
// in the chain. If any error in the chain matches target, the assertion
// will pass.
func ErrorIs(
	tb TestingTB,
	target error,
	err error,
	msgs ...interface{},
) {
	if errors.Is(err, target) {
		return
	}

	errMsg := parseMsg(fmt.Sprintf("error %s not in chain", target), msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s: %v", filepath.Base(file), line, errMsg, err)
}

// ErrorContains asserts that actual error contains wanted message.
func ErrorContains(
	tb TestingTB,
	want string,
	err error,
	msgs ...interface{},
) {
	if err != nil && strings.Contains(err.Error(), want) {
		return
	}

	errMsg := parseMsg("Expected error not returned", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s, got: %v, want: %s",
		filepath.Base(file), line, errMsg, err, want)
}

// Nil asserts that passed value is nil.
func Nil(tb TestingTB, obj interface{}, msgs ...interface{}) {
	if isNil(obj) {
		return
	}

	errMsg := parseMsg("Unexpected non-nil value", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s, got %v", filepath.Base(file), line, obj, errMsg)
}

// NotNil asserts that passed value is not nil.
func NotNil(tb TestingTB, obj interface{}, msgs ...interface{}) {
	if !isNil(obj) {
		return
	}

	errMsg := parseMsg("Unexpected nil value", msgs...)
	_, file, line, _ := runtime.Caller(1)
	tb.Fatalf("%s:%d %s", filepath.Base(file), line, errMsg)
}

// NotEmpty checks that the object fields are not empty. This method also checks all of the
// pointer fields to ensure none of those fields are empty.
func NotEmpty(tb TestingTB, obj interface{}, msgs ...interface{}) {
	_, isProto := obj.(proto.Message)
	notEmpty(tb, obj, isProto, []string{} /*fields*/, 0 /*stackSize*/, msgs...)
}

// LogsContain checks that the desired string is a subset of the current log output.
func LogsContain(
	tb TestingTB,
	hook *test.Hook,
	want string,
	msgs ...interface{},
) {
	logsContain(tb, hook, want, true, msgs...)
}

// LogsDoNotContain is the inverse check of LogsContain.
func LogsDoNotContain(
	tb TestingTB,
	hook *test.Hook,
	want string,
	msgs ...interface{},
) {
	logsContain(tb, hook, want, false, msgs...)
}

func parseMsg(defaultMsg string, msgs ...interface{}) string {
	if len(msgs) >= 1 {
		msgFormat, ok := msgs[0].(string)
		if !ok {
			return defaultMsg
		}
		return fmt.Sprintf(msgFormat, msgs[1:]...)
	}
	return defaultMsg
}

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f func()) (bool, interface{}) {

	didPanic := false
	var message interface{}
	func() {
		defer func() {
			if message = recover(); message != nil {
				didPanic = true
			}
		}()

		// call the target function
		f()
	}()

	return didPanic, message
}

func isDeepEqual(expected, actual interface{}) bool {
	_, isProto := expected.(proto.Message)
	if isProto {
		return proto.Equal(expected.(proto.Message), actual.(proto.Message))
	}
	return reflect.DeepEqual(expected, actual)
}

// isNil checks that underlying value of obj is nil.
func isNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return value.IsNil()
	}
	return false
}

// notEmpty checks all fields are not zero, including pointer field references to other structs.
// This method has the option to ignore fields without struct tags, which is helpful for checking
// protobuf messages that have internal fields.
func notEmpty(
	tb TestingTB,
	obj interface{},
	ignoreFieldsWithoutTags bool,
	fields []string,
	stackSize int,
	msgs ...interface{},
) {
	var v reflect.Value
	if vo, ok := obj.(reflect.Value); ok {
		v = reflect.Indirect(vo)
	} else {
		v = reflect.Indirect(reflect.ValueOf(obj))
	}

	if len(fields) == 0 {
		fields = []string{v.Type().Name()}
	}

	fail := func(fields []string) {
		m := parseMsg("", msgs...)
		errMsg := fmt.Sprintf("empty/zero field: %s", strings.Join(fields, "."))
		if len(m) > 0 {
			m = strings.Join([]string{m, errMsg}, ": ")
		} else {
			m = errMsg
		}
		_, file, line, _ := runtime.Caller(4 + stackSize)
		tb.Fatalf("%s:%d %s", filepath.Base(file), line, m)
	}

	if v.Kind() != reflect.Struct {
		if v.IsZero() {
			fail(fields)
		}
		return
	}

	for i := 0; i < v.NumField(); i++ {
		if ignoreFieldsWithoutTags && len(v.Type().Field(i).Tag) == 0 {
			continue
		}
		fields := append(fields, v.Type().Field(i).Name)

		switch k := v.Field(i).Kind(); k {
		case reflect.Ptr:
			notEmpty(tb, v.Field(i), ignoreFieldsWithoutTags, fields, stackSize+1, msgs...)
		case reflect.Slice:
			f := v.Field(i)
			if f.Len() == 0 {
				fail(fields)
			}
			for i := 0; i < f.Len(); i++ {
				notEmpty(tb, f.Index(i), ignoreFieldsWithoutTags, fields, stackSize+1, msgs...)
			}
		default:
			if v.Field(i).IsZero() {
				fail(fields)
			}
		}
	}
}

// logsContain checks whether a given substring is a part of logs. If flag=false, inverse is checked.
func logsContain(
	tb TestingTB,
	hook *test.Hook,
	want string,
	flag bool,
	msgs ...interface{},
) {
	_, file, line, _ := runtime.Caller(2)
	entries := hook.AllEntries()
	var logs []string
	match := false
	for _, e := range entries {
		msg, err := e.String()
		if err != nil {
			tb.Fatalf("%s:%d Failed to format log entry to string: %v", filepath.Base(file), line, err)
			return
		}
		if strings.Contains(msg, want) {
			match = true
		}
		for _, field := range e.Data {
			fieldStr, ok := field.(string)
			if !ok {
				continue
			}
			if strings.Contains(fieldStr, want) {
				match = true
			}
		}
		logs = append(logs, msg)
	}
	var errMsg string
	if flag && !match {
		errMsg = parseMsg("Expected log not found", msgs...)
	} else if !flag && match {
		errMsg = parseMsg("Unexpected log found", msgs...)
	}
	if errMsg != "" {
		tb.Fatalf("%s:%d %s: %v\nSearched logs:\n%v", filepath.Base(file), line, errMsg, want, logs)
	}
}
