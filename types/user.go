package types

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/vivasoft-ltd/go-ems/consts"
)

type (
	CreateUserReq struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RoleID    int    `json:"role_id"`
	}

	UserInfo struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RoleID    int    `json:"role_id"`
		Role      string `json:"role,omitempty" gorm:"-"`
	}

	CurrentUser struct {
		ID          int    `json:"id"`
		Email       string `json:"email"`
		RoleID      int    `json:"role_id"`
		Role        string `json:"role"`
		AccessUuid  string `json:"access_uuid"`
		RefreshUuid string `json:"refresh_uuid"`
	}
)

func (req *CreateUserReq) Validate() error {
	return v.ValidateStruct(req,
		v.Field(&req.Email, v.Required, is.Email),
		v.Field(&req.Password, v.Required),
		v.Field(&req.FirstName, v.Required, v.Length(0, 50)),
		v.Field(&req.LastName, v.Required, v.Length(0, 50)),
		v.Field(&req.RoleID, v.Required, v.In(consts.RoleIdAttendee, consts.RoleIdAdmin, consts.RoleIdManager)),
	)
}
