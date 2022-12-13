package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"buf.build/gen/go/bufbuild/eliza/bufbuild/connect-go/buf/connect/demo/eliza/v1/elizav1connect"
	elizav1 "buf.build/gen/go/bufbuild/eliza/protocolbuffers/go/buf/connect/demo/eliza/v1"
	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type GreetServer struct {
	elizav1connect.UnimplementedElizaServiceHandler
}

func (s *GreetServer) Say(ctx context.Context, req *connect.Request[elizav1.SayRequest]) (*connect.Response[elizav1.SayResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&elizav1.SayResponse{
		Sentence: fmt.Sprintf("Hello, %s!", req.Msg.Sentence),
	})
	res.Header().Set("Eliza-Version", "v1")
	return res, nil
}

func main() {
	greeter := &GreetServer{}
	mux := http.NewServeMux()
	path, handler := elizav1connect.NewElizaServiceHandler(greeter)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
