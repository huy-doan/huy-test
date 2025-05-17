package outputdata

import "github.com/huydq/test/internal/domain/model/user"

// UserListOutputData represents paginated user list results
type UserListOutputData struct {
	Users       []*user.User `json:"users"`
	TotalPages  int          `json:"total_pages"`
	TotalCount  int          `json:"total_count"`
	CurrentPage int          `json:"current_page"`
}
