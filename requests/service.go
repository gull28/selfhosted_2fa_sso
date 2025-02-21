package requests

type UnlinkRequest struct {
	UserServiceLinkID uint `json:"id" binding:"required"`
}
