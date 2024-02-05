package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"github.com/msecommerce/api_gateway/pkg/graaphql"
	"github.com/msecommerce/api_gateway/pkg/middleware"
	"github.com/sidsreeram/msproto/pb"
	"google.golang.org/grpc"
)

func main() {
	userconn, err := grpc.Dial("localhost:40004", grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
	}


	productConn, err := grpc.Dial("localhost:40002", grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
	}

	cartConn, err := grpc.Dial("localhost:40003", grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
	}

	

	defer func() {
		userconn.Close()
		productConn.Close()
		cartConn.Close()
	}()

	userRes := pb.NewUserServiceClient(userconn)
	productRes := pb.NewProductServiceClient(productConn)
	cartRes := pb.NewCartServiceClient(cartConn)

	if err := godotenv.Load("app.env"); err != nil {
		log.Fatal(err.Error())
	}
	secretstr := os.Getenv("SECRET")

	graaphql.Initialize(userRes, productRes, cartRes)
	graaphql.RetrieveSecret(secretstr)
	middleware.InitMiddleware(secretstr)

	h := handler.New(&handler.Config{
		Schema: &graaphql.Schema,
		Pretty: true,
	})
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		// Add the http.ResponseWriter to the context.
		ctx := context.WithValue(r.Context(), "httpResponseWriter", w)
		ctx = context.WithValue(ctx, "request", r)

		// Update the request's context.
		r = r.WithContext(ctx)

		// Call the GraphQL handler.
		h.ContextHandler(ctx, w, r)
	})

	log.Println("listening on port :8080 of api gateway")

	http.ListenAndServe(":8080", nil)

}
