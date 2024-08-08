package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexandrecpedro/ama-room/backend/internal/store/pgstore"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// SHARED FUNCTIONS

// (a) READ ROOM
func (apiHandler apiHandler) readRoom(
	respWriter http.ResponseWriter,
	req *http.Request,
) (room pgstore.Room, rawRoomID string, roomID uuid.UUID, ok bool) {
	// Get ID
	rawRoomID = chi.URLParam(req, "room_id")
	// Verify if roomID is valid
	roomID, err := uuid.Parse(rawRoomID)
	if err != nil {
		http.Error(respWriter, ErrInvalidRoomID, http.StatusBadRequest)
		return pgstore.Room{}, "", uuid.UUID{}, false
	}

	// Check if this room exists at DB
	room, err = apiHandler.query.GetRoom(req.Context(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(respWriter, ErrRoomNotFound, http.StatusNotFound)
			return pgstore.Room{}, "", uuid.UUID{}, false
		}

		slog.Error(ErrFailedToGetRoom, "error", err)
		http.Error(respWriter, ErrSomethingWentWrong, http.StatusInternalServerError)
		return pgstore.Room{}, "", uuid.UUID{}, false
	}

	return room, rawRoomID, roomID, true
}

// (b) SEND JSON
func sendJSON(respWriter http.ResponseWriter, rawData any) {
	// Encoding JSON from response
	data, err := json.Marshal(rawData)
	if err != nil {
		http.Error(respWriter, ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	// Response return
	respWriter.Header().Set("Content-Type", "application/json")
	_, err = respWriter.Write(data)
	if err != nil {
		slog.Error(ErrFailedToReturnRegisteredRoom, "error", err)
		http.Error(respWriter, ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}
}

// (c) NOTIFY CLIENTS
func (apiHandler apiHandler) notifyClients(msg Message) {
	// Lock the map to iterate an get each client
	apiHandler.mu.Lock()
	defer apiHandler.mu.Unlock()

	// Get subscribers from a room
	subscribers, ok := apiHandler.subscribers[msg.RoomID]
	if !ok || len(subscribers) == 0 {
		return
	}

	// Notify each client
	for connection, cancel := range subscribers {
		if err := connection.WriteJSON(msg); err != nil {
			// keep the log
			slog.Error(ErrFailedToNotifyClient, "error", err)
			// cancel the connection with client
			// Ps: it is not necessary to remove from the map (the signal will be received and the client will be remove as done before)
			cancel()
		}
	}
}
