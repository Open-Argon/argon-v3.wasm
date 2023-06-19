package main

import (
	"strings"
	"syscall/js"
)

var runID = 0

func argonToJsValid(argon any) any {
	switch x := argon.(type) {
	case number:
		f, _ := x.Float64()
		return f
	case ArObject:
		if typeof(x) == "array" {
			arr := js.Global().Get("Array").New()
			for i, v := range x.obj["__value__"].([]any) {
				arr.SetIndex(i, argonToJsValid(v))
			}
			return arr
		} else if typeof(x) == "string" {
			return x.obj["__value__"].(string)
		}

		obj := js.Global().Get("Object").New()
		for k, v := range x.obj {
			obj.Set(anyToArgon(k, false, false, 3, 0, false, 0), argonToJsValid(v))
		}
		return obj
	case bool, string:
		return x
	default:
		return nil
	}
}

func wasmRun(code string) (any, ArErr) {
	runID++
	initRandom()
	global := makeGlobal()
	lines := strings.Split(code, "\n")
	codelines := []UNPARSEcode{}
	for i := 0; i < len(lines); i++ {
		codelines = append(codelines, UNPARSEcode{
			lines[i],
			lines[i],
			i + 1,
			"<wasm>",
		})
	}

	translated, translationerr := translate(codelines)
	if translationerr.EXISTS {
		return nil, translationerr
	}
	local := newscope()
	localvars := Map(anymap{
		"program": Map(anymap{
			"args":   []any{},
			"origin": "",
			"import": builtinFunc{"import", func(args ...any) (any, ArErr) {
				return nil, ArErr{"Import Error", "Cannot Import in WASM", 0, "<wasm>", "", true}
			}},
			"cwd": "",
			"exc": "",
			"file": Map(anymap{
				"name": "<wasm>",
				"path": "",
			}),
			"main":  true,
			"scope": global,
		}),
	})
	return ThrowOnNonLoop(run(translated, stack{global, localvars, local}))
}

func await(awaitable js.Value) ([]js.Value, []js.Value) {
	then := make(chan []js.Value)
	defer close(then)
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		then <- args
		return nil
	})
	defer thenFunc.Release()

	catch := make(chan []js.Value)
	defer close(catch)
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		catch <- args
		return nil
	})
	defer catchFunc.Release()

	awaitable.Call("then", thenFunc).Call("catch", catchFunc)

	select {
	case result := <-then:
		return result, nil
	case err := <-catch:
		return nil, err
	}
}
