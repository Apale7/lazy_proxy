package main

import (
	"fmt"
	"testing"
)

func TestDefaultProxyGetter_CheckProxy(t *testing.T) {
	p := &DefaultProxyGetter{}
	fmt.Printf("p.CheckProxy(\"178.134.208.126:50824\"): %v\n", p.CheckProxy("178.134.208.126:50824"))
}
