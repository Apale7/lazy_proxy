package proxy_pool

import (
	"fmt"
	"testing"
)

func TestDefaultProxyPool_CheckProxy(t *testing.T) {
	p := &DefaultProxyPool{}
	fmt.Printf("p.CheckProxy(\"178.134.208.126:50824\"): %v\n", p.CheckProxy("178.134.208.126:50824"))
}
