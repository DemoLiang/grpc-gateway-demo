//package main

//
//import (
//	"context"
//	"flag"
//	"net/http"
//
//	"github.com/golang/glog"
//	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
//
//	gw "github.com/DemoLiang/grpc-gateway-demo/gen/gw/hello" // Update
//)
//
//var (
//	// command-line options:
//	// gRPC server endpoint
//	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
//)
//
//func run() error {
//	ctx := context.Background()
//	ctx, cancel := context.WithCancel(ctx)
//	defer cancel()
//
//	// Register gRPC server endpoint
//	// Note: Make sure the gRPC server is running properly and accessible
//	mux := runtime.NewServeMux()
//	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
//	err := gw.RegisterHelloHTTPHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
//	if err != nil {
//		return err
//	}
//
//	// Start HTTP server (and proxy calls to gRPC server endpoint)
//	return http.ListenAndServe(":8081", mux)
//}
//
//func main() {
//	flag.Parse()
//	defer glog.Flush()
//
//	if err := run(); err != nil {
//		glog.Fatal(err)
//	}
//}

package main

import (
	"context"
	"fmt"
	hello "github.com/DemoLiang/grpc-gateway-demo/gen/go/hello" // Update
	"github.com/felixge/httpsnoop"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	// importing generated stubs
	gw "github.com/DemoLiang/grpc-gateway-demo/gen/gw/hello" // Update
)

// GreeterServerImpl will implement the service defined in protocol buffer definitions
type GreeterServerImpl struct {
	//gw.UnimplementedGreeterServer
	hello.UnsafeHelloHTTPServer
}

// SayHello is the implementation of RPC call defined in protocol definitions.
// This will take HelloRequest message and return HelloReply
func (g *GreeterServerImpl) SayHello(ctx context.Context, request *hello.HelloHTTPRequest) (*hello.HelloHTTPResponse, error) {
	//if err := request.Validate(); err != nil {
	//	return nil, err
	//}
	return &hello.HelloHTTPResponse{
		Message: fmt.Sprintf("hello %s", request.Name),
	}, nil
}
func main() {
	// create new gRPC server
	grpcSever := grpc.NewServer()
	// register the GreeterServerImpl on the gRPC server
	hello.RegisterHelloHTTPServer(grpcSever, &GreeterServerImpl{})
	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux()
	// setting up a dail up for gRPC service by specifying endpoint/target url
	err := gw.RegisterHelloHTTPHandlerFromEndpoint(context.Background(), mux, "localhost:8081", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatal(err)
	}
	// Creating a normal HTTP server
	server := http.Server{
		Handler: withLogger(mux),
	}
	// creating a listener for server
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	m := cmux.New(l)
	// a different listener for HTTP1
	httpL := m.Match(cmux.HTTP1Fast())
	// a different listener for HTTP2 since gRPC uses HTTP2
	grpcL := m.Match(cmux.HTTP2())
	// start server
	// passing dummy listener
	go server.Serve(httpL)
	// passing dummy listener
	go grpcSever.Serve(grpcL)
	// actual listener
	m.Serve()
}
func withLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}
