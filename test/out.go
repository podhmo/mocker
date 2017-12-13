// Code generated by mocker; DO NOT EDIT
// github.com/travisjeffery/mocker
package test

import (
	av1 "github.com/travisjeffery/mocker/test/a"
	bv1 "github.com/travisjeffery/mocker/test/b"
	"github.com/travisjeffery/mocker/test/c"
	"sync"
)

var (
	lockMockIfaceFour  sync.RWMutex
	lockMockIfaceOne   sync.RWMutex
	lockMockIfaceThree sync.RWMutex
	lockMockIfaceTwo   sync.RWMutex
)

// MockIface is a mock implementation of Iface.
//
//     func TestSomethingThatUsesIface(t *testing.T) {
//
//         // make and configure a mocked Iface
//         mockedIface := &MockIface{
//             FourFunc: func(in1 c.Int)  {
// 	               panic("TODO: mock out the Four method")
//             },
//             OneFunc: func(str string,variadic ...string) (string, []string) {
// 	               panic("TODO: mock out the One method")
//             },
//             ThreeFunc: func(in1 av1.Int) bv1.Str {
// 	               panic("TODO: mock out the Three method")
//             },
//             TwoFunc: func(in1 int,in2 int) int {
// 	               panic("TODO: mock out the Two method")
//             },
//         }
//
//         // TODO: use mockedIface in code that requires Iface
//         //       and then make assertions.
//
//     }
type MockIface struct {
	// FourFunc mocks the Four method.
	FourFunc func(in1 c.Int)

	// OneFunc mocks the One method.
	OneFunc func(str string, variadic ...string) (string, []string)

	// ThreeFunc mocks the Three method.
	ThreeFunc func(in1 av1.Int) bv1.Str

	// TwoFunc mocks the Two method.
	TwoFunc func(in1 int, in2 int) int

	// calls tracks calls to the methods.
	calls struct {
		// Four holds details about calls to the Four method.
		Four []struct {
			// In1 is the in1 argument value.
			In1 c.Int
		}
		// One holds details about calls to the One method.
		One []struct {
			// Str is the str argument value.
			Str string
			// Variadic is the variadic argument value.
			Variadic []string
		}
		// Three holds details about calls to the Three method.
		Three []struct {
			// In1 is the in1 argument value.
			In1 av1.Int
		}
		// Two holds details about calls to the Two method.
		Two []struct {
			// In1 is the in1 argument value.
			In1 int
			// In2 is the in2 argument value.
			In2 int
		}
	}
}

// Reset resets the calls made to the mocked APIs.
func (mock *MockIface) Reset() {
	lockMockIfaceFour.Lock()
	mock.calls.Four = nil
	lockMockIfaceFour.Unlock()
	lockMockIfaceOne.Lock()
	mock.calls.One = nil
	lockMockIfaceOne.Unlock()
	lockMockIfaceThree.Lock()
	mock.calls.Three = nil
	lockMockIfaceThree.Unlock()
	lockMockIfaceTwo.Lock()
	mock.calls.Two = nil
	lockMockIfaceTwo.Unlock()
}

// Four calls FourFunc.
func (mock *MockIface) Four(in1 c.Int) {
	if mock.FourFunc == nil {
		panic("moq: MockIface.FourFunc is nil but Iface.Four was just called")
	}
	callInfo := struct {
		In1 c.Int
	}{
		In1: in1,
	}
	lockMockIfaceFour.Lock()
	mock.calls.Four = append(mock.calls.Four, callInfo)
	lockMockIfaceFour.Unlock()
	mock.FourFunc(in1)
}

// FourCalled returns true if at least one call was made to Four.
func (mock *MockIface) FourCalled() bool {
	lockMockIfaceFour.RLock()
	defer lockMockIfaceFour.RUnlock()
	return len(mock.calls.Four) > 0
}

