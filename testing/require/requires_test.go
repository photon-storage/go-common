package require_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbc "github.com/photon-storage/photon-proto/consensus"
	pbtest "github.com/photon-storage/photon-proto/util"

	"github.com/photon-storage/go-common/testing/require"
)

func Test_Equal(t *testing.T) {
	type args struct {
		tb       *tbMock
		expected interface{}
		actual   interface{}
		msgs     []interface{}
	}
	tests := []struct {
		name        string
		args        args
		expectedErr string
	}{
		{
			name: "equal values",
			args: args{
				tb:       &tbMock{},
				expected: 42,
				actual:   42,
			},
		},
		{
			name: "equal values different types",
			args: args{
				tb:       &tbMock{},
				expected: uint64(42),
				actual:   42,
			},
			expectedErr: "Values are not equal, want: 42 (uint64), got: 42 (int)",
		},
		{
			name: "non-equal values",
			args: args{
				tb:       &tbMock{},
				expected: 42,
				actual:   41,
			},
			expectedErr: "Values are not equal, want: 42 (int), got: 41 (int)",
		},
		{
			name: "custom error message",
			args: args{
				tb:       &tbMock{},
				expected: 42,
				actual:   41,
				msgs:     []interface{}{"Custom values are not equal"},
			},
			expectedErr: "Custom values are not equal, want: 42 (int), got: 41 (int)",
		},
		{
			name: "custom error message with params",
			args: args{
				tb:       &tbMock{},
				expected: 42,
				actual:   41,
				msgs:     []interface{}{"Custom values are not equal (for slot %d)", 12},
			},
			expectedErr: "Custom values are not equal (for slot 12), want: 42 (int), got: 41 (int)",
		},
	}
	for _, tt := range tests {
		verify := func() {
			if tt.expectedErr == "" && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %v", tt.args.tb.FatalfMsg)
			} else if !strings.Contains(tt.args.tb.FatalfMsg, tt.expectedErr) {
				t.Errorf("got: %q, want: %q", tt.args.tb.FatalfMsg, tt.expectedErr)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			require.Equal(tt.args.tb, tt.args.expected, tt.args.actual, tt.args.msgs...)
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			require.Equal(tt.args.tb, tt.args.expected, tt.args.actual, tt.args.msgs...)
			verify()
		})
	}
}

func Test_NotEqual(t *testing.T) {
	type args struct {
		tb       *tbMock
		expected interface{}
		actual   interface{}
		msgs     []interface{}
	}
	tests := []struct {
		name        string
		args        args
		expectedErr string
	}{
		{
			name: "equal values",
			args: args{
				tb:       &tbMock{},
				expected: 42,
				actual:   42,
			},
			expectedErr: "Values are equal, both values are equal: 42 (int)",
		},
		{
			name: "equal values different types",
			args: args{
				tb:       &tbMock{},
				expected: uint64(42),
				actual:   42,
			},
		},
		{
			name: "non-equal values",
			args: args{
				tb:       &tbMock{},
				expected: 42,
				actual:   41,
			},
		},
		{
			name: "custom error message",
			args: args{
				tb:       &tbMock{},
				expected: 42,
				actual:   42,
				msgs:     []interface{}{"Custom values are equal"},
			},
			expectedErr: "Custom values are equal, both values are equal",
		},
	}
	for _, tt := range tests {
		verify := func() {
			if tt.expectedErr == "" && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %v", tt.args.tb.FatalfMsg)
			} else if !strings.Contains(tt.args.tb.FatalfMsg, tt.expectedErr) {
				t.Errorf("got: %q, want: %q", tt.args.tb.FatalfMsg, tt.expectedErr)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			require.NotEqual(tt.args.tb, tt.args.expected, tt.args.actual, tt.args.msgs...)
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			require.NotEqual(tt.args.tb, tt.args.expected, tt.args.actual, tt.args.msgs...)
			verify()
		})
	}
}

