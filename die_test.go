package die

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestNilLog(t *testing.T) {
	defer Log("TestLog")
	log := &bytes.Buffer{}
	SetLogger(log)
	OnErr(nil)
	if log.Len() > 0 {
		t.Fail()
	}
}

type testBuffer interface {
	fmt.Stringer
	io.Writer
}

func TestLog(t *testing.T) {
	log := tLog(t)
	if log.String() != "Panic caught in func: tLog, err: Testing Error!\n" {
		t.Fail()
		t.Log(`"` + log.String() + `"`)
	}
}

func tLog(t *testing.T) (lg *bytes.Buffer) {
	log := &bytes.Buffer{}
	SetLogger(log)
	defer func() { lg = log }()
	defer Log("tLog")
	OnErr(errors.New("Testing Error!"))
	return
}

func TestErrLog(t *testing.T) {
	log, err := tErrLog()
	if log.Len() <= 0 || err != loggedErr {
		t.Fail()
	}
}

var loggedErr = errors.New("Error to log!")

func tErrLog() (lg *bytes.Buffer, err error) {
	log := &bytes.Buffer{}
	SetLogger(log)
	defer func() { lg = log }()
	defer LogErr("tErrLog", &err)
	OnErr(loggedErr)
	return
}

func TestErrLogSetReturns(t *testing.T) {
	t.Fail()
	t.Log(tErrReturnValLog())
}

func tErrReturnValLog() (val int, lg *bytes.Buffer, err error) {
	log := &bytes.Buffer{}
	val = 0
	SetLogger(log)
	defer LogSettingReturns("tErrReturnValLog", &err, func() { lg = log; val = 1 })
	OnErr(loggedErr)
	return
}
