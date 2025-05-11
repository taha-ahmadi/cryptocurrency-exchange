package handler

import (
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/exchanges"
)

// Handler handles HTTP requests
type Handler struct {
	Exchange *exchanges.Exchange
}

// New creates a new handler
func New(exchange *exchanges.Exchange) *Handler {
	return &Handler{Exchange: exchange}
}
