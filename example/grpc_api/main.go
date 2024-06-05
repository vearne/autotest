package main

import (
	"context"
	pb "github.com/vearne/autotest/example/grpc_api/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

var books = []pb.Book{
	{Id: 1, Title: "The Go Programming Language", Author: "Alan A. A. Donovan and Brian W. Kernighan"},
	{Id: 2, Title: "Effective Go", Author: "The Go Authors"},
}

type BookServer struct {
	lock    sync.RWMutex
	counter int64
}

func NewBookServer() *BookServer {
	var b BookServer
	b.counter = 3
	return &b
}

func (b BookServer) AddBook(ctx context.Context, in *pb.AddBookRequest) (*pb.AddBookReply, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	var book pb.Book
	book.Id = b.counter
	b.counter++
	book.Title = in.Title
	book.Author = in.Author
	books = append(books, book)

	out := new(pb.AddBookReply)
	out.Code = pb.CodeEnum_Success
	out.Data = &book

	return out, nil
}

func (b BookServer) DeleteBook(ctx context.Context, in *pb.DeleteBookRequest) (*pb.DeleteBookReply, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	out := new(pb.DeleteBookReply)
	for i, book := range books {
		if book.Id == in.Id {
			books = append(books[:i], books[i+1:]...)
			out.Code = pb.CodeEnum_Success
			return out, nil
		}
	}
	out.Code = pb.CodeEnum_NoDataFound
	return out, nil
}

func (b BookServer) UpdateBook(ctx context.Context, in *pb.UpdateBookRequest) (*pb.UpdateBookReply, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	out := new(pb.UpdateBookReply)
	for i, book := range books {
		if book.Id == in.Id {
			if len(in.Title) > 0 {
				books[i].Title = in.Title
			}
			if len(in.Author) > 0 {
				books[i].Author = in.Author
			}
			out.Code = pb.CodeEnum_Success
			data := books[i]
			out.Data = &data
			return out, nil
		}
	}
	out.Code = pb.CodeEnum_NoDataFound
	return out, nil
}

func (b BookServer) GetBook(ctx context.Context, in *pb.GetBookRequest) (*pb.GetBookReply, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	out := new(pb.GetBookReply)
	for i, book := range books {
		if book.Id == in.Id {
			out.Code = pb.CodeEnum_Success
			data := books[i]
			out.Data = &data
			return out, nil
		}
	}

	out.Code = pb.CodeEnum_NoDataFound
	return out, nil
}

func (b BookServer) ListBook(ctx context.Context, in *pb.ListBookRequest) (*pb.ListBookReply, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	out := new(pb.ListBookReply)
	out.Code = pb.CodeEnum_Success
	out.List = make([]*pb.Book, 0)
	for i := 0; i < len(books); i++ {
		record := books[i]
		out.List = append(out.List, &record)
	}
	return out, nil
}

func main() {
	server := grpc.NewServer()
	pb.RegisterBookstoreServer(server, NewBookServer())
	// Register reflection service on gRPC server.
	reflection.Register(server)
	lis, err := net.Listen("tcp", ":50031")
	if err != nil {
		log.Fatalf("failed to listen, %v\n", err)
	}
	log.Println("starting...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("server.Serve, %v\n", err)
	}
}