func TestAssert_DeepEqual(t *testing.T) {
	type args struct {
		tb       *tbMock
		expected interface{}
		actual   interface{}
		msgs     []interface{}
	}
	tests := []struct {
		name        string
		args        args
		expectedErr string
	}{
		{
			name: "equal values",
			args: args{
				tb:       &tbMock{},
				expected: struct{ i int }{42},
				actual:   struct{ i int }{42},
			},
		},
		{
			name: "non-equal values",
			args: args{
				tb:       &tbMock{},
				expected: struct{ i int }{42},
				actual:   struct{ i int }{41},
			},
			expectedErr: "Values are not equal, want: struct { i int }{i:42}, got: struct { i int }{i:41}",
		},
		{
			name: "equal values with proto message",
			args: args{
				tb:       &tbMock{},
				expected: &pbtest.TestMessage{Foo: "foo"},
				actual:   &pbtest.TestMessage{Foo: "foo"},
			},
		},
		{
			name: "equal values with empty proto message",
			args: args{
				tb:       &tbMock{},
				expected: &pbtest.TestMessage{},
				actual:   &pbtest.TestMessage{},
			},
		},
		{
			name: "not-equal values with proto message",
			args: args{
				tb:       &tbMock{},
				expected: &pbtest.TestMessage{Foo: "foo"},
				actual:   &pbtest.TestMessage{Foo: "foo-test"},
				msgs:     []interface{}{"Custom values are not equal"},
			},
			expectedErr: "Custom values are not equal",
		},
		{
			name: "equal values with slice proto message",
			args: args{
				tb: &tbMock{},
				expected: []*pbtest.TestMessage{
					&pbtest.TestMessage{Foo: "foo"},
					&pbtest.TestMessage{Bar: "bar"},
				},
				actual: []*pbtest.TestMessage{
					&pbtest.TestMessage{Foo: "foo"},
					&pbtest.TestMessage{Bar: "bar"},
				},
			},
		},
		{
			name: "equal values with empty slice proto message",
			args: args{
				tb:       &tbMock{},
				expected: []*pbtest.TestMessage{},
				actual:   []*pbtest.TestMessage{},
			},
		},
		{
			name: "not-equal values with slice proto message",
			args: args{
				tb: &tbMock{},
				expected: []*pbtest.TestMessage{
					&pbtest.TestMessage{Foo: "foo"},
					&pbtest.TestMessage{Bar: "bar"},
				},
				actual: []*pbtest.TestMessage{
					&pbtest.TestMessage{Foo: "foo-test"},
					&pbtest.TestMessage{Bar: "bar-test"},
				},
				msgs: []interface{}{"Custom values are not equal"},
			},
			expectedErr: "Custom values are not equal",
		},
		{
			name: "custom error message",
			args: args{
				tb:       &tbMock{},
				expected: struct{ i int }{42},
				actual:   struct{ i int }{41},
				msgs:     []interface{}{"Custom values are not equal"},
			},
			expectedErr: "Custom values are not equal, want: struct { i int }{i:42}, got: struct { i int }{i:41}",
		},
		{
			name: "custom error message with params",
			args: args{
				tb:       &tbMock{},
				expected: struct{ i int }{42},
				actual:   struct{ i int }{41},
				msgs:     []interface{}{"Custom values are not equal (for slot %d)", 12},
			},
			expectedErr: "Custom values are not equal (for slot 12), want: struct { i int }{i:42}, got: struct { i int }{i:41}",
		},
	}
	for _, tt := range tests {
		verify := func() {
			if tt.expectedErr == "" && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %v", tt.args.tb.FatalfMsg)
			} else if !strings.Contains(tt.args.tb.FatalfMsg, tt.expectedErr) {
				t.Errorf("got: %q, want: %q", tt.args.tb.FatalfMsg, tt.expectedErr)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			require.DeepEqual(tt.args.tb, tt.args.expected, tt.args.actual, tt.args.msgs...)
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			require.DeepEqual(tt.args.tb, tt.args.expected, tt.args.actual, tt.args.msgs...)
			verify()
		})
	}
}

func TestAssert_DeepNotEqual(t *testing.T) {
	type args struct {
		tb       *tbMock
		expected interface{}
		actual   interface{}
		msgs     []interface{}
	}
	tests := []struct {
		name        string
		args        args
		expectedErr string
	}{
		{
			name: "equal values",
			args: args{
				tb:       &tbMock{},
				expected: struct{ i int }{42},
				actual:   struct{ i int }{42},
			},
			expectedErr: "Values are equal, want: struct { i int }{i:42}, got: struct { i int }{i:42}",
		},
		{
			name: "non-equal values",
			args: args{
				tb:       &tbMock{},
				expected: struct{ i int }{42},
				actual:   struct{ i int }{41},
			},
		},
		{
			name: "equal values with proto message",
			args: args{
				tb:       &tbMock{},
				expected: &pbtest.TestMessage{Foo: "foo"},
				actual:   &pbtest.TestMessage{Foo: "foo"},
			},
			expectedErr: "Values are equal",
		},
		{
			name: "equal values with empty proto message",
			args: args{
				tb:       &tbMock{},
				expected: &pbtest.TestMessage{},
				actual:   &pbtest.TestMessage{},
			},
			expectedErr: "Values are equal",
		},
		{
			name: "not-equal values with proto message",
			args: args{
				tb:       &tbMock{},
				expected: &pbtest.TestMessage{Foo: "foo"},
				actual:   &pbtest.TestMessage{Foo: "foo-test"},
			},
		},
		{
			name: "equal values with slice proto message",
			args: args{
				tb: &tbMock{},
				expected: []*pbtest.TestMessage{
					&pbtest.TestMessage{Foo: "foo"},
					&pbtest.TestMessage{Bar: "bar"},
				},
				actual: []*pbtest.TestMessage{
					&pbtest.TestMessage{Foo: "foo"},
					&pbtest.TestMessage{Bar: "bar"},
				},
			},
			expectedErr: "Values are equal",
		},
		{
			name: "equal values with empty slice proto message",
			args: args{
				tb:       &tbMock{},
				expected: []*pbtest.TestMessage{},
				actual:   []*pbtest.TestMessage{},
			},
			expectedErr: "Values are equal",
		},
		{
			name: "not-equal values with slice proto message",
			args: args{
				tb: &tbMock{},
				expected: []*pbtest.TestMessage{
					&pbtest.TestMessage{Foo: "foo"},
					&pbtest.TestMessage{Bar: "bar"},
				},
				actual: []*pbtest.TestMessage{
					&pbtest.TestMessage{Foo: "foo-test"},
					&pbtest.TestMessage{Bar: "bar-test"},
				},
			},
		},
		{
			name: "custom error message",
			args: args{
				tb:       &tbMock{},
				expected: struct{ i int }{42},
				actual:   struct{ i int }{42},
				msgs:     []interface{}{"Custom values are equal"},
			},
			expectedErr: "Custom values are equal, want: struct { i int }{i:42}, got: struct { i int }{i:42}",
		},
		{
			name: "custom error message with params",
			args: args{
				tb:       &tbMock{},
				expected: struct{ i int }{42},
				actual:   struct{ i int }{42},
				msgs:     []interface{}{"Custom values are equal (for slot %d)", 12},
			},
			expectedErr: "Custom values are equal (for slot 12), want: struct { i int }{i:42}, got: struct { i int }{i:42}",
		},
	}
	for _, tt := range tests {
		verify := func() {
			if tt.expectedErr == "" && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %v", tt.args.tb.FatalfMsg)
			} else if !strings.Contains(tt.args.tb.FatalfMsg, tt.expectedErr) {
				t.Errorf("got: %q, want: %q", tt.args.tb.FatalfMsg, tt.expectedErr)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			require.DeepNotEqual(tt.args.tb, tt.args.expected, tt.args.actual, tt.args.msgs...)
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			require.DeepNotEqual(tt.args.tb, tt.args.expected, tt.args.actual, tt.args.msgs...)
			verify()
		})
	}
}

