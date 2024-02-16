package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/msecommerce_product_service/pkg/interfaces"
	"github.com/msecommerce_product_service/pkg/models"
	"github.com/sidsreeram/msproto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductUsecase struct {
	Adapter interfaces.ProductAdapter
	pb.UnimplementedProductServiceServer
}

func NewProductUsecase(adapter interfaces.ProductAdapter) *ProductUsecase {
	return &ProductUsecase{
		Adapter: adapter,
	}
}
func (p *ProductUsecase) Add(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	if req.Name == "" {
		return &pb.ProductResponse{}, fmt.Errorf("the Product name cannot be empty")
	}
	reqmodel := models.Product{
		Name:        req.Name,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Description: req.Description,
		InStock:     req.Instock,
	}
	res, err := p.Adapter.Add(reqmodel)
	if err != nil {
		return nil, fmt.Errorf("something Wrong with add Product : %w", err)
	}
	log.Println(res.Id)
	return &pb.ProductResponse{
		Id:          res.Id,
		Name:        res.Name,
		Quantity:    res.Quantity,
		Price:       res.Price,
		Description: res.Description,
		Instock:     res.InStock,
	}, err
}
func (p *ProductUsecase) Update(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {

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
		Id:        req.Id,
		Price:     req.Price,
		Quantity:  newQuantity,
		Increased: req.Increased,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ProductResponse{
		Id:          updatedProduct.Id,
		Name:        updatedProduct.Name,
		Quantity:    updatedProduct.Quantity,
		Price:       updatedProduct.Price,
		Description: updatedProduct.Description,
		Instock:     updatedProduct.InStock,
	}, nil
}
func (p *ProductUsecase) Delete(ctx context.Context, req *pb.ProductIdRequest) (*pb.ProductSuccessResponse, error) {
	msg, err := p.Adapter.Delete(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ProductSuccessResponse{Msg: msg}, nil
}

func (p *ProductUsecase) Get(ctx context.Context, req *pb.ProductIdRequest) (*pb.ProductResponse, error) {
	product, err := p.Adapter.Get(req.Id)
	if err != nil {
		return nil, err
	}
	
	return &pb.ProductResponse{
		Id:          product.Id,
		Name:        product.Name,
		Quantity:    product.Quantity,
		Price:       product.Price,
		Description: product.Description,
		Instock:     product.InStock,
	}, nil
}

func (p *ProductUsecase) GetAll(empty *emptypb.Empty, stream pb.ProductService_GetAllServer) error {
	products, err := p.Adapter.GetAll()
	if err != nil {
		return err
	}
	for _, product := range products {
		response := &pb.ProductResponse{
			Id:          product.Id,
			Name:        product.Name,
			Quantity:    product.Quantity,
			Price:       product.Price,
			Description: product.Description,
			Instock:     product.InStock,
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	return nil
}
func (p *ProductUsecase) GetMultiple(ctx context.Context, req *pb.ProductMultipleRequest) (*pb.ProductMultipleResponse, error) {
	var productIDs []uint64

	for _, idReq := range req.Ids {
		productIDs = append(productIDs, idReq.Id)
	}

	// Call the adapter function to get multiple products from the database
	products, err := p.Adapter.GetMultiple(productIDs)
	if err != nil {
		return nil, err
	}

	// Convert the database models to the protobuf response and send them through the server stream
	res := &pb.ProductMultipleResponse{}
	for _, product := range products {
		productResponse := &pb.ProductResponse{
			Id:          product.Id,
			Name:        product.Name,
			Quantity:    product.Quantity,
			Price:       product.Price,
			Description: product.Description,
			Instock:     product.InStock,
		}
		res.Multipleresponse = append(res.Multipleresponse, productResponse)
	}

	return &pb.ProductMultipleResponse{
		Multipleresponse: res.Multipleresponse,
	}, nil
}
