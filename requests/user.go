package requests

type VerifyRequest struct {
	ServiceUserID string `json:"serviceUserId" binding:"required"`
	Code          string `json:"code" binding:"required"`
}

type CreateRequest struct {
	Username string `json:"username" binding:"required"`
}

type UpdateRequest struct {
	Username    string `json:"username" binding:"required"`
	OldUsername string `json:"oldUsername" binding:"required"`
	Code        string `json:"code" binding:"required"`
}

type CheckValidRequest struct {
	ServiceID uint   `json:"serviceId" binding:"required"`
	UserID    string `json:"userId" binding:"required"`
}
