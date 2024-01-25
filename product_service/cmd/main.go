package main

import (
	"fmt"
	"log"
	"net"
	"os"
	pb"github.com/msecommerce_product_service/grpc/pb"

	"github.com/joho/godotenv"
	"github.com/msecommerce_product_service/pkg/db"
	"github.com/msecommerce_product_service/pkg/initializer"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err.Error())
	}

	addr := os.Getenv("DB_URL")

	DB, err := db.ConnectDatabase(addr)
	if err != nil {
		log.Fatalf("unable to connect with database : %v", err)
	}
	
	usecase := initializer.Initialize(DB)
	
	server := grpc.NewServer()
	
	pb.RegisterProductServiceServer(server,usecase)
	
	listener, err := net.Listen("tcp", ":40001")
	
	if err != nil {
		log.Println(fmt.Errorf("Failed to listen on port 40001: %w", err))
	}
	
	log.Printf("Product Server is listening on port")
	
	if err = server.Serve(listener); err != nil {
		log.Println(fmt.Errorf("Failed to listen on port 40001: %w", err))
	}
}	