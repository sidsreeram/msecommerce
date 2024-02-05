package helpers

import "golang.org/x/crypto/bcrypt"
func HashPassword(password string)(string,error){
	pass:=[]byte(password)
	hasedpass,err:=bcrypt.GenerateFromPassword(pass,bcrypt.DefaultCost)
	return string(hasedpass),err
}
func VerifyPassword(password,checkpassword string)(error){
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(checkpassword))
}