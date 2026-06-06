package worker

import (
	"context"
	"encoding/json"
	"notification-service/models"
	redisclient "notification-service/redis"
	"sync"

	"github.com/gorilla/websocket"
)

//definir el Hub

type Hub struct {
	Clientes map[uint]*websocket.Conn
	mu       sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clientes: make(map[uint]*websocket.Conn),
	}
}

func (h *Hub) Registrar(usuarioID uint, conn *websocket.Conn) {

	h.mu.Lock()
	h.Clientes[usuarioID] = conn
	h.mu.Unlock()

}

func (h *Hub) Eliminar(usuarioID uint) {
	h.mu.Lock()
	delete(h.Clientes, usuarioID)
	h.mu.Unlock()
}

func (h *Hub) Enviar(usuarioID uint, mensaje []byte) {
	h.mu.Lock()
	con, existe := h.Clientes[usuarioID]
	h.mu.Unlock()

	if existe {
		con.WriteMessage(websocket.TextMessage, mensaje)
	}

}

func StartWorker(hub *Hub) {
	//susbcribirse al canal de eventos de redis
	pubsub := redisclient.RDB.Subscribe(context.Background(), "eventos")

	ch := pubsub.Channel()
	// 2. goroutine que escucha en background
	go func() {
		for msg := range ch {
			// 3. deserializar el JSON
			var evento models.Evento
			json.Unmarshal([]byte(msg.Payload), &evento)
			// 4. enviar al usuario conectado
			hub.Enviar(evento.UsuarioID, []byte(msg.Payload))
		}
	}()

}
