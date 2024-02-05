package models

type User struct {
	Id       uint64 `gorm:"primaryKey;column:id"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	Mobile   uint32 `gorm:"column:mobile"`
	Password string `gorm:"column:password"`
	Is_Admin bool   `gorm:"column:is_admin"`
}

type Admins struct {
	Id       uint64 `gorm:"primaryKey"`
	Name     string
	Email    string
	Mobile   uint32
	Password string
	Is_Admin bool
}
