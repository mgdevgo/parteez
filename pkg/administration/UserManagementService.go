package administration

type UserManagementService interface {
	CreateUser(username string, password string, role int) (string, error)
	GrantRole(userID string, role int) error
}

type UserManagementLocalService struct {
}

func NewUserManagementService() *UserManagementLocalService {
	return &UserManagementLocalService{}
}
