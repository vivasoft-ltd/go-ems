package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/middlewares"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/go-ems/utils/msgutil"
	"net/http"
)

type UserController struct {
	userSvc domain.UserService
}

func NewUserController(userSvc domain.UserService) *UserController {
	return &UserController{
		userSvc: userSvc,
	}
}

func (ctrl *UserController) Signup(c echo.Context) error {
	var req types.CreateUserReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	req.RoleID = consts.RoleIdAttendee
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	if err := ctrl.userSvc.CreateUser(&req); err != nil {
		if errors.Is(err, errutil.ErrUserIsAlreadyExists) {
			return c.JSON(http.StatusConflict, msgutil.UserAlreadyExists())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusCreated, msgutil.UserCreatedSuccessfully())
}

func (ctrl *UserController) GetProfile(c echo.Context) error {
	currentUser, err := middlewares.CurrentUserFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
	}

	user, err := ctrl.userSvc.ReadUser(currentUser.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, &types.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RoleID:    user.RoleID,
		Role:      consts.RoleMap[user.RoleID],
	})
}

func (ctrl *UserController) CreateUser(c echo.Context) error {
	var req types.CreateUserReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	if err := ctrl.userSvc.CreateUser(&req); err != nil {
		if errors.Is(err, errutil.ErrUserIsAlreadyExists) {
			return c.JSON(http.StatusConflict, msgutil.UserAlreadyExists())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusCreated, msgutil.UserCreatedSuccessfully())
}
