// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.uber.org/yarpc"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/examples/thrift/keyvalue/kv/yarpc/keyvalueclient"
	"go.uber.org/yarpc/transport/http"
	"go.uber.org/yarpc/transport/tchannel"
)

func main() {
	outboundName := ""
	flag.StringVar(
		&outboundName,
		"outbound", "", "name of the outbound to use (http/tchannel)",
	)

	flag.Parse()

	httpTransport := http.NewTransport()
	tchannelTransport := tchannel.NewChannelTransport(tchannel.ServiceName("keyvalue-client"))

	var outbound transport.UnaryOutbound
	switch strings.ToLower(outboundName) {
	case "http":
		outbound = httpTransport.NewSingleOutbound("http://127.0.0.1:24034")
	case "tchannel":
		outbound = tchannelTransport.NewSingleOutbound("localhost:28941")
	default:
		log.Fatalf("invalid outbound: %q\n", outboundName)
	}

	cache := NewCacheOutboundMiddleware()
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: "keyvalue-client",
		Outbounds: yarpc.Outbounds{
			"keyvalue": {
				Unary: outbound,
			},
		},
		OutboundMiddleware: yarpc.OutboundMiddleware{
			Unary: cache,
		},
	})
	if err := dispatcher.Start(); err != nil {
		log.Fatalf("failed to start Dispatcher: %v", err)
	}
	defer dispatcher.Stop()

	client := keyvalueclient.New(dispatcher.ClientConfig("keyvalue"))

	scanner := bufio.NewScanner(os.Stdin)
	rootCtx := context.Background()
	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Split(line, " ")
		if len(args) < 1 || len(args[0]) < 3 {
			continue
		}

		cmd := args[0]
		args = args[1:]

		switch cmd {

		case "get":
			if len(args) != 1 {
				fmt.Println("usage: get key")
				continue
			}
			key := args[0]

			ctx, cancel := context.WithTimeout(rootCtx, 100*time.Millisecond)
			defer cancel()

			if value, _, err := client.GetValue(ctx, nil, &key); err != nil {
				fmt.Printf("get %q failed: %s\n", key, err)
			} else {
				fmt.Println(key, "=", value)
			}
			continue

		case "set":
			if len(args) != 2 {
				fmt.Println("usage: set key value")
				continue
			}
			key, value := args[0], args[1]

			cache.Invalidate()
			ctx, cancel := context.WithTimeout(rootCtx, 100*time.Millisecond)
			defer cancel()

			if _, err := client.SetValue(ctx, nil, &key, &value); err != nil {
				fmt.Printf("set %q = %q failed: %v\n", key, value, err.Error())
			}
			continue

		case "exit":
			return

		default:
			fmt.Println("invalid command", cmd)
			fmt.Println("valid commansd are: get, set, exit")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error:", err.Error())
	}
}
