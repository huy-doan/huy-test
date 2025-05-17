package outputdata

import (
	userModel "github.com/huydq/test/internal/domain/model/user"
)

// UserProfileOutputData represents the output data structure for user profile
type UserProfileOutputData struct {
	User *userModel.User `json:"user"`
}
