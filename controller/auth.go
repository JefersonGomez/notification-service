package controller

import (
	"notification-service/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

/* 1. Leer el body (nombre, email, password)
2. Encriptar el password con bcrypt
3. Guardar el usuario en la BD
4. Responder con mensaje de éxito */

// @Summary Registro de usuario
// @Description Crea un nuevo usuario en sistema
// @Tags Auth
// @Accept json
// @Produce json
// @Param usuario body models.Usuario true "Datos del usuario"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /registro [post]
func Registro(c *gin.Context) {

	var usuario models.Usuario

	if err := c.BindJSON(&usuario); err != nil {
		c.JSON(400, gin.H{"error": "Los datos son incorrectos"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(usuario.Password), 10)

	if err != nil {
		c.JSON(400, gin.H{"error": "No se pudo encriptar la contrasea"})
		return
	}

	usuario.Password = string(hash)

	DB.Create(&usuario)

	c.JSON(201, gin.H{"message": "Se creo el usuario correctamente"})

}

/*
 1. Leer el body (email, password)

2. Buscar el usuario por email en la BD
3. Si no existe → 404
4. Comparar el password con bcrypt
5. Si no coincide → 401
6. Generar token JWT con el ID del usuario que expire en 24 horas
7. Responder con el token
*/

// @Summary Login
// @Description Inicia sesion y retorna un JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body object{email=string,password=string} true "Credenciales"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failture 404 {object} map[string]string
// @Router /login [post]

func Login(c *gin.Context) {

	var usuario models.Usuario
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Los datos son incorrectos"})
		return
	}

	resultado := DB.Where("email = ?", body.Email).First(&usuario)

	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "No se encontro considencia con email"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(body.Password))

	if err != nil {
		c.JSON(404, gin.H{"error": "No se encontro considencia con las contraseñas"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  usuario.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(400, gin.H{"error": "No se pudo crear el token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})

}