// FourCalls gets all the calls that were made to Four.
// Check the length with:
//     len(mockedIface.FourCalls())
func (mock *MockIface) FourCalls() []struct {
	In1 c.Int
} {
	var calls []struct {
		In1 c.Int
	}
	lockMockIfaceFour.RLock()
	calls = mock.calls.Four
	lockMockIfaceFour.RUnlock()
	return calls
}

// One calls OneFunc.
func (mock *MockIface) One(str string, variadic ...string) (string, []string) {
	if mock.OneFunc == nil {
		panic("moq: MockIface.OneFunc is nil but Iface.One was just called")
	}
	callInfo := struct {
		Str      string
		Variadic []string
	}{
		Str:      str,
		Variadic: variadic,
	}
	lockMockIfaceOne.Lock()
	mock.calls.One = append(mock.calls.One, callInfo)
	lockMockIfaceOne.Unlock()
	return mock.OneFunc(str, variadic...)
}

// OneCalled returns true if at least one call was made to One.
func (mock *MockIface) OneCalled() bool {
	lockMockIfaceOne.RLock()
	defer lockMockIfaceOne.RUnlock()
	return len(mock.calls.One) > 0
}

// OneCalls gets all the calls that were made to One.
// Check the length with:
//     len(mockedIface.OneCalls())
func (mock *MockIface) OneCalls() []struct {
	Str      string
	Variadic []string
} {
	var calls []struct {
		Str      string
		Variadic []string
	}
	lockMockIfaceOne.RLock()
	calls = mock.calls.One
	lockMockIfaceOne.RUnlock()
	return calls
}

// Three calls ThreeFunc.
func (mock *MockIface) Three(in1 av1.Int) bv1.Str {
	if mock.ThreeFunc == nil {
		panic("moq: MockIface.ThreeFunc is nil but Iface.Three was just called")
	}
	callInfo := struct {
		In1 av1.Int
	}{
		In1: in1,
	}
	lockMockIfaceThree.Lock()
	mock.calls.Three = append(mock.calls.Three, callInfo)
	lockMockIfaceThree.Unlock()
	return mock.ThreeFunc(in1)
}

// ThreeCalled returns true if at least one call was made to Three.
func (mock *MockIface) ThreeCalled() bool {
	lockMockIfaceThree.RLock()
	defer lockMockIfaceThree.RUnlock()
	return len(mock.calls.Three) > 0
}

// ThreeCalls gets all the calls that were made to Three.
// Check the length with:
//     len(mockedIface.ThreeCalls())
func (mock *MockIface) ThreeCalls() []struct {
	In1 av1.Int
} {
	var calls []struct {
		In1 av1.Int
	}
	lockMockIfaceThree.RLock()
	calls = mock.calls.Three
	lockMockIfaceThree.RUnlock()
	return calls
}

// Two calls TwoFunc.
func (mock *MockIface) Two(in1 int, in2 int) int {
	if mock.TwoFunc == nil {
		panic("moq: MockIface.TwoFunc is nil but Iface.Two was just called")
	}
	callInfo := struct {
		In1 int
		In2 int
	}{
		In1: in1,
		In2: in2,
	}
	lockMockIfaceTwo.Lock()
	mock.calls.Two = append(mock.calls.Two, callInfo)
	lockMockIfaceTwo.Unlock()
	return mock.TwoFunc(in1, in2)
}

// TwoCalled returns true if at least one call was made to Two.
func (mock *MockIface) TwoCalled() bool {
	lockMockIfaceTwo.RLock()
	defer lockMockIfaceTwo.RUnlock()
	return len(mock.calls.Two) > 0
}

// TwoCalls gets all the calls that were made to Two.
// Check the length with:
//     len(mockedIface.TwoCalls())
func (mock *MockIface) TwoCalls() []struct {
	In1 int
	In2 int
} {
	var calls []struct {
		In1 int
		In2 int
	}
	lockMockIfaceTwo.RLock()
	calls = mock.calls.Two
	lockMockIfaceTwo.RUnlock()
	return calls
}
