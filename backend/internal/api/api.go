package api

import (
	// NATIVE PACKAGES
	// Context package
	"context"
	// JSON encoder/decoder
	"encoding/json"
	// Deal with errors
	"errors"
	// Keep Log
	"log/slog"
	// Native package to deal with HTTP
	"net/http"
	// Sync package
	"sync"

	// INTERNAL PACKAGES
	// Internal package that handles PostgreSQL database operations
	"github.com/alexandrecpedro/ama-room/backend/internal/store/pgstore"

	// EXTERNAL PACKAGES
	// Middleware for router
	"github.com/go-chi/chi/middleware"
	// Router for creating HTTP handlers
	"github.com/go-chi/chi/v5"
	// CORS for router
	"github.com/go-chi/cors"
	// UUID from Google
	"github.com/google/uuid"
	// Websocket
	"github.com/gorilla/websocket"
	// PGX v5
	"github.com/jackc/pgx/v5"
)

// Part 1: Interface structure
type apiHandler struct {
	//? TODO: use an interface on this parameter "query"
	// (a) instance of pgstore.Queries, which presumably contains database queries
	query *pgstore.Queries
	// (b) router for managing HTTP routes
	router *chi.Mux
	// (c) upgrader: upgrade HTTP request to websocket
	upgrader websocket.Upgrader
	// (d) subscribers: store all opened connections with clients
	// Ps1: map[string] for map[opened_connections]; for each connection, save context.CancelFunc (option to cancel any operation running on Go)
	// Ps2: maps is not thread safe (data race occurs)
	subscribers map[string]map[*websocket.Conn]context.CancelFunc
	// (e) mux: mutex (mutual exclusion) block the data race
	mu *sync.Mutex
}

// Part 2: Method from interface
// mandatory method for any type that implements the http.Handler interface
func (handler apiHandler) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	// delegates request handling to the chi router
	handler.router.ServeHTTP(respWriter, req)
}

// Part 3: Function that creates and returns a new HTTP handler
func NewHandler(query *pgstore.Queries) http.Handler {
	// Instantiate apiHandler type
	apiHandler := apiHandler{
		query: query,
		// Ps: CheckOrigin is a closure
		upgrader: websocket.Upgrader{CheckOrigin: func(req *http.Request) bool { return true }},
		// need to initialize the map (default = start as null)
		subscribers: make(map[string]map[*websocket.Conn]context.CancelFunc),
		mu:          &sync.Mutex{},
	}

	// Create new router
	router := chi.NewRouter()

	// Set middlewares on router
	// (a) RequestID => assigns a unique ID to each request
	// (b) Recoverer => prevents the server from crashing in case of a panic
	// (c) Logger => keep a log of all requests
	router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	// Set CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Set routes

	// -- Websockets route
	// Ps: client will connect to an specific room and
	// receives any changes that happens
	router.Get("/subscribe/{room_id}", apiHandler.handleSubscribe)

	// -- API routes
	router.Route("/api", func(apiRouter chi.Router) {
		// Set subroutes
		//* Ps: all routes will be associated with an specific room
		// (a) Rooms
		apiRouter.Route("/rooms", func(roomRouter chi.Router) {
			// i. Register a room
			roomRouter.Post("/", apiHandler.handleCreateRoom)
			// ii. Get all rooms
			roomRouter.Get("/", apiHandler.handleGetRooms)

			// (b) Specific room
			roomRouter.Route("/{room_id}", func(r chi.Router) {
				// i. Get a room
				roomRouter.Get("/", apiHandler.handleGetRoom)

				// (c) Room Messages
				roomRouter.Route("/messages", func(messageRoomRouter chi.Router) {
					// i. Register message from a room
					messageRoomRouter.Post("/", apiHandler.handleCreateRoomMessage)
					// ii. Get all messages from a room
					messageRoomRouter.Get("/", apiHandler.handleGetRoomMessages)

					// (d) Specific Message
					messageRoomRouter.Route("/{message_id}", func(specMessageRoomRouter chi.Router) {
						// i. Get room message
						specMessageRoomRouter.Get("/", apiHandler.handleGetRoomMessage)
						// ii. Mark an specific room message as answered
						specMessageRoomRouter.Patch("/answer", apiHandler.handleMarkRoomMessageAsAnswered)
						// iii. React to an specific room message
						specMessageRoomRouter.Patch("/react", apiHandler.handleReactToRoomMessage)
						// iv. Delete a reaction from a room message
						specMessageRoomRouter.Delete("/react", apiHandler.handleRemoveReactFromRoomMessage)
					})
				})
			})
		})
	})

	// Assigns the router r to the apiHandler structure
	apiHandler.router = router
	// Returns the handler
	return apiHandler
}

