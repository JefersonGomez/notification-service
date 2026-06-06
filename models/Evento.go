package models

import (
	"time"

	"gorm.io/gorm"
)

/* Evento {
  ID        string
  Tipo      string   // "tarea_asignada", "miembro_agregado", "estado_cambiado"
  Payload   string   // JSON con los datos del evento
  UsuarioID uint     // a quién va dirigido
  Leido     bool
  CreadoEn  time.Time
} */

type Evento struct {
	gorm.Model
	Tipo      string    `json:"tipo"`
	Payload   string    `json:"payload"`
	UsuarioID uint      `json:"usuarioID"`
	Leido     bool      `json:"leido"`
	CreadoEn  time.Time `json:"creadoEn"`
}
