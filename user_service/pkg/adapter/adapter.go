package adapter

import (
	"fmt"

	"github.com/msecommerce/user_service/pkg/interfaces"
	"github.com/msecommerce/user_service/pkg/models"
	"gorm.io/gorm"
)

type UserAdapter struct {
	DB *gorm.DB
}

func NewUserAdapter(DB *gorm.DB) interfaces.UserAdapter {
	return &UserAdapter{DB}
}
func (u *UserAdapter) UserSignup(req models.User) (models.User, error) {
	var user models.User
	query := `INSERT INTO users (name,email,mobile,password)VALUES($1,$2,$3,$4) RETURNING id,name,email,mobile`
	err := u.DB.Raw(query, req.Name, req.Email, req.Mobile, req.Password).Scan(&user).Error
	if err != nil {
		return user, fmt.Errorf("Error in inserting user: %w", err)
	}
	return user, nil

}
func (u *UserAdapter) UserLogin(email, password string) (models.User, error) {
	var user models.User

	// Use Where to build the query and First to retrieve a single result
	err := u.DB.Where("email = ? AND password = ?", email, password).First(&user).Error
	if err != nil {
		return models.User{}, fmt.Errorf("error in UserLogin at adapter: %v", err)
	}

	return user, nil
}
func (u *UserAdapter) AdminLogin(email, password string) (models.Admins, error) {
	var admin models.Admins

	// Use Where to build the query and First to retrieve a single result
	err := u.DB.Where("email = ? AND password = ?", email, password).First(&admin).Error
	if err != nil {
		return models.Admins{}, fmt.Errorf("error in UserLogin at adapter: %v", err)
	}

	return admin, nil
}

func (u *UserAdapter) AddAdmin(admin models.Admins) (models.Admins, error) {
	if err := u.DB.Create(&admin).Error; err != nil {
		return models.Admins{}, err
	}
	return admin, nil
}

func (u *UserAdapter) Getuser(id uint64) (models.User, error) {
	var user models.User
	if err := u.DB.First(&user, id).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *UserAdapter) GetAdmin(id uint64) (models.Admins, error) {
	var admin models.Admins
	if err := u.DB.First(&admin, id).Error; err != nil {
		return models.Admins{}, err
	}
	return admin, nil
}

func (u *UserAdapter) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := u.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserAdapter) GetAllAdmins() ([]models.Admins, error) {
	var admins []models.Admins
	if err := u.DB.Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}
func (u *UserAdapter) FindByEmail(email string, isadmin bool) (string, error) {
	var user models.User
	var admin models.Admins

	if isadmin {
		if err := u.DB.Where("email = ?", email).First(&admin).Error; err == nil {
			return admin.Password, nil
		}
	}

	if err := u.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return user.Password, nil
	}
	return "", fmt.Errorf("Email not found")
}
