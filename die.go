//The die package wraps panic, recover and defer to help make go a little more scripty
//for those times when scriptiness is what you need. It's designed to fit my needs
//and hopefully, yours. The code is distrubuted under the MIT license, in hopes that it will
//help others when they need to bang some code together
//
//This package also allows for a chaos monkey, so that you can have higher error rates
//in order to flesh out error handling code
//
//On important note is this. In go, recover only works inside a directly deferred function.
//
//This means that the below will work:
//
//    defer die.Log("main")
//
//But the following won't
//
//    defer func() {die.Log("main")}()
//
//This is a simplistic case, in real code it could be more complicated,
//but this is a more subtle detail that you need to know it before using this library
//
/*
With that out of the way, examples:
The simplest demonstration of the package [Playgound link]

	package main

	import (
        "fmt"
        "github.com/yumaikas/die"
	)

    func main() {
         defer die.Log("main")
         die.OnErr(fmt.Errorf("Error!"))
    }

Should output:

    Panic caught in func: main err: Error!

Something a little more complicated:

	package main

	import (
        "fmt"
        "math/rand"
        "github.com/yumaikas/die"
	)

    func main() {
    	err := updateInfo()
	}

	//Use a named return value to be able to set errors from LogErr
	func updateInfo() (err error) {
		defer die.LogErr("updateInfo", &err)
		die.OnErr(runDatabaseCall())
	}

	func runDatabaseCall() error {
		//Pretending to fail 50% of the time. Your database is definitely more reliable than this. :)
		if rand.NextInt(2) = 1 {
			return fmt.Errorf("Database connection timed out")
		} else {
			return nil
		}
	}

*/

package die

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
)

var logger io.Writer = os.Stdout
var mu sync.Mutex

//Save three lines of code when you need to panic.
func OnErr(err error) {
	if err != nil {
		panic(err)
	}
}

//This allows you to log to any io.Writer you want to.
func SetLogger(w io.Writer) {
	mu.Lock()
	logger = w
	mu.Unlock()
}

//Used by the logging functions.
func logForFunc(funcName string, errVal interface{}) {
	if errVal != nil {
		buf := &bytes.Buffer{}
		fmt.Fprintln(buf, "Panic caught in func:", funcName+",", "err:", errVal)
		mu.Lock()
		buf.WriteTo(logger)
		mu.Unlock()
	}
}

//Log if an error occurs inside the function named by the string.
func Log(funcName string) {
	errVal := recover()
	logForFunc(funcName, errVal)
}

func anyToErr(val interface{}) error {
	if val == nil {
		return nil
	}
	switch val.(type) {
	case error:
		return val.(error)
	default:
		return fmt.Errorf("Error: %v", val)
	}
}

//Like log, but sets an error value when a panic happens
func LogErr(funcName string, err *error) {
	errVal := recover()
	*err = anyToErr(errVal)
	logForFunc(funcName, errVal)
}

//Logs a an error, if one occurs, and then calls a function that handles
//any custom tracing or setting of return values that might be needed
//by the caller
func LogSettingReturns(funcName string, err *error, setReturns func()) {
	errVal := recover()
	logForFunc(funcName, errVal)
	*err = anyToErr(errVal)
	if errVal != nil {
		setReturns()
	}
}

//For the 5% that the above functions do not satisfy, below is the power to take any action
//with the value from recover. It is verbosely named to discourage use, but here if needed.
//It also locks the err logger
func CustomLoggingWithDiePackage(CustomFunc func(io.Writer, interface{})) {
	errVal := recover()
	mu.Lock()
	CustomFunc(logger, errVal)
	mu.Unlock()
}

//If you need a chaos monkey in your code to flesh out error handling paths, this is for you
func Chaos(err error) error {
	if rand.Intn(1) == 1 {
		return err
	} else {
		return errors.New("The chaos monkey found you!")
	}
}
