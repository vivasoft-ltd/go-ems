package domain

import "github.com/vivasoft-ltd/go-ems/types"

type (
	AuthService interface {
		Login(req *types.LoginReq) (*types.LoginResp, error)
		VerifyAccessToken(accessToke string) (*types.UserInfo, *types.Token, error)
	}
)
