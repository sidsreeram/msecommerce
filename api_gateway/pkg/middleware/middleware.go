package middleware

import (
	"context"
	"fmt"
	

	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/msecommerce/api_gateway/pkg/authorize"
)

var secret []byte

func InitMiddleware(secretstr string) {
	secret = []byte(secretstr)
}

func UserMiddleware(next graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		r := p.Context.Value("request").(*http.Request)
		cookiee, err := r.Cookie("jwtToken")
		if err != nil {
			return nil, err
		}
		if cookiee == nil {
			return nil, fmt.Errorf("You are not logged in please login")
		}
		ctx := p.Context
		token := cookiee.Value
		auth, err := authorize.ValidateToken(token, secret)
		if err != nil {
			return nil, fmt.Errorf("error in validating jwt : %w", err)
		}
		useridVal := auth["userID"].(uint64)

		if useridVal < 1 {
			return nil, fmt.Errorf("UserID is invalid")
		}
		ctx = context.WithValue(ctx, "userID", useridVal)
		p.Context = ctx
		return next(p)
	}
}
func AdminMiddleware(next graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		r := p.Context.Value("request").(*http.Request)
		cookiee, err := r.Cookie("jwtToken")
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		if cookiee == nil {
			return nil, fmt.Errorf("You are not logged in please login")
		}
		ctx := p.Context
		token := cookiee.Value
		auth, err := authorize.ValidateToken(token, secret)
		if err != nil {
			return nil, fmt.Errorf("error in validating jwt : %w", err)
		}

		useridVal := auth["userID"].(uint64)

		if useridVal < 1 {
			return nil, fmt.Errorf("UserID is invalid")
		}
		
		if !auth["isadmin"].(bool){
			return nil, fmt.Errorf("Ooops You don't have admin Previlages ! :%w", err)
		}

		ctx = context.WithValue(ctx, "userID", useridVal)
		p.Context = ctx
		return next(p)
	}
}
