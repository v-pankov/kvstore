package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/urfave/cli"

	apiGrpc "github.com/vdrpkv/kvstore/generated/api/grpc"
)

func main() {
	DoMain()
}

func DoMain() {
	doMain(grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func doMain(opts ...grpc.DialOption) {
	ctx := context.Background()
	const target = "localhost:8080"

	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Fatalf("failed to dial grpc service: %v", err)
	}
	defer conn.Close()

	client := apiGrpc.NewKVStoreClient(conn)

	app := cli.App{
		Name: "client",

		Commands: []cli.Command{
			{
				Name:      "set",
				ArgsUsage: "key val",
				Action: func(c *cli.Context) error {
					args := c.Args()
					key := args[0]
					val := args[1]

					rsp, err := client.Set(
						ctx, &apiGrpc.SetRequest{
							Key: key,
							Val: []byte(val),
						},
					)
					fmt.Println(rsp, err)
					return err
				},
			},
			{
				Name:      "get",
				ArgsUsage: "key",
				Action: func(c *cli.Context) error {
					args := c.Args()
					key := args[0]
					rsp, err := client.Get(
						ctx, &apiGrpc.GetRequest{
							Key: key,
						},
					)
					fmt.Println(rsp, err)
					return err
				},
			},
			{
				Name:      "delete",
				ArgsUsage: "key",
				Action: func(c *cli.Context) error {
					args := c.Args()
					key := args[0]
					rsp, err := client.Delete(
						ctx, &apiGrpc.DeleteRequest{
							Key: key,
						},
					)
					fmt.Println(rsp, err)
					return err
				},
			},
		},
	}
	app.Run(os.Args)
}
