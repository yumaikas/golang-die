package die

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

//This file will support per goroutine loggers that work without locks, once it's complete.
//A doom context allows us to reconstrcut where we where.
//It is intended to be used per goroutine, and typically per function call
//It is NOT multithread safe. It's not intended to be, but this is a judgement call that can be challeged with a strong case for it to be locked
type DoomContext struct {
	logger   io.Writer
	numCalls int
	funcName string
}

func Context(funcName string) *DoomContext {
	return &DoomContext{logger: os.Stdout, numCalls: 0, funcName: funcName}
}

//Save three lines of code when you need to panic.
func onErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Automatically increment a counter so that we know where we where
func (dc *DoomContext) OnErr(err error) {
	dc.numCalls++
	onErr(err)
}

func (dc *DoomContext) OnErrAt(err error, errNum int) {
	dc.numCalls = errNum
	onErr(err)
}

//This allows you to log to any io.Writer you want to.
func (dc *DoomContext) SetLogger(w io.Writer) {
	dc.logger = w
}

//Used by the logging functions.
func (dc *DoomContext) logForFunc(errVal interface{}) {
	if errVal != nil {
		buf := &bytes.Buffer{}
		fmt.Fprintln(buf,
			"Panic caught in func:", dc.funcName,
			fmt.Sprintf(" at failure point %v,", dc.numCalls),
			"err:", errVal)
		buf.WriteTo(logger)
	}
}

//Log if an error occurs inside the function named by the string.
func (dc *DoomContext) Log() {
	errVal := recover()
	dc.logForFunc(errVal)
}

//Like log, but sets an error value when a panic happens
func (dc DoomContext) LogErr(funcName string, err *error) {
	errVal := recover()
	if err != nil {
		switch errVal.(type) {
		case error:
			*err = errVal.(error)
		default:
			*err = fmt.Errorf("Error: %v", errVal)
		}
	}
	dc.logForFunc(errVal)
}

//Logs a an error, if one occurs, and then calls a function that handles
//any custom tracing or setting of return values that might be needed
//by the caller
func (dc DoomContext) LogSettingReturns(funcName string, setReturns func()) {
	errVal := recover()
	if errVal != nil {
		dc.logForFunc(errVal)
		setReturns()
	}
}