func TestAssert_DeepSSZEqual(t *testing.T) {
	type args struct {
		tb       *tbMock
		expected interface{}
		actual   interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult bool
	}{
		{
			name: "equal values",
			args: args{
				tb:       &tbMock{},
				expected: struct{ I uint64 }{42},
				actual:   struct{ I uint64 }{42},
			},
			expectedResult: true,
		},
		{
			name: "equal structs",
			args: args{
				tb: &tbMock{},
				expected: &pbc.Checkpoint{
					Epoch: 5,
				},
				actual: &pbc.Checkpoint{
					Epoch: 5,
				},
			},
			expectedResult: true,
		},
		{
			name: "non-equal values",
			args: args{
				tb:       &tbMock{},
				expected: struct{ I uint64 }{42},
				actual:   struct{ I uint64 }{41},
			},
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		verify := func() {
			if tt.expectedResult && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %s %v", tt.name, tt.args.tb.FatalfMsg)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			require.DeepSSZEqual(tt.args.tb, tt.args.expected, tt.args.actual)
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			require.DeepSSZEqual(tt.args.tb, tt.args.expected, tt.args.actual)
			verify()
		})
	}
}

func TestAssert_DeepNotSSZEqual(t *testing.T) {
	type args struct {
		tb       *tbMock
		expected interface{}
		actual   interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult bool
	}{
		{
			name: "equal values",
			args: args{
				tb:       &tbMock{},
				expected: struct{ I uint64 }{42},
				actual:   struct{ I uint64 }{42},
			},
			expectedResult: true,
		},
		{
			name: "non-equal values",
			args: args{
				tb:       &tbMock{},
				expected: struct{ I uint64 }{42},
				actual:   struct{ I uint64 }{41},
			},
			expectedResult: false,
		},
		{
			name: "not equal structs",
			args: args{
				tb: &tbMock{},
				expected: &pbc.Checkpoint{
					Epoch: 4,
				},
				actual: &pbc.Checkpoint{
					Epoch: 5,
				},
			},
			expectedResult: true,
		},
	}
	for _, tt := range tests {
		verify := func() {
			if !tt.expectedResult && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %s %v", tt.name, tt.args.tb.FatalfMsg)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			require.DeepNotSSZEqual(tt.args.tb, tt.args.expected, tt.args.actual)
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			require.DeepNotSSZEqual(tt.args.tb, tt.args.expected, tt.args.actual)
			verify()
		})
	}
}

func TestAssert_NoError(t *testing.T) {
	type args struct {
		tb   *tbMock
		err  error
		msgs []interface{}
	}
	tests := []struct {
		name        string
		args        args
		expectedErr string
	}{
		{
			name: "nil error",
			args: args{
				tb: &tbMock{},
			},
		},
		{
			name: "non-nil error",
			args: args{
				tb:  &tbMock{},
				err: errors.New("failed"),
			},
			expectedErr: "Unexpected error: failed",
		},
		{
			name: "custom non-nil error",
			args: args{
				tb:   &tbMock{},
				err:  errors.New("failed"),
				msgs: []interface{}{"Custom error message"},
			},
			expectedErr: "Custom error message: failed",
		},
		{
			name: "custom non-nil error with params",
			args: args{
				tb:   &tbMock{},
				err:  errors.New("failed"),
				msgs: []interface{}{"Custom error message (for slot %d)", 12},
			},
			expectedErr: "Custom error message (for slot 12): failed",
		},
	}
	for _, tt := range tests {
		verify := func() {
			if tt.expectedErr == "" && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %v", tt.args.tb.FatalfMsg)
			} else if !strings.Contains(tt.args.tb.FatalfMsg, tt.expectedErr) {
				t.Errorf("got: %q, want: %q", tt.args.tb.FatalfMsg, tt.expectedErr)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			require.NoError(tt.args.tb, tt.args.err, tt.args.msgs...)
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			require.NoError(tt.args.tb, tt.args.err, tt.args.msgs...)
			verify()
		})
	}
}

