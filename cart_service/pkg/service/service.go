package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/msecommerce/cart_service/pkg/adapter/interfaces"
	"github.com/msecommerce/cart_service/pkg/client"
	"github.com/msecommerce/cart_service/pkg/models"
	"github.com/sidsreeram/msproto/pb"
)

var (
	ProductClient pb.ProductServiceClient
)

type CartServiceServer struct {
	Adapter interfaces.CartAdapter
	Product *client.ProductClient
    pb.UnimplementedCartServiceServer
}

func NewCartServiceServer(adapter interfaces.CartAdapter,product *client.ProductClient) *CartServiceServer {
	return &CartServiceServer{Adapter: adapter,Product: product}
}
func (c *CartServiceServer) CreateCart(ctx context.Context, req *pb.CartRequest) (*pb.CartResponse, error) {
	res, err := c.Adapter.CreateCart(models.Cart{UserID: req.UserId})
	if err != nil {
		return nil, fmt.Errorf("Error occured in creating user's cart:%w", err)
	}
	return &pb.CartResponse{CartId: uint64(res.ID), UserId: res.UserID, IsEmpty: true}, nil
}
func (c *CartServiceServer) Get(req *pb.CartRequest, stream pb.CartService_GetServer) error {
	var ids []*pb.ProductIdRequest
	res, err := c.Adapter.GetCart(req.UserId)
	if err != nil {
		return fmt.Errorf("error happened in getting cart @service: %w", err)
	}
	if len(res) == 0 {
		return nil
	}
	for _, pro := range res {
		ids = append(ids, &pb.ProductIdRequest{Id: pro.ProductID})
	}
	productStream, err := ProductClient.GetMultiple(context.Background(), &pb.ProductMultipleRequest{Ids: ids})
	if err != nil {
		return fmt.Errorf("error happened in getting multiple products: %w", err)
	}

	for i, product := range productStream.Multipleresponse {
		addtocart := &pb.AddToCartResponse{
			Product:  product,
			Quantity: res[i].Quantity,
		}
		if err := stream.Send(addtocart); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (c *CartServiceServer) AddtoCart(ctx context.Context, req *pb.AddTOCartRequest) (*pb.AddToCartResponse, error) {
	productRes, err := ProductClient.Get(context.TODO(), &pb.ProductIdRequest{Id: req.ProductId})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if productRes.Name == "" {
		return nil, errors.New("Sorry, the product doesn't exist")
	}

	item, err := c.Adapter.AddToCart(models.CartItems{
		ProductID: req.ProductId,
		Quantity:  req.Quantity,
	}, req.UserId)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &pb.AddToCartResponse{
		Product: &pb.ProductResponse{
			Id:       productRes.Id,
			Name:     productRes.Name,
			Price:    productRes.Price,
			Quantity: item.Quantity,
		},
	}, nil
}

func (c *CartServiceServer) Delete(ctx context.Context, req *pb.AddTOCartRequest) (*pb.AddToCartResponse, error) {
	productResult, err := ProductClient.Get(context.TODO(), &pb.ProductIdRequest{Id: req.ProductId})
	if err != nil {
		return nil, fmt.Errorf("Error in geting productId :%w", err)
	}
	if productResult.Name == "" {
		return nil, fmt.Errorf("A product doesn't find :%w", err)
	}
	err = c.Adapter.DeleteCartItem(models.CartItems{ProductID: productResult.Id}, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("Error in Deleting product :%w", err)
	}
	return &pb.AddToCartResponse{}, nil
}
func (c *CartServiceServer) UpdateQuantity(ctx context.Context, req *pb.UpdateQuantityRequest) (*pb.AddToCartResponse, error) {
	var res models.CartItems
	var err error
	if req.IsIncreased {
		res, err = c.Adapter.IncrementQuantity(models.CartItems{
			ProductID: req.ProductId,
			Quantity:  req.Quantity,
		}, req.UserId)
		if err != nil {
			return nil, fmt.Errorf("Error in incrementing quantity:%w", err)
		}
	} else {
		res, err = c.Adapter.DecrementQuantity(models.CartItems{
			ProductID: req.ProductId,
			Quantity:  req.Quantity,
		}, req.UserId)
		if err != nil {
			return nil, fmt.Errorf("Error in decrement of quantity :%w", err)
		}
	}
	prod, err := ProductClient.Get(ctx, &pb.ProductIdRequest{Id: res.ProductID})
	if err != nil {
		return nil, fmt.Errorf("Error in getting product client:%w", err)
	}
	return &pb.AddToCartResponse{Product: prod, Quantity: req.Quantity}, nil
}

/*
CreateCart(context.Context, *CartRequest) (*CartResponse, error)
    Get(*CartRequest, CartService_GetServer) error
    AddtoCart(context.Context, *AddTOCartRequest) (*AddToCartResponse, error)
    Delete(context.Context, *AddTOCartRequest) (*AddToCartResponse, error)
    UpdateQuantity(context.Context, *UpdateQuantityRequest) (*CartResponse, error)
    mustEmbedUnimplementedCartServiceServer()
*/
