package main

import "scanner"

func main() {
	engine := scanner.GetEngine()
	err := engine.Run()
	if err != nil {
		panic(err)
	}
}
