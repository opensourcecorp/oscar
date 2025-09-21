// Package main
package main

import "fmt"

//go:generate go build -o ../build/protoc-gen-oscarcfg .
//go:generate buf generate

func main() {
	fmt.Println("not implemented, but this generator should generate a populated oscar.yaml")
}
