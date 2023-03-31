package main

import (
	"context"
	"grpc-demo/product"
	"log"
	"net"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const port = ":50051"

type server struct {
	product.ProductInfoServer
	productMap map[string]*product.Product
}

// 添加商品
func (s *server) AddProduct(ctx context.Context, req *product.Product) (resp *product.ProductId, err error) {
	resp = &product.ProductId{}
	out, err := uuid.NewV4()
	if err != nil {
		return resp, status.Errorf(codes.Internal, "err while generate the uuid ", err)
	}

	req.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*product.Product)
	}

	s.productMap[req.Id] = req
	resp.Value = req.Id
	return
}

// 获取商品
func (s *server) GetProduct(ctx context.Context, req *product.ProductId) (resp *product.ResponseSingleDTO, err error) {
	if s.productMap == nil {
		s.productMap = make(map[string]*product.Product)
	}

	resp = &product.ResponseSingleDTO{
		Code: 200,
		Msg:  "ok",
		Data: s.productMap[req.Value],
	}

	return
}

// list商品
func (s *server) ListProduct(ctx context.Context, req *product.QueryRequest) (resp *product.ResponseDTO, err error) {
	products := make([]*product.Product, 0, len(s.productMap))

	for _, v := range s.productMap {
		products = append(products, v)
	}

	resp = &product.ResponseDTO{
		Code: 200,
		Msg:  "ok",
		Data: products[req.PageNo:req.PageSize],
	}
	return
}

// 删除商品
func (s *server) DelProduct(ctx context.Context, req *product.ProductId) (resp *product.ResponseBool, err error) {
	log.Printf("Deleting the key named %s from the productMap", req.Value)
	if _, ok := s.productMap[req.Value]; ok {
		delete(s.productMap, req.Value)
		resp = &product.ResponseBool{
			Code: 200,
			Msg:  "ok",
			Data: true,
		}
	}
	resp = &product.ResponseBool{
		Code: 400,
		Msg:  "Data not found",
		Data: false,
	}
	return
}

// 更新商品
func (s *server) UpdateProduct(ctx context.Context, req *product.Product) (resp *product.ResponseSingleDTO, err error) {
	if _, ok := s.productMap[req.Id]; ok {
		s.productMap[req.Id] = req
		resp = &product.ResponseSingleDTO{
			Code: 200,
			Msg:  "ok",
			Data: s.productMap[req.Id],
		}
	}
	resp = &product.ResponseSingleDTO{
		Code: 400,
		Msg:  "data not found",
		Data: req,
	}
	return
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	s := grpc.NewServer()
	product.RegisterProductInfoServer(s, &server{})
	log.Println("start gRPC listen on port " + port)
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}
