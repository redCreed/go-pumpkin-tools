package main

import (
	"context"
	"go-do/etcd/client_balance"
	"go-do/etcd/client_balance/proto"
	"log"
	"net"
	"google.golang.org/grpc"
)

// SimpleService 定义我们的服务
type SimpleService struct{}

const (
	// Address 监听地址
	Address string = "localhost:8002"
	// Network 网络通信协议
	Network string = "tcp"
	// SerName 服务名称
	SerName string = "simple_grpc"
)

// EtcdEndpoints etcd地址
var EtcdEndpoints = []string{"localhost:2379"}

func main() {
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	log.Println(Address + " net.Listing...")
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务

	proto.RegisterSimpleServer(grpcServer, &SimpleService{})
	//把服务注册到etcd
	ser, err := client_balance.NewRegister(EtcdEndpoints, SerName, Address, 5)
	if err != nil {
		log.Fatalf("register service err: %v", err)
	}
	defer ser.Close()
	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}

// Route 实现Route方法
func (s *SimpleService) Route(ctx context.Context, req *proto.SimpleRequest) (*proto.SimpleResponse, error) {
	log.Println("receive: " + req.Data)
	res := proto.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return &res, nil
}
