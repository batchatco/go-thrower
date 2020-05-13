# go-thrower
Package go-thrower implements a simple throw/catch exception wrapper around
panic. It catches its own panics, but lets the others through.
## Functions
### func DisableCatching()
*DisableCatching* will prevent thrown errors from being caught, and so they
will become regular panics. Do not use this in production code; it is for
debugging only. It is against Go style to let panics cross API boundaries. All thrown errors should be caught by *RecoverError* normally.
### func RecoverError(err \*error)
*RecoverError* catches a thrown error. The pointer passed in can be `nil` if you
don't care what the thrown error was.

Use it as follows:

        func doSomething() (err error) {  
           // This will catch thrown errors and set the return value to the thrown error.
           defer thrower.RecoverError(&err)
           // Do some things that might call thrower.Throw() eventually.
           // For example:
           r := somethingThatCanReturnError()
           thrower.ThrowIfError(r)  // 'r' becomes the function's return value
        }
### func ReEnableCatching()
*ReEnableCatching* re-enables catching of panics if they were disabled by *DisableCatching*.
### func Throw(err error)
*Throw* throws the given error, which should be caught by *RecoverError* normally.
### func ThrowIfError(err error)
*ThrowIfError* throws an error only if `err` is not `nil`.

