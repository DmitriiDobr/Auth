package repository

import jwt "github.com/golang-jwt/jwt/v4"

type User struct {
	Id       int    `json:"-" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Registrator struct {
	Username string `form:"username" json:"username" binding:"required" db:"username"`
	Password string `form:"password" json:"password" binding:"required" db:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
