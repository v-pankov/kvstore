package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiGrpc "github.com/vdrpkv/kvstore/generated/api/grpc"
)

func main() {
	DoMain()
}

func DoMain() {
	doMain(grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func doMain(opts ...grpc.DialOption) {
	const target = "localhost:8080"

	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Fatalf("failed to dial grpc service: %v", err)
	}
	defer conn.Close()

	client := apiGrpc.NewKVStoreClient(conn)

	key := strconv.FormatInt(int64(time.Now().Minute()), 10)
	//key := "   \t\r\n "
	val := []byte(strconv.FormatInt(int64(time.Now().Second()), 10))
	ctx := context.Background()

	getRespBeforeSet, err := client.Get(
		ctx, &apiGrpc.GetRequest{
			Key: key,
		},
	)
	log.Printf("get before set\t| [%+v] %v", getRespBeforeSet, err)

	setResp, err := client.Set(
		ctx, &apiGrpc.SetRequest{
			Key: key,
			Val: val,
		},
	)
	log.Printf("set\t\t\t| [%+v] %v", setResp, err)

	getRespAfterSet, err := client.Get(
		ctx, &apiGrpc.GetRequest{
			Key: key,
		},
	)
	log.Printf("get after set\t| [%+v] %v", getRespAfterSet, err)

	deleteResp, err := client.Delete(
		ctx, &apiGrpc.DeleteRequest{
			Key: key,
		},
	)
	log.Printf("delete\t\t| [%+v] %v", deleteResp, err)

	getRespAfterDelete, err := client.Get(
		ctx, &apiGrpc.GetRequest{
			Key: key,
		},
	)
	log.Printf("get after delete\t| [%v] %v", getRespAfterDelete, err)
}
