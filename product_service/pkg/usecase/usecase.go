package usecase

import (
	"context"
	"fmt"
  pb"github.com/msecommerce_product_service/grpc/pb"
	"github.com/msecommerce_product_service/pkg/interfaces"
	"github.com/msecommerce_product_service/pkg/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductUsecase struct {
	Adapter interfaces.ProductAdapter
	pb.UnimplementedProductServiceServer
}
func NewProductUsecase(adapter interfaces.ProductAdapter)*ProductUsecase{
	return &ProductUsecase{
		Adapter: adapter,
	}
}
func (p *ProductUsecase) Add(ctx context.Context, req *pb.ProductRequest)(*pb.ProductResponse,error){
  if req.Name == ""{
     return &pb.ProductResponse{}, fmt.Errorf("The Product name cannot be empty")
  }
  reqmodel:= models.Product{
	Name: req.Name,
	Price: req.Price,
	Quantity: req.Quantity,
	Description: req.Description,
	InStock: req.Instock,
  }
  res,err:=p.Adapter.Add(reqmodel)
  if err!=nil{
    return	nil,fmt.Errorf("Something Wrong with add Product : %w",err)
  }
  return &pb.ProductResponse{
	Id:          res.Id,
	Name:        res.Name,
	Quantity:    res.Quantity,
	Price:       res.Price,
	Description: res.Description,
	Instock:     res.InStock,
}, err
}
func (p *ProductUsecase) Update(ctx context.Context, req *pb.UpdateProductRequest)(*pb.ProductResponse,error){
  
    currentProduct, err := p.Adapter.Get(req.Id)
    if err != nil {
        return nil, err
    }

    newQuantity := currentProduct.Quantity
    if req.Increased {
        newQuantity += req.Quantity
    } else {
        newQuantity -= req.Quantity
    }

    updatedProduct, err := p.Adapter.Update(models.UpdateProduct{
        Id: req.Id,
        Price: req.Price,
        Quantity: newQuantity,
        Increased: req.Increased,
    })
    if err != nil {
        return nil, err
    }

    return &pb.ProductResponse{
        Id: updatedProduct.Id,
        Name: updatedProduct.Name,
        Quantity: updatedProduct.Quantity,
        Price: updatedProduct.Price,
        Description: updatedProduct.Description,
        Instock: updatedProduct.InStock,
    }, nil
}
func(p *ProductUsecase) Delete(ctx context.Context,req *pb.ProductIdRequest)(*pb.SuccessResponse,error){
    msg, err := p.Adapter.Delete(req.Id)
    if err != nil {
        return nil, err
    }
    return &pb.SuccessResponse{Msg: msg}, nil
}

func (p*ProductUsecase) Get(ctx context.Context,req *pb.ProductIdRequest)(*pb.ProductResponse,error) {
    product, err := p.Adapter.Get(req.Id)
    if err != nil {
        return nil, err
    }
    return &pb.ProductResponse{
        Id: product.Id,
        Name: product.Name,
        Quantity: product.Quantity,
        Price: product.Price,
        Description: product.Description,
        Instock: product.InStock,
    }, nil
}

func (p *ProductUsecase) GetAll(empty *emptypb.Empty, stream pb.ProductService_GetAllServer) error {
    products, err := p.Adapter.GetAll()
    if err != nil {
        return err
    }
    for _, product := range products {
        response := &pb.ProductResponse{
            Id: product.Id,
            Name: product.Name,
            Quantity: product.Quantity,
            Price: product.Price,
            Description: product.Description,
            Instock: product.InStock,
        }
        if err := stream.Send(response); err != nil {
            return err
        }
    }
    return nil
}

// func(p*ProductUsecase) GetMultiple(ctx context.Context, id *pb.ProductIdRequest)([]*models.Product,error){

// }