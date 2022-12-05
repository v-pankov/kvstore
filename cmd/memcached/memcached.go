package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached"
)

func main() {
	conn, err := memcached.OpenTCP(net.IPv4(127, 0, 0, 1), 11211)
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Set("x", 0, 0, time.Now().String()); err != nil {
		log.Fatal(err)
	}

	items, err := conn.Get("x")
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		log.Println(item)
	}

	if err := conn.Delete("x"); err != nil {
		log.Fatal(err)
	}

	fmt.Println(conn.Delete("x"))

	fmt.Println(conn.Get("x"))
}
