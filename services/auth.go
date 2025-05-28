package services

import (
	"github.com/vivasoft-ltd/go-ems/consts"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userSvc  domain.UserService
	tokenSvc domain.TokenService
}

func NewAuthServiceImpl(userSvc domain.UserService, tokenSvc domain.TokenService) *AuthServiceImpl {
	return &AuthServiceImpl{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
	}
}

func (authSvc *AuthServiceImpl) Login(req *types.LoginReq) (*types.LoginResp, error) {
	// user fetch user by email
	user, err := authSvc.userSvc.ReadUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// compare the entered password with the one in the db
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errutil.ErrInvalidLoginCredentials
	}

	token, err := authSvc.tokenSvc.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	if err := authSvc.tokenSvc.StoreTokenUUID(token); err != nil {
		logger.Error(err)
		return nil, err
	}

	userInfo := &types.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RoleID:    user.RoleID,
		Role:      consts.RoleMap[user.RoleID],
	}

	resp := &types.LoginResp{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         userInfo,
	}

	go func() {
		if err := authSvc.userSvc.StoreInCache(userInfo); err != nil {
			logger.Error(err)
		}
	}()

	return resp, nil
}

func (authSvc *AuthServiceImpl) VerifyAccessToken(accessToken string) (*types.UserInfo, *types.Token, error) {
	token, err := authSvc.tokenSvc.ParseAccessToken(accessToken)
	if err != nil {
		return nil, nil, err
	}

	cachedUserID, err := authSvc.tokenSvc.ReadUserIDFromAccessTokenUUID(token.AccessUuid)
	if err != nil {
		return nil, nil, err
	}

	if cachedUserID != token.UserID {
		return nil, nil, errutil.ErrInvalidAccessToken
	}

	user, err := authSvc.userSvc.ReadUser(cachedUserID)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	userInfo := &types.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RoleID:    user.RoleID,
		Role:      consts.RoleMap[user.RoleID],
	}

	return userInfo, token, nil
}
