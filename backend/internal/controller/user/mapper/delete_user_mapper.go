package mapper

// DeleteUserResponse represents the HTTP response for user deletion
type DeleteUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ToResponseMap converts the DeleteUserResponse to a map for JSON response
func (r *DeleteUserResponse) ToResponseMap() map[string]interface{} {
	return map[string]interface{}{
		"success": r.Success,
		"message": r.Message,
	}
}

// NewDeleteUserResponse creates a new DeleteUserResponse
func NewDeleteUserResponse() *DeleteUserResponse {
	return &DeleteUserResponse{
		Success: true,
		Message: "ユーザーが正常に削除されました",
	}
}

// NewDeleteUserErrorResponse creates a new error response for user deletion
func NewDeleteUserErrorResponse(err error) *DeleteUserResponse {
	return &DeleteUserResponse{
		Success: false,
		Message: "ユーザーの削除に失敗しました: " + err.Error(),
	}
}
