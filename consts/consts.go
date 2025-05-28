package consts

const (
	RoleIdAdmin    = 1
	RoleIdManager  = 2
	RoleIdAttendee = 3

	RoleAdmin    = "ADMIN"
	RoleManager  = "MANAGER"
	RoleAttendee = "ATTENDEE"

	PermissionUserCreate = "user.create" // Permission to create a new user
	PermissionUserUpdate = "user.update" // Permission to update an existing user's information
	PermissionUserFetch  = "user.fetch"  // Permission to fetch a specific user's data
	PermissionUserList   = "user.list"   // Permission to list all users
	PermissionUserDelete = "user.delete" // Permission to delete a user

	PermissionEventCreate = "event.create" // Permission to create a new event
	PermissionEventUpdate = "event.update" // Permission to update an existing event
	PermissionEventFetch  = "event.fetch"  // Permission to fetch a specific event
	PermissionEventList   = "event.list"   // Permission to list events
	PermissionEventDelete = "event.delete" // Permission to delete an event
)

var RoleMap = map[int]string{
	RoleIdAdmin:    RoleAdmin,
	RoleIdManager:  RoleManager,
	RoleIdAttendee: RoleAttendee,
}
