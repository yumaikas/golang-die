package die

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
)

//Save three lines of code when you need to panic.
func OnErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("Error in function: %s; Details: %s", getFunctionName(), err))
	}
}

func getFunctionName() string {
	// Get the caller of OnErr's name.
	if pc, _, _, ok := runtime.Caller(2); ok {
		return runtime.FuncForPC(pc).Name()
	}
	return "function not found"
}

//If you need a chaos monkey in your code to flesh out error handling paths, this is for you
func Chaos(err error) error {
	if rand.Intn(1) == 1 {
		return err
	} else {
		return errors.New("The chaos monkey found you!")
	}
}
