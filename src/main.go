package main

import (
	"fmt"
	"os"
	"syscall/js"
)

// args without the program path
var Args = os.Args[1:]

type stack = []ArObject

const VERSION = "3.0.0"

func newscope() ArObject {
	return Map(anymap{})
}

func main() {
	debugInit()
	c := make(chan bool)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("There was a fundamental error in argon v3 that caused it to crash.")
			fmt.Println()
			if fork {
				fmt.Println("This is a fork of Open-Argon. Please report this to the fork's maintainer.")
				fmt.Println("Fork repo:", forkrepo)
				fmt.Println("Fork issue page:", forkissuesPage)
				fmt.Println()
				fmt.Println("website:", website)
				fmt.Println("docs:", docs)
				fmt.Println()
				if fork {
					fmt.Println("This is a fork of Open-Argon. Please report this to the fork's maintainer.")
					fmt.Println("Fork repo:", forkrepo)
					fmt.Println("Fork issue page:", forkissuesPage)
					fmt.Println()
				} else {
					fmt.Println("Please report this to the Open-Argon team.")
					fmt.Println("Main repo:", mainrepo)
					fmt.Println("Issue page:", mainissuesPage)
					fmt.Println()
				}
				fmt.Println("please include the following information:")
				fmt.Println("panic:", r)
				os.Exit(1)
			}
		}
	}()
	initRandom()
	garbageCollect()
	obj := js.Global().Get("Object").New()
	obj.Set("eval", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		code := ""
		if len(args) >= 1 {
			code = args[0].String()
		}
		val, err := wasmRun(code)
		if err.EXISTS {
			panicErr(err)
			return js.Null()
		}

		return js.ValueOf(argonToJsValid(val))
	}))
	js.Global().Set("Ar", obj)
	<-c
}
