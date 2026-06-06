package models

import "gorm.io/gorm"

/*
Usuario {
  gorm.Model
  Nombre   string
  Email    string
  Password string
} */

type Usuario struct {
	gorm.Model
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
