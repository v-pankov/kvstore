package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spf13/cobra"

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

	var defMsg = func(msg string) string {
		if msg == "" {
			return "OK"
		}
		return msg
	}

	var (
		rootCmd = &cobra.Command{
			Short: "client",
			Long:  "kvstore client",
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		}

		setCmd = &cobra.Command{
			Short: "set",
			Use:   "set key val",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				key := args[0]
				val := args[1]

				rsp, err := client.Set(
					ctx, &apiGrpc.SetRequest{
						Key: key,
						Val: []byte(val),
					},
				)

				fmt.Printf("STATUS: %d, MESSAGE: %s\n", rsp.Status, defMsg(rsp.Message))
				return err
			},
		}

		getCmd = &cobra.Command{
			Short: "get",
			Use:   "get key",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				key := args[0]
				rsp, err := client.Get(
					ctx, &apiGrpc.GetRequest{
						Key: key,
					},
				)

				fmt.Printf("STATUS: %d, MESSAGE: %s, VALUE: %s\n", rsp.Status, defMsg(rsp.Message), string(rsp.Val))
				return err
			},
		}

		deleteCmd = &cobra.Command{
			Short: "delete",
			Use:   "delete key",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				key := args[0]
				rsp, err := client.Delete(
					ctx, &apiGrpc.DeleteRequest{
						Key: key,
					},
				)

				fmt.Printf("STATUS: %d, MESSAGE: %s\n", rsp.Status, defMsg(rsp.Message))
				return err
			},
		}
	)

	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(deleteCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Println("DONE")
}
