package response

import "time"

type CreateReservationResponse struct {
	ID            uint64    `json:"id"`
	BookID        uint64    `json:"book_id"`
	BookTitle     string    `json:"book_title"`
	UserID        uint64    `json:"user_id"`
	Status        string    `json:"status"`
	QueuePosition int       `json:"queue_position"`
	ReservedAt    time.Time `json:"reserved_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

type GetMyReservationsResponse struct {
	Reservations []ReservationItem `json:"reservations"`
}

type ReservationItem struct {
	ID            uint64                  `json:"id"`
	Book          ReservationBookResponse `json:"book"`
	Status        string                  `json:"status"`
	QueuePosition int                     `json:"queue_position"`
	ReservedAt    time.Time               `json:"reserved_at"`
	ExpiresAt     *time.Time              `json:"expires_at"`
}

type ReservationBookResponse struct {
	ID       uint64 `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	CoverURL string `json:"cover_url"`
}


