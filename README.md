golang-die
==========

The die package wraps panic, recover and defer to help make go a little more scripty/hacky
for those times when scriptiness is what you need. It's designed to fit my needs
and hopefully, yours. The code is distrubuted under the MIT license, in hopes that it will
help others when they need to bang some code together.

This package also provides a chaos monkey, so that you can have higher error rates
in order to flesh out error handling code. 

On important note is this. In go, recover only works inside a directly deferred function.

This means that the below will work:

    defer die.Log("main")

But the following won't

    defer func() {die.Log("main")}()

This is a simplistic case, in real code it could be more complicated,
but this is a subtle detail that you need to know it before using this library

Simplest example below:


```go
	package main

	import (
        "fmt"
        "github.com/yumaikas/die"
	)

    func main() {
         defer die.Log("main")
         
         die.OnErr(fmt.Errorf("Error!"))
    }
```
Which should output: "Panic caught in func: main err: Error!"

In the more complicated case, it looks like this:


Something a little more complicated:

```go

    package main

    import (
        "fmt"
        "math/rand"
        "github.com/yumaikas/golang-die"
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
```

And finally, you can pass in a `func` that performs failure cleanup (rolling back transactions, zeroing return values, etc) using ``.

```go
    package main

    import (
        "fmt"
        "math/rand"
        "github.com/yumaikas/golang-die"
    )

    func main() {
        err := updateInfo()
    }

    //Use a named return value to be able to set errors from LogErr
    func updateInfo() (err error) {
        defer die.LogErr("updateInfo", &err)
        numRows, err := countRowsInTable()
        die.OnErr(err)
        fmt.Println("Counted", numRows, "rows in database")
    }

    func countRowsInTable() (numRows int, err error) {
        defer die.LogSettingReturns("countRowsInTable", &err, func(){ numRows = 0 })
        if rand.NextInt(2) = 1 {
            die.OnErr(fmt.Errorf("Database connection timed out"))
        } else {
            return 10, nil
        }
    }
```