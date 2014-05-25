golang-die
==========

The die package wraps panic, recover and defer to help make go a little more scripty
for those times when scriptiness is what you need. It's designed to fit my needs
and hopefully, yours. The code is distrubuted under the MIT license, in hopes that it will
help others when they need to bang some code together

This package also allows for a chaos monkey, so that you can have higher error rates
in order to flesh out error handling code

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
