package requests

type CreateBindRequest struct {
	ServiceID string `json:"serviceId" binding:"required"`
	// used for hooks in service so that user doesnt have to enter an id or username on each totp request
	ServiceUserID string `json:"userId" binding:"required"`   // userId created in the service server
	Username      string `json:"username" binding:"required"` // userId created in the 2fa server
}

type FetchActiveBindRequests struct {
	userID string // 2fa user id
}

type ActionBindRequest struct {
	userID    string // 2fa user id
	serviceID string
}
