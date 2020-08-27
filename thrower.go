/*
Package thrower implements a simple throw/catch exception wrapper around panic.
It only catches the panics that are thrown. Other panics are passed through.
*/
package thrower

var (
	disabled bool // disables catching for debugging only
)

// We only catch our own thrown errors.  Any other panics are passed through.
type thrown struct {
	error
}

// CatchState is passed to SetCatching
type CatchState bool

const (
	Catch CatchState = iota == 1
	DontCatch
)

func newThrown(err error) thrown {
	return thrown{err}
}

func (th thrown) toError() error {
	return th.error
}

func (th thrown) Error() string {
	return th.error.Error()
}

// Throw throws the given error, which can be caught by RecoverError potentially.
func Throw(err error) {
	if disabled {
		panic(err)
	}
	panic(newThrown(err))
}

// ThrowIfError throws an error only if err is not nil.
func ThrowIfError(err error) {
	if err != nil {
		Throw(err)
	}
}

// RecoverError catches a thrown error. The pointer
// passed in can be nil if you don't care what the
// thrown error was.
//
// Use it as follows:
//
//  func doSomething() (err error) {
//     defer thrower.RecoverError(&err)
//     // do some things that might call thrower.Throw() eventually.
//     r := somethingThatCanReturnError()
//     thrower.ThrowIfError(r)
//  }
func RecoverError(err *error) {
	if disabled {
		return
	}
	// Attempt to convert the panic to an error
	if r := recover(); r != nil {
		th, has := r.(thrown)
		if has {
			// This is our panic
			if err == nil {
				return
			}
			*err = th.toError()
			return
		}
		// This is someone else's panic.
		panic(r)
	}
}

// SetCatching is for debugging only. It sets thrown error catching on or off and
// returns the  previous state. The default is Catch and does not need to be set explicitly
// Disabled thrown errors become just regular panics and are not caught.
// Do not use this in production code. It is for debugging only.
func SetCatching(state CatchState) CatchState {
	old := disabled
	disabled = bool(state)
	return CatchState(old)
}
