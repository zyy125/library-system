package request
type CreateReservationRequest struct {
    BookID uint64 `json:"book_id" binding:"required"`
}