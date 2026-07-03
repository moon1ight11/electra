package requesthandlers

import "electra/internal/api/handlers/interfaces"

type RequestHandler struct {
	requestService interfaces.RequestService
}

func NewRequestHandler(requestService interfaces.RequestService) *RequestHandler {
	return &RequestHandler{requestService: requestService}
}
