//The die package works wraps around error, panic and recover
//OnErr panics if the error passed to it isn't nil,

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

//Call with recover() inside
func Log(val interface{}) interface{} {
	if val != nil {
		buf := []byte(fmt.Sprintln(val))
		mu.Lock()
		logger.Write(buf)
		mu.Unlock()
		return val
	}
	return nil
}
