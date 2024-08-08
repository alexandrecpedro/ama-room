package api

// TYPE STRUCTURES
// (a) MessageMessageReactionIncreased
type MessageMessageReactionIncreased struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
}

// (b) MessageMessageReactionDecreased
type MessageMessageReactionDecreased struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
}

// (c) MessageMessageAnswered
type MessageMessageAnswered struct {
	ID string `json:"id"`
}

// (d) MessageMessageCreated
type MessageMessageCreated struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// (e) Message
type Message struct {
	Kind  string `json:"kind"`
	Value any    `json:"value"`
	// "-" => JSON package will not encode this value
	RoomID string `json:"-"`
}
