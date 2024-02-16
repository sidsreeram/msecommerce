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
		return nil, fmt.Errorf("error occured in creating user's cart:%w", err)
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
	
	productStream, err := c.Product.Client.GetMultiple(context.Background(), &pb.ProductMultipleRequest{Ids: ids})
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
	if c == nil {
        log.Println("CartServiceServer is nil")
        return nil,errors.New("error is heree")
    }
	productResult, err := c.Product.Client.Get(context.TODO(), &pb.ProductIdRequest{Id: req.ProductId})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	
    //    if productResult.Name == ""{
	// 	return nil,fmt.Errorf("sorrry the product doesn't exist")
	//    }


	item, err := c.Adapter.AddToCart(models.CartItems{
		ProductID: req.ProductId,
		Quantity:  req.Quantity,
	}, req.UserId)
  log.Println(productResult.Name)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &pb.AddToCartResponse{
		Product: &pb.ProductResponse{
		 Id: productResult.Id,
		 Name: productResult.Name,
		 Price: productResult.Price*item.Quantity,
		 Description: productResult.Description,
		 Instock: productResult.Instock,
		},
		Quantity: item.Quantity,
	}, nil
}

func (c *CartServiceServer) Delete(ctx context.Context, req *pb.AddTOCartRequest) (*pb.AddToCartResponse, error) {
	productResult, err := c.Product.Client.Get(context.TODO(), &pb.ProductIdRequest{Id: req.ProductId})
	log.Println(productResult)
	if err != nil {
		return nil, fmt.Errorf("error in geting productId :%w", err)
	}
	if productResult.Name == "" {
		return nil, fmt.Errorf(" product doesn't find :%w", err)
	}


	cartitems,err := c.Adapter.DeleteCartItem(models.CartItems{ProductID: productResult.Id}, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("error in Deleting product :%w", err)
	}
	log.Println(cartitems.Quantity)
	return &pb.AddToCartResponse{Product: &pb.ProductResponse{
		Id: productResult.Id,
		Name: productResult.Name,
		Price: productResult.Price*cartitems.Quantity,
		Description: productResult.Description,
		Instock:productResult.Instock,
		Quantity: cartitems.Quantity,
	},Quantity: cartitems.Quantity,}, nil
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
			return nil, fmt.Errorf("error in incrementing quantity:%w", err)
		}
	} else {
		res, err = c.Adapter.DecrementQuantity(models.CartItems{
			ProductID: req.ProductId,
			Quantity:  req.Quantity,
		}, req.UserId)
		if err != nil {
			return nil, fmt.Errorf("error in decrement of quantity :%w", err)
		}
	}
	prod, err := c.Product.Client.Get(ctx, &pb.ProductIdRequest{Id: res.ProductID})
	if err != nil {
		return nil, fmt.Errorf("error in getting product client:%w", err)
	}
	return &pb.AddToCartResponse{Product: &pb.ProductResponse{
		Id: prod.Id,
		Name: prod.Name,
		Quantity: req.Quantity,
		Price: prod.Price*req.Quantity,
		Description: prod.Description,
		Instock: prod.Instock,
	}, Quantity: req.Quantity}, nil
}

/*
CreateCart(context.Context, *CartRequest) (*CartResponse, error)
    Get(*CartRequest, CartService_GetServer) error
    AddtoCart(context.Context, *AddTOCartRequest) (*AddToCartResponse, error)
    Delete(context.Context, *AddTOCartRequest) (*AddToCartResponse, error)
    UpdateQuantity(context.Context, *UpdateQuantityRequest) (*CartResponse, error)
    mustEmbedUnimplementedCartServiceServer()
*/