// Part 4: CONSTANT VARIABLES
// (a) Message object constants
const (
	MessageKindMessageAnswered          = "message_answered"
	MessageKindMessageCreated           = "message_created"
	MessageKindMessageReactionDecreased = "message_reaction_decreased"
	MessageKindMessageReactionIncreased = "message_reaction_increased"
)

// Part 6: Functions related to each HTTP method
// -- WEBSOCKETS ROUTE
// i. GET: handleSubscribe
func (apiHandler apiHandler) handleSubscribe(respWriter http.ResponseWriter, req *http.Request) {
	// Verify if a room exists
	_, rawRoomID, _, ok := apiHandler.readRoom(respWriter, req)
	if !ok {
		return
	}

	// Upgrade connection with client
	connection, err := apiHandler.upgrader.Upgrade(respWriter, req, nil)
	if err != nil {
		// the client is not able to upgrade to websocket
		slog.Warn(ErrUpgradeToWebsocketConnection, "error", err)
		// return error to user
		http.Error(respWriter, ErrUpgradeToWebsocketConnection, http.StatusBadRequest)
		return
	}

	// Clean up connection with client (clean resources)
	defer connection.Close()

	connectionContext, cancel := context.WithCancel(req.Context())

	// Store this connection on the connection pool
	apiHandler.mu.Lock()
	if _, ok := apiHandler.subscribers[rawRoomID]; !ok {
		// initialize the map
		apiHandler.subscribers[rawRoomID] = make(map[*websocket.Conn]context.CancelFunc)
	}
	// Keep logs
	//? TODO: check the best form to get client IP (change from req.RemoteAddr)
	slog.Info("New client connected!", "room_id", rawRoomID, "client_ip", req.RemoteAddr)
	// Create this room on the map
	apiHandler.subscribers[rawRoomID][connection] = cancel
	// Enable to touch on mutex again
	apiHandler.mu.Unlock()

	// Keep this function running until it ends by client or server
	<-connectionContext.Done()

	// Remove our subscribe from subscriber list
	// Ps: the context has being cancel
	apiHandler.mu.Lock()
	delete(apiHandler.subscribers[rawRoomID], connection)
	apiHandler.mu.Unlock()
}

// -- API ROUTES
// (a) ROOMS
// i. POST: handleCreateRoom
func (apiHandler apiHandler) handleCreateRoom(respWriter http.ResponseWriter, req *http.Request) {
	// body
	type _body struct {
		Theme string `json:"theme"`
	}

	// Create variable from _body type
	var body _body
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		// return error to user
		http.Error(respWriter, ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	// Insert room at DB
	roomID, err := apiHandler.query.InsertRoom(req.Context(), body.Theme)
	if err != nil {
		// keep log
		slog.Error(ErrFailedToRegisterRoom, "error", err)
		http.Error(respWriter, ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	// Response type to user
	type response struct {
		ID string `json:"id"`
	}

	sendJSON(respWriter, response{ID: roomID.String()})
}

// ii. GET MANY: handleGetRooms
func (apiHandler apiHandler) handleGetRooms(respWriter http.ResponseWriter, req *http.Request) {
	rooms, err := apiHandler.query.GetRooms(req.Context())
	if err != nil {
		slog.Error(ErrFailedToGetRooms, "error", err)
		http.Error(respWriter, ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	if rooms == nil {
		rooms = []pgstore.Room{}
	}

	sendJSON(respWriter, rooms)
}

// iii. GET ONE: handleGetRoom
func (apiHandler apiHandler) handleGetRoom(respWriter http.ResponseWriter, req *http.Request) {
	room, _, _, ok := apiHandler.readRoom(respWriter, req)
	if !ok {
		return
	}

	sendJSON(respWriter, room)
}

// (b) ROOM MESSAGES
// i. POST: handleCreateRoomMessage
func (apiHandler apiHandler) handleCreateRoomMessage(respWriter http.ResponseWriter, req *http.Request) {
	// Verify if a room exists
	_, rawRoomID, roomID, ok := apiHandler.readRoom(respWriter, req)
	if !ok {
		return
	}

	// body
	type _body struct {
		Message string `json:"message"`
	}

	// Create variable from _body type
	var body _body
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		// return error to user
		http.Error(respWriter, ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	messageID, err := apiHandler.query.InsertMessage(req.Context(), pgstore.InsertMessageParams{RoomID: roomID, Message: body.Message})
	if err != nil {
		// log the error
		slog.Error(ErrFailedToInsertMessage, "error", err)
		// return error to user
		http.Error(respWriter, ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	// Response type to user
	type response struct {
		ID string `json:"id"`
	}

	sendJSON(respWriter, response{ID: messageID.String()})

	// Notify all clients asynchronously
	go apiHandler.notifyClients(Message{
		Kind:   MessageKindMessageCreated,
		RoomID: rawRoomID,
		Value: MessageMessageCreated{
			ID:      messageID.String(),
			Message: body.Message,
		},
	})
}

// ii. GET MANY: handleGetRooms
func (apiHandler apiHandler) handleGetRoomMessages(respWriter http.ResponseWriter, req *http.Request) {
	_, _, roomID, ok := apiHandler.readRoom(respWriter, req)
	if !ok {
		return
	}

	messages, err := apiHandler.query.GetRoomMessages(req.Context(), roomID)
	if err != nil {
		slog.Error(ErrFailedToGetRoomMessages, "error", err)
		http.Error(respWriter, ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	if messages == nil {
		messages = []pgstore.Message{}
	}

	sendJSON(respWriter, messages)
}

// (c) SPECIFIC ROOM MESSAGE
// i. GET ONE: handleGetRoomMessage
func (apiHandler apiHandler) handleGetRoomMessage(respWriter http.ResponseWriter, req *http.Request) {
	_, _, _, ok := apiHandler.readRoom(respWriter, req)
	if !ok {
		return
	}

	rawMessageID := chi.URLParam(req, "message_id")
	messageID, err := uuid.Parse(rawMessageID)
	if err != nil {
		http.Error(respWriter, ErrInvalidMessageID, http.StatusBadRequest)
		return
	}

	messages, err := apiHandler.query.GetMessage(req.Context(), messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(respWriter, ErrMessageNotFound, http.StatusNotFound)
			return
		}

		slog.Error(ErrFailedToGetRoomMessage, "error", err)
		http.Error(respWriter, ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	sendJSON(respWriter, messages)
}

// ii. PATCH: handleMarkRoomMessageAsAnswered
func (apiHandler apiHandler) handleMarkRoomMessageAsAnswered(respWriter http.ResponseWriter, req *http.Request) {
	_, rawRoomID, _, ok := apiHandler.readRoom(respWriter, req)
	if !ok {
		return
	}

	rawMessageID := chi.URLParam(req, "message_id")
	messageID, err := uuid.Parse(rawMessageID)
	if err != nil {
		http.Error(respWriter, ErrInvalidMessageID, http.StatusBadRequest)
		return
	}

	err = apiHandler.query.MarkMessageAsAnswered(req.Context(), messageID)
	if err != nil {
		slog.Error(ErrFailedToMarkAsAnswered, "error", err)
		http.Error(respWriter, ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	// Response return
	respWriter.WriteHeader(http.StatusOK)

	go apiHandler.notifyClients(Message{
		Kind:   MessageKindMessageAnswered,
		RoomID: rawRoomID,
		Value: MessageMessageAnswered{
			ID: rawMessageID,
		},
	})

}

// iii. PATCH: handleReactToRoomMessage
func (apiHandler apiHandler) handleReactToRoomMessage(respWriter http.ResponseWriter, req *http.Request) {
	_, rawRoomID, _, ok := apiHandler.readRoom(respWriter, req)
	if !ok {
		return
	}

	rawMessageID := chi.URLParam(req, "message_id")
	messageID, err := uuid.Parse(rawMessageID)
	if err != nil {
		http.Error(respWriter, ErrInvalidMessageID, http.StatusBadRequest)
		return
	}

	count, err := apiHandler.query.RemoveReactionFromMessage(req.Context(), messageID)
	if err != nil {
		slog.Error(ErrFailedToReactToMessage, "error", err)
		http.Error(respWriter, ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	type response struct {
		Count int64 `json:"count"`
	}

	sendJSON(respWriter, response{Count: count})

	go apiHandler.notifyClients(Message{
		Kind:   MessageKindMessageReactionIncreased,
		RoomID: rawRoomID,
		Value: MessageMessageReactionIncreased{
			ID:    rawMessageID,
			Count: count,
		},
	})
}

// iv. DELETE: handleRemoveReactFromRoomMessage
func (apiHandler apiHandler) handleRemoveReactFromRoomMessage(respWriter http.ResponseWriter, req *http.Request) {
	_, rawRoomID, _, ok := apiHandler.readRoom(respWriter, req)
	if !ok {
		return
	}

	rawMessageID := chi.URLParam(req, "message_id")
	messageID, err := uuid.Parse(rawMessageID)
	if err != nil {
		http.Error(respWriter, ErrInvalidMessageID, http.StatusBadRequest)
		return
	}

	count, err := apiHandler.query.RemoveReactionFromMessage(req.Context(), messageID)
	if err != nil {
		slog.Error(ErrFailedToRemoveReaction, "error", err)
		http.Error(respWriter, ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	type response struct {
		Count int64 `json:"count"`
	}

	sendJSON(respWriter, response{Count: count})

	go apiHandler.notifyClients(Message{
		Kind:   MessageKindMessageReactionDecreased,
		RoomID: rawRoomID,
		Value: MessageMessageReactionDecreased{
			ID:    rawMessageID,
			Count: count,
		},
	})
}
