//The die package works wraps around error, panic and recover
//OnErr panics if the error passed to it isn't nil,
//Log takes recovers writes

package die

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var logger io.Writer = os.Stdout
var mu sync.Mutex

func SetLogger(w io.Writer) {
	mu.Lock()
	logger = w
	mu.Unlock()
}

func OnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Log() interface{} {
	val := recover()
	if val != nil {
		buf := []byte(fmt.Sprintln(val))
		mu.Lock()
		logger.Write(buf)
		mu.Unlock()
		return val
	}
	return nil
}

func RecoverTo(err *error) {
	*err = recover().(error)
}