func TestAssert_ErrorContains(t *testing.T) {
	type args struct {
		tb   *tbMock
		want string
		err  error
		msgs []interface{}
	}
	tests := []struct {
		name        string
		args        args
		expectedErr string
	}{
		{
			name: "nil error",
			args: args{
				tb:   &tbMock{},
				want: "some error",
			},
			expectedErr: "Expected error not returned, got: <nil>, want: some error",
		},
		{
			name: "unexpected error",
			args: args{
				tb:   &tbMock{},
				want: "another error",
				err:  errors.New("failed"),
			},
			expectedErr: "Expected error not returned, got: failed, want: another error",
		},
		{
			name: "expected error",
			args: args{
				tb:   &tbMock{},
				want: "failed",
				err:  errors.New("failed"),
			},
			expectedErr: "",
		},
		{
			name: "custom unexpected error",
			args: args{
				tb:   &tbMock{},
				want: "another error",
				err:  errors.New("failed"),
				msgs: []interface{}{"Something wrong"},
			},
			expectedErr: "Something wrong, got: failed, want: another error",
		},
		{
			name: "expected error",
			args: args{
				tb:   &tbMock{},
				want: "failed",
				err:  errors.New("failed"),
				msgs: []interface{}{"Something wrong"},
			},
			expectedErr: "",
		},
		{
			name: "custom unexpected error with params",
			args: args{
				tb:   &tbMock{},
				want: "another error",
				err:  errors.New("failed"),
				msgs: []interface{}{"Something wrong (for slot %d)", 12},
			},
			expectedErr: "Something wrong (for slot 12), got: failed, want: another error",
		},
		{
			name: "expected error with params",
			args: args{
				tb:   &tbMock{},
				want: "failed",
				err:  errors.New("failed"),
				msgs: []interface{}{"Something wrong (for slot %d)", 12},
			},
			expectedErr: "",
		},
	}
	for _, tt := range tests {
		verify := func() {
			if tt.expectedErr == "" && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %v", tt.args.tb.FatalfMsg)
			} else if !strings.Contains(tt.args.tb.FatalfMsg, tt.expectedErr) {
				t.Errorf("got: %q, want: %q", tt.args.tb.FatalfMsg, tt.expectedErr)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			require.ErrorContains(tt.args.tb, tt.args.want, tt.args.err, tt.args.msgs...)
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			require.ErrorContains(tt.args.tb, tt.args.want, tt.args.err, tt.args.msgs...)
			verify()
		})
	}
}

func Test_NotNil(t *testing.T) {
	type args struct {
		tb   *tbMock
		obj  interface{}
		msgs []interface{}
	}
	var nilBlock *pbc.SignedBlock = nil
	tests := []struct {
		name        string
		args        args
		expectedErr string
	}{
		{
			name: "nil",
			args: args{
				tb: &tbMock{},
			},
			expectedErr: "Unexpected nil value",
		},
		{
			name: "nil custom message",
			args: args{
				tb:   &tbMock{},
				msgs: []interface{}{"This should not be nil"},
			},
			expectedErr: "This should not be nil",
		},
		{
			name: "nil custom message with params",
			args: args{
				tb:   &tbMock{},
				msgs: []interface{}{"This should not be nil (for slot %d)", 12},
			},
			expectedErr: "This should not be nil (for slot 12)",
		},
		{
			name: "not nil",
			args: args{
				tb:  &tbMock{},
				obj: "some value",
			},
			expectedErr: "",
		},
		{
			name: "nil value of dynamic type",
			args: args{
				tb:  &tbMock{},
				obj: nilBlock,
			},
			expectedErr: "Unexpected nil value",
		},
		{
			name: "make sure that assertion works for basic type",
			args: args{
				tb:  &tbMock{},
				obj: 15,
			},
		},
	}
	for _, tt := range tests {
		verify := func() {
			if tt.expectedErr == "" && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %v", tt.args.tb.FatalfMsg)
			} else if !strings.Contains(tt.args.tb.FatalfMsg, tt.expectedErr) {
				t.Errorf("got: %q, want: %q", tt.args.tb.FatalfMsg, tt.expectedErr)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			require.NotNil(tt.args.tb, tt.args.obj, tt.args.msgs...)
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			require.NotNil(tt.args.tb, tt.args.obj, tt.args.msgs...)
			verify()
		})
	}
}

