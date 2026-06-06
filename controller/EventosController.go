package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"notification-service/models"

	redisclient "notification-service/redis"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
POST /eventos

Recibe tipo, payload y usuarioID en el body
Guarda el evento en PostgreSQL
Publica el evento en Redis para que el worker lo procese
Responde con el evento creado
*/
var DB *gorm.DB

// @Summary      Publicar evento
// @Description  Publica un evento en Redis y lo guarda en la BD
// @Tags         Eventos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        evento  body      models.Evento  true  "Datos del evento"
// @Success      200     {object}  models.Evento
// @Failure      400     {object}  map[string]string
// @Router       /eventos [post]
func PublicarEvento(c *gin.Context) {

	var evento models.Evento

	if err := c.BindJSON(&evento); err != nil {
		c.JSON(400, gin.H{"error": "Los datos son incorrectos"})
		return
	}
	DB.Create(&evento)
	ctx := context.Background()
	eventoJson, err := json.Marshal(evento)

	if err != nil {
		fmt.Println("Error al serializar a JSON:", err)
		return
	}

	redisclient.RDB.Publish(ctx, "eventos", eventoJson)

	c.JSON(200, evento)
}

/*
 1. Obtener el usuarioID del token con c.GetUint("id")

2. Buscar todos los eventos donde usuario_id = usuarioID
3. Ordenar por fecha más reciente
4. Responder con la lista
*/

// @Summary      Obtener eventos
// @Description  Retorna todos los eventos del usuario autenticado
// @Tags         Eventos
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Evento
// @Failure      400  {object}  map[string]string
// @Router       /eventos [get]
func ObtenerEventos(c *gin.Context) {
	idUsuario := c.GetUint("id")

	if idUsuario <= 0 {
		c.JSON(400, gin.H{"error": "id incorrecto"})
		return
	}

	var eventos []models.Evento
	resulado := DB.Where("usuario_id = ?", idUsuario).Order("created_at desc").Find(&eventos)

	if resulado.Error != nil {
		c.JSON(404, gin.H{"error": "No se encotro eventos asociados con el usuario"})
		return
	}

	c.JSON(200, eventos)

}

/*
1. Obtener el id del evento de la URL
2. Buscar el evento en la BD
3. Si no existe → 404
4. Verificar que el evento pertenece al usuario autenticado
5. Actualizar leido = true con DB.Model(&evento).Update("leido", true)
6. Responder con mensaje de éxito */

// @Summary      Marcar evento como leído
// @Description  Actualiza el campo leido a true
// @Tags         Eventos
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID del evento"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /eventos/{id} [patch]
func MarcarLeido(c *gin.Context) {

	idUsuario := c.GetUint("id")
	idEvento := c.Param("id")
	var evento models.Evento

	resultado := DB.First(&evento, idEvento)

	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "No se encotro ningun evento relacionado a ese id de evento"})
		return
	}

	if evento.UsuarioID != idUsuario {
		c.JSON(403, gin.H{"error": "El evento no pertenece a este usuario"})
		return
	}

	DB.Model(&evento).Update("leido", true)

	c.JSON(200, gin.H{"message": "Se actualizo a leido correctamente"})

}
