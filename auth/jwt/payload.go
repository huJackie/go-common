package jwt

import "github.com/dgrijalva/jwt-go"

type Payload struct {
	jwt.StandardClaims
	Uid      string `json:"uid,omitempty"`      //用户id
	Account  string `json:"account,omitempty"`  //账号
	Role     int64  `json:"role,omitempty"`     //角色
	Platform int64  `json:"platform,omitempty"` //平台
}