func Test_LogsContainDoNotContain(t *testing.T) {
	type args struct {
		tb   *tbMock
		want string
		flag bool
		msgs []interface{}
	}
	tests := []struct {
		name        string
		args        args
		updateLogs  func(log *logrus.Logger)
		expectedErr string
	}{
		{
			name: "should contain not found",
			args: args{
				tb:   &tbMock{},
				want: "here goes some expected log string",
				flag: true,
			},
			expectedErr: "Expected log not found: here goes some expected log string",
		},
		{
			name: "should contain found",
			args: args{
				tb:   &tbMock{},
				want: "here goes some expected log string",
				flag: true,
			},
			updateLogs: func(log *logrus.Logger) {
				log.Info("here goes some expected log string")
			},
			expectedErr: "",
		},
		{
			name: "should contain not found custom message",
			args: args{
				tb:   &tbMock{},
				msgs: []interface{}{"Waited for logs"},
				want: "here goes some expected log string",
				flag: true,
			},
			expectedErr: "Waited for logs: here goes some expected log string",
		},
		{
			name: "should contain not found custom message with params",
			args: args{
				tb:   &tbMock{},
				msgs: []interface{}{"Waited for %d logs", 10},
				want: "here goes some expected log string",
				flag: true,
			},
			expectedErr: "Waited for 10 logs: here goes some expected log string",
		},
		{
			name: "should not contain and not found",
			args: args{
				tb:   &tbMock{},
				want: "here goes some unexpected log string",
			},
			expectedErr: "",
		},
		{
			name: "should not contain but found",
			args: args{
				tb:   &tbMock{},
				want: "here goes some unexpected log string",
			},
			updateLogs: func(log *logrus.Logger) {
				log.Info("here goes some unexpected log string")
			},
			expectedErr: "Unexpected log found: here goes some unexpected log string",
		},
		{
			name: "should not contain but found custom message",
			args: args{
				tb:   &tbMock{},
				msgs: []interface{}{"Dit not expect logs"},
				want: "here goes some unexpected log string",
			},
			updateLogs: func(log *logrus.Logger) {
				log.Info("here goes some unexpected log string")
			},
			expectedErr: "Dit not expect logs: here goes some unexpected log string",
		},
		{
			name: "should not contain but found custom message with params",
			args: args{
				tb:   &tbMock{},
				msgs: []interface{}{"Dit not expect %d logs", 10},
				want: "here goes some unexpected log string",
			},
			updateLogs: func(log *logrus.Logger) {
				log.Info("here goes some unexpected log string")
			},
			expectedErr: "Dit not expect 10 logs: here goes some unexpected log string",
		},
	}
	for _, tt := range tests {
		verify := func() {
			if tt.expectedErr == "" && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %v", tt.args.tb.FatalfMsg)
			} else if !strings.Contains(tt.args.tb.FatalfMsg, tt.expectedErr) {
				t.Errorf("got: %q, want: %q", tt.args.tb.FatalfMsg, tt.expectedErr)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			log, hook := test.NewNullLogger()
			if tt.updateLogs != nil {
				tt.updateLogs(log)
			}
			if tt.args.flag {
				require.LogsContain(tt.args.tb, hook, tt.args.want, tt.args.msgs...)
			} else {
				require.LogsDoNotContain(tt.args.tb, hook, tt.args.want, tt.args.msgs...)
			}
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			log, hook := test.NewNullLogger()
			if tt.updateLogs != nil {
				tt.updateLogs(log)
			}
			if tt.args.flag {
				require.LogsContain(tt.args.tb, hook, tt.args.want, tt.args.msgs...)
			} else {
				require.LogsDoNotContain(tt.args.tb, hook, tt.args.want, tt.args.msgs...)
			}
			verify()
		})
	}
}

func TestAssert_NotEmpty(t *testing.T) {
	type args struct {
		tb     *tbMock
		input  interface{}
		actual interface{}
		msgs   []interface{}
	}
	tests := []struct {
		name        string
		args        args
		expectedErr string
	}{
		{
			name: "literal value int",
			args: args{
				tb:    &tbMock{},
				input: 42,
			},
		}, {
			name: "literal value int",
			args: args{
				tb:    &tbMock{},
				input: 0,
			},
			expectedErr: "empty/zero field: int",
		}, {
			name: "literal value slice",
			args: args{
				tb:    &tbMock{},
				input: []uint64{42},
			},
		}, {
			name: "literal value string",
			args: args{
				tb:    &tbMock{},
				input: "42",
			},
		}, {
			name: "simple populated struct",
			args: args{
				tb: &tbMock{},
				input: struct {
					foo int
					bar string
				}{
					foo: 42,
					bar: "42",
				},
			},
		}, {
			name: "simple partially empty struct",
			args: args{
				tb: &tbMock{},
				input: struct {
					foo int
					bar string
				}{
					bar: "42",
				},
			},
			expectedErr: "empty/zero field: .foo",
		}, {
			name: "simple empty struct",
			args: args{
				tb: &tbMock{},
				input: struct {
					foo int
					bar string
				}{},
			},
			expectedErr: "empty/zero field",
		}, {
			name: "simple populated protobuf",
			args: args{
				tb: &tbMock{},
				input: &pbtest.Puzzle{
					Challenge: "what do you get when protobufs have internal fields?",
					Answer:    "Complicated reflect logic!",
				},
			},
		}, {
			name: "simple partially empty protobuf",
			args: args{
				tb: &tbMock{},
				input: &pbtest.Puzzle{
					Challenge: "what do you get when protobufs have internal fields?",
				},
			},
			expectedErr: "empty/zero field: Puzzle.Answer",
		}, {
			name: "complex populated protobuf",
			args: args{
				tb: &tbMock{},
				input: &pbtest.AddressBook{
					People: []*pbtest.Person{
						{
							Name:  "Foo",
							Id:    42,
							Email: "foo@bar.com",
							Phones: []*pbtest.Person_PhoneNumber{
								{
									Number: "+1 111-111-1111",
									Type:   pbtest.Person_WORK, // Note: zero'th enum value will count as empty.
								},
							},
							LastUpdated: timestamppb.Now(),
						},
					},
				},
			},
		}, {
			name: "complex partially empty protobuf with empty slices",
			args: args{
				tb: &tbMock{},
				input: &pbtest.AddressBook{
					People: []*pbtest.Person{
						{
							Name:        "Foo",
							Id:          42,
							Email:       "foo@bar.com",
							Phones:      []*pbtest.Person_PhoneNumber{},
							LastUpdated: timestamppb.Now(),
						},
					},
				},
			},
			expectedErr: "empty/zero field: AddressBook.People.Phones",
		}, {
			name: "complex partially empty protobuf with empty string",
			args: args{
				tb: &tbMock{},
				input: &pbtest.AddressBook{
					People: []*pbtest.Person{
						{
							Name:  "Foo",
							Id:    42,
							Email: "",
							Phones: []*pbtest.Person_PhoneNumber{
								{
									Number: "+1 111-111-1111",
									Type:   pbtest.Person_WORK, // Note: zero'th enum value will count as empty.
								},
							},
							LastUpdated: timestamppb.Now(),
						},
					},
				},
			},
			expectedErr: "empty/zero field: AddressBook.People.Email",
		},
	}
	for _, tt := range tests {
		verify := func() {
			if tt.expectedErr == "" && tt.args.tb.FatalfMsg != "" {
				t.Errorf("Unexpected error: %v", tt.args.tb.FatalfMsg)
			} else if !strings.Contains(tt.args.tb.FatalfMsg, tt.expectedErr) {
				t.Errorf("got: %q, want: %q", tt.args.tb.FatalfMsg, tt.expectedErr)
			}
		}
		t.Run(fmt.Sprintf("Assert/%s", tt.name), func(t *testing.T) {
			require.NotEmpty(tt.args.tb, tt.args.input, tt.args.msgs...)
			verify()
		})
		t.Run(fmt.Sprintf("Require/%s", tt.name), func(t *testing.T) {
			require.NotEmpty(tt.args.tb, tt.args.input, tt.args.msgs...)
			verify()
		})
	}
}

// tbMock exposes enough testing.TB methods.
type tbMock struct {
	FatalfMsg string
}

// Errorf writes testing logs to FatalfMsg.
func (tb *tbMock) Errorf(format string, args ...interface{}) {
	panic("not implemented")
}

// Fatalf writes testing logs to FatalfMsg.
func (tb *tbMock) Fatalf(format string, args ...interface{}) {
	tb.FatalfMsg = fmt.Sprintf(format, args...)
}
