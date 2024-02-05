package interfaces

import "github.com/msecommerce/user_service/pkg/models"

type UserAdapter interface {
	UserSignup(req models.User) (models.User, error)
	UserLogin(email, password string) (models.User, error)
	AddAdmin(admin models.Admins) (models.Admins, error)
	Getuser(id uint64) (models.User, error)
	GetAdmin(id uint64) (models.Admins, error)
	GetAllUsers() ([]models.User, error)
	GetAllAdmins() ([]models.Admins, error)
	FindByEmail(email string,isadmin bool) (string, error)
	AdminLogin(email, password string) (models.Admins ,error)
}
