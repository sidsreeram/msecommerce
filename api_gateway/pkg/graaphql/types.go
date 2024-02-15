package graaphql

import (
	"context"
	"fmt"
	"io"
	"log"

	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/msecommerce/api_gateway/pkg/authorize"
	"github.com/msecommerce/api_gateway/pkg/middleware"
	"github.com/sidsreeram/msproto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	UsersrvCon    pb.UserServiceClient
	ProductsrvCon pb.ProductServiceClient
	CartsrvCon    pb.CartServiceClient
	Secret        []byte
)

func RetrieveSecret(secretString string) {
	Secret = []byte(secretString)

}
func Initialize(userconn pb.UserServiceClient, prodconn pb.ProductServiceClient, cartconn pb.CartServiceClient) {
	UsersrvCon = userconn
	ProductsrvCon = prodconn
	CartsrvCon = cartconn

}

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "user",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"mobile": &graphql.Field{
				Type: graphql.Int,
			},
			"password": &graphql.Field{
				Type: graphql.String,
			},
			"isadmin": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)
var ProductType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"quantity": &graphql.Field{
				Type: graphql.Int,
			},
			"price": &graphql.Field{
				Type: graphql.Int,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"instock": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)
var CartItemType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "cartItem",
		Fields: graphql.Fields{
			"product_id": &graphql.Field{
				Type: graphql.Int,
			},
			"product": &graphql.Field{
				Type: ProductType,
			},
			"quantity": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
var RootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: middleware.AdminMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					return UsersrvCon.Getuser(context.Background(), &pb.UserRequest{Id: uint64(p.Args["id"].(int))})
				}),
			},
			"admin": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: middleware.AdminMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					return UsersrvCon.GetAdmin(context.Background(), &pb.UserRequest{Id: uint64(p.Args["id"].(int))})
				}),
			},
			"userdetails": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: middleware.UserMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					return UsersrvCon.Getuser(context.Background(), &pb.UserRequest{Id: uint64(p.Args["id"].(int))})
				}),
			},
			"allusers": &graphql.Field{
				Type: graphql.NewList(UserType), // Assuming UserType is a list type
				Resolve: middleware.AdminMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					stream, err := UsersrvCon.GetAllUsers(context.Background(), &emptypb.Empty{})
					if err != nil {
						return nil, fmt.Errorf("unable to get all users @graphql: %w", err)

					}

					var users []map[string]interface{}
					for {
						user, err := stream.Recv()
						if err == io.EOF {
							break
						}
						if err != nil {
							return nil, fmt.Errorf("error in receiving all users: %w", err)
						}

						userMap := map[string]interface{}{
							"id":     user.Id,
							"name":   user.Name,
							"email":  user.Email,
							"mobile": user.Mobile,
						}

						users = append(users, userMap)
					}

					return users, nil
				}),
			},
			"product": &graphql.Field{
				Type: ProductType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					productdetails,err := ProductsrvCon.Get(context.Background(), &pb.ProductIdRequest{
						Id: uint64(p.Args["id"].(int)),
					})
					if err !=nil {
						return nil ,fmt.Errorf("error in getting product details : %w",err)
					}
					
					log.Println(productdetails.Name)
					return productdetails ,nil
				},
			},
			"allproducts": &graphql.Field{
				Type: graphql.NewList(ProductType), // Assuming ProductType is a list type
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					stream, err := ProductsrvCon.GetAll(context.Background(), &emptypb.Empty{})
					if err != nil {
						return nil, fmt.Errorf("unable to get all products @graphql: %w", err)
					}

					var products []map[string]interface{}
					for {
						product, err := stream.Recv()
						if err == io.EOF {
							break
						}
						if err != nil {
							return nil, fmt.Errorf("error receiving product: %w", err)
						}

						productMap := map[string]interface{}{
							"id":          product.GetId(),
							"name":        product.GetName(),
							"quantity":    product.GetQuantity(),
							"price":       product.GetPrice(),
							"description": product.GetDescription(),
							"instock":     product.GetInstock(),
						}

						products = append(products, productMap)
					}

					return products, nil
				},
			},
			"cart": &graphql.Field{
				Type: graphql.NewList(CartItemType),
				Resolve: middleware.UserMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					userId := p.Context.Value("userID").(uint64)
					req := &pb.CartRequest{UserId: userId}
					stream, err := CartsrvCon.Get(context.Background(), req)
					if err != nil {
						return nil, fmt.Errorf("unable to get all products @graphql: %w", err)
					}

					var cartItem []map[string]interface{}
					for {
						cart, err := stream.Recv()
						if err == io.EOF {
							break
						}
						if err != nil {
							return nil, fmt.Errorf("error receiving product: %w", err)
						}

						cartMap := map[string]interface{}{
							"product":  cart.GetProduct(),
							"quantity": cart.GetQuantity(),
						}

						cartItem = append(cartItem, cartMap)
					}
					return cartItem, nil

				}),
			},
		},
	},
)
var Mutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"signup": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"mobile": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := UsersrvCon.UserSignup(context.Background(), &pb.SignupRequest{
						Name:     p.Args["name"].(string),
						Email:    p.Args["email"].(string),
						Mobile:   uint32(p.Args["mobile"].(int)),
						Password: p.Args["password"].(string),
					})
					if err != nil {
						return nil, fmt.Errorf("error in passing arguments for user signup :%w", err)
					}

					_, err = CartsrvCon.CreateCart(context.TODO(), &pb.CartRequest{UserId: user.Id})
					if err != nil {
						return nil , fmt.Errorf("error in creating cart for user: %v",err)
					} 

					w := p.Context.Value("httpResponseWriter").(http.ResponseWriter)
					// ctx := context.WithValue(context.Background(), "httpResponseWriter", w)

					tokenstr, err := authorize.GenerateJwt(user.Id, user.IsAdmin, Secret)
					if err != nil {
						return nil, fmt.Errorf("error in generating jwt token at signup :%w", err)

					}

					http.SetCookie(w, &http.Cookie{
						Name:  "jwtToken",
						Value: tokenstr,
						Path:  "/",
					})

					return user, nil
				},
			},
			"loginuser": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := UsersrvCon.UserLogin(context.Background(), &pb.LoginRequest{
						Email:    p.Args["email"].(string),
						Password: p.Args["password"].(string),
						IsAdmin:  false,
					})
					
					if err != nil {
						return nil, fmt.Errorf("error in passing paramter into userlogin :%w", err)
					}
					w := p.Context.Value("httpResponseWriter").(http.ResponseWriter)

					tokenString, err := authorize.GenerateJwt((user.Id), false, Secret)

					if err != nil {
						return nil, fmt.Errorf("error in generating jwt :%w", err)
					}
					http.SetCookie(w, &http.Cookie{
						Name:  "jwtToken",
						Value: tokenString,
						Path:  "/",
					})
				

					return user, nil
				},
			},
			"loginadmin": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					admin, err := UsersrvCon.UserLogin(context.Background(), &pb.LoginRequest{
						Email:    p.Args["email"].(string),
						Password: p.Args["password"].(string),
						IsAdmin:  true,
					})

					if err != nil {
						return nil, fmt.Errorf("error in passing parameter into adminlogin :%w", err)
					}

					w := p.Context.Value("httpResponseWriter").(http.ResponseWriter)
					tokenString, err := authorize.GenerateJwt(admin.Id, true, Secret)
					if err != nil {
						return nil, fmt.Errorf("error in generating jwt :%w", err)
					}

					http.SetCookie(w, &http.Cookie{
						Name:  "jwtToken",
						Value: tokenString,
						Path:  "/",
					})

					return admin, nil
				},
			},
			"addadmin": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"mobile": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: middleware.AdminMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					admin, err := UsersrvCon.AddAdmin(context.Background(), &pb.SignupRequest{
						Name:     p.Args["name"].(string),
						Email:    p.Args["email"].(string),
						Password: p.Args["password"].(string),
						Mobile:   uint32(p.Args["mobile"].(int)),
					})
					if err != nil {
						return nil, fmt.Errorf("error in adding new admin :%w", err)
					}
					return admin, nil
				}),
			},
			"addproduct": &graphql.Field{
				Type: ProductType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"quantity": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"price": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"instock": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Boolean),
					},
				},
				Resolve: middleware.AdminMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					products, err := ProductsrvCon.Add(context.Background(), &pb.ProductRequest{
						Name:        p.Args["name"].(string),
						Quantity:    uint64(p.Args["quantity"].(int)),
						Price:       uint64(p.Args["price"].(int)),
						Description: p.Args["description"].(string),
						Instock:     p.Args["instock"].(bool),
					})
					if err != nil {
						return nil, fmt.Errorf("error in passing arguments to add products :%w", err)
					}
					return products, nil
				},
				),
			},
			"updatestock": &graphql.Field{
				Type: ProductType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"quantity": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"price": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"increase": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Boolean),
					},
				},
				Resolve: middleware.AdminMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					return ProductsrvCon.Update(context.Background(), &pb.UpdateProductRequest{
						Id:        uint64(p.Args["id"].(int)),
						Quantity:  uint64(p.Args["quantity"].(int)),
						Price:     uint64(p.Args["price"].(int)),
						Increased: p.Args["increase"].(bool),
					})
				}),
			},
			"addtocart": &graphql.Field{
				Type: CartItemType,
				Args: graphql.FieldConfigArgument{
					"product_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"quantity": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: middleware.UserMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					userIDVal, ok := p.Context.Value("userID").(uint64)
					if !ok {
						return nil, fmt.Errorf("userID not set in context or not of type uint64")
					}

					cart, err := CartsrvCon.AddtoCart(context.TODO(), &pb.AddTOCartRequest{
						UserId:    userIDVal,
						ProductId: uint64(p.Args["product_id"].(int)),
						Quantity:  uint64(p.Args["quantity"].(int)),
					})
					
					if err != nil {
						return nil, fmt.Errorf("error in passing arguments to cart :%w", err)
					}
					return cart, nil

				}),
			},
			"removefromcart": &graphql.Field{
				Type: CartItemType,
				Args: graphql.FieldConfigArgument{
					"product_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: middleware.UserMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					userIdval := p.Context.Value("userID").(uint64)
					cart, err := CartsrvCon.Delete(context.TODO(), &pb.AddTOCartRequest{
						UserId:    uint64(userIdval),
						ProductId: uint64(p.Args["product_id"].(int)),
					})
					if err != nil {
						return nil, fmt.Errorf("error in removing from cart :%w", err)

					}
				

					return cart, nil
				}),
			},
			"updateCartItemQty": &graphql.Field{
				Type: CartItemType,
				Args: graphql.FieldConfigArgument{
					"product_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"quantity": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"isIncreasing": &graphql.ArgumentConfig{
						Type: graphql.Boolean,
					},
				},
				Resolve: middleware.UserMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
					userIdVal := p.Context.Value("userID").(uint64)
					cart, err := CartsrvCon.UpdateQuantity(context.TODO(), &pb.UpdateQuantityRequest{
						UserId:      uint64(userIdVal),
						ProductId:   uint64(p.Args["product_id"].(uint64)),
						Quantity:    uint64(p.Args["quantity"].(uint64)),
						IsIncreased: p.Args["isIncreasing"].(bool),
					})
					if err != nil {
						return nil, fmt.Errorf("error in accessing params of update qauntity :%w", err)
					}
					return cart, nil
				}),
			},
		},
	},
)
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    RootQuery,
	Mutation: Mutation,
})
