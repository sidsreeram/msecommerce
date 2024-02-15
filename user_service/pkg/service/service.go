package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/msecommerce/user_service/pkg/helpers"
	"github.com/msecommerce/user_service/pkg/interfaces"
	"github.com/msecommerce/user_service/pkg/models"
	"github.com/sidsreeram/msproto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	Adapter interfaces.UserAdapter
	pb.UnimplementedUserServiceServer
}

func NewUserService(adapter interfaces.UserAdapter) *UserService {
	return &UserService{
		Adapter: adapter,
	}
}

func (u *UserService) UserSignup(ctx context.Context, req *pb.SignupRequest) (*pb.UserResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("Name Cannot be empty")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("Email cannot be empty")
	}
	if req.Mobile == 0 {
		return nil, fmt.Errorf("Mobile cannot be empty")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("Password cannot be empty")

	}

	password, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("Error in hashing password :%w", err)
	}
	reqest := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: password,
	}
	res, err := u.Adapter.UserSignup(reqest)
	if err != nil {
		return nil, fmt.Errorf("Error in User signup :%w", err)
	}
	return &pb.UserResponse{
		Id:  res.Id,
		Name:    res.Name,
		Email:   res.Email,
		Mobile:  res.Mobile,
		IsAdmin: res.Is_Admin,
	}, err
}

func (u *UserService) Getuser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := u.Adapter.Getuser(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		Id: user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Mobile: user.Mobile,
	}, nil
}

func (u *UserService) GetAdmin(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	admin, err := u.Adapter.GetAdmin(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:  admin.Id,
		Name:    admin.Name,
		Email:   admin.Email,
		Mobile:  admin.Mobile,
		IsAdmin: admin.Is_Admin,
	}, nil
}
func (u *UserService) GetAllUsers(empty *emptypb.Empty, stream pb.UserService_GetAllUsersServer) error {
	users, err := u.Adapter.GetAllUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		userResponse := &pb.UserResponse{
			Id: user.Id,
			Name:   user.Name,
			Email:  user.Email,
			Mobile: uint32(user.Mobile),
		}
		if err := stream.Send(userResponse); err != nil {
			return err
		}
	}
	return nil
}

func (u *UserService) GetAllAdmins(empty *emptypb.Empty, stream pb.UserService_GetAllAdminsServer) error {
	admins, err := u.Adapter.GetAllAdmins()
	if err != nil {
		return err
	}
	for _, admin := range admins {
		adminResponse := &pb.UserResponse{
			Id:  admin.Id,
			Name:    admin.Name,
			Email:   admin.Email,
			Mobile:  admin.Mobile,
			IsAdmin: admin.Is_Admin,
		}
		if err := stream.Send(adminResponse); err != nil {
			return err
		}
	}
	return nil
}

func (u *UserService) UserLogin(ctx context.Context, req *pb.LoginRequest) (*pb.UserResponse, error) {

	// Search for users in the user table, admins in the admin table
	if req.IsAdmin {
		// Admin login
		adminres, err := u.Adapter.AdminLogin(req.Email, req.Password)
		if err != nil {
			return nil, err
		}
		return &pb.UserResponse{
			Id:  adminres.Id,
			Name:    adminres.Name,
			Email:   adminres.Email,
			Mobile:  adminres.Mobile,
			IsAdmin: true,
		}, nil
	} else {
		// Regular user login
		pass, err := u.Adapter.FindByEmail(req.Email, false) // Search in the user table
		if err != nil {
			return nil, err
		}
		err = helpers.VerifyPassword(pass, req.Password)
		if err != nil {
			return nil, errors.New("invalid password")
		}
		userres, err := u.Adapter.UserLogin(req.Email, pass)
		log.Println(userres.Id) // Use UserLogin for regular users
		return &pb.UserResponse{
			Id:  userres.Id,
			Name:    userres.Name,
			Email:   userres.Email,
			Mobile:  userres.Mobile,
			IsAdmin: false,
		}, nil
	}
}

func (u *UserService) AddAdmin(ctx context.Context, req *pb.SignupRequest) (*pb.UserResponse, error) {
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	admin := models.Admins{
		Name:     req.Name,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: hashedPassword,
	}

	admin, err = u.Adapter.AddAdmin(admin)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:  admin.Id,
		Name:    admin.Name,
		Email:   admin.Email,
		Mobile:  admin.Mobile,
		IsAdmin: admin.Is_Admin,
	}, nil
}
