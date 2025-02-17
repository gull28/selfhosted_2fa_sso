package requests

type BindRequest struct {
	ServiceID string `json:"serviceId" binding:"required"`
	// used for hooks in service so that user doesnt have to enter an id or username on each totp request
	UserID   string `json:"userId" binding:"required"`   // userId created in the service server
	Username string `json:"username" binding:"required"` // userId created in the 2fa server
}

type FetchActiveBindRequests struct {
	userID string
}

type ActionBindRequest struct {
	userID    string
	serviceID string
}
