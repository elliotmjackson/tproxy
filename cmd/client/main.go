package main

import (
	"context"
	"log"
	"net/http"

	"buf.build/gen/go/bufbuild/eliza/bufbuild/connect-go/buf/connect/demo/eliza/v1/elizav1connect"
	elizav1 "buf.build/gen/go/bufbuild/eliza/protocolbuffers/go/buf/connect/demo/eliza/v1"
	"github.com/bufbuild/connect-go"
)

func main() {
	client := elizav1connect.NewElizaServiceClient(
		http.DefaultClient,
		"http://localhost:56006",
		// connect.WithGRPC(),
	)
	res, err := client.Say(
		context.Background(),
		connect.NewRequest(&elizav1.SayRequest{Sentence: "Jane"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Sentence)
}
