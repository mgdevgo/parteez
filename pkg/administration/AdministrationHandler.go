package administration

type AdministrationHandler struct {
	eventManagementService EventManagementService
	userManagementService  UserManagementService
}

func NewAdministrationHandler(eventManagementService EventManagementService, userManagementService UserManagementService) *AdministrationHandler {
	return &AdministrationHandler{
		eventManagementService: eventManagementService,
		userManagementService:  userManagementService,
	}
}
