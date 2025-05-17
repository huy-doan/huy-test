package inputdata

// GetPermissionsByIDsInput represents the input for retrieving permissions by their IDs
type GetPermissionsByIDsInput struct {
	IDs []int `json:"ids" binding:"required"`
}
