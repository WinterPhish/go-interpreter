//go:build js && wasm

package main

import (
	"syscall/js"

	"selenite/session"
)

func main() {
	s := session.New()

	js.Global().Set("seleniteRun", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return ""
		}
		return s.Run(args[0].String())
	}))

	select {} // keep the program alive so the page can keep calling it
}