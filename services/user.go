package services

import (
	"errors"
	"fmt"
	"github.com/vivasoft-ltd/go-ems/config"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/go-ems/models"
	"github.com/vivasoft-ltd/go-ems/types"
	"github.com/vivasoft-ltd/go-ems/utils/errutil"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

type UserServiceImpl struct {
	userRepo domain.UserRepository
	redisSvc *RedisService
}

func NewUserServiceImpl(userRepo domain.UserRepository, redisSvc *RedisService) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
		redisSvc: redisSvc,
	}
}

func (svc *UserServiceImpl) CreateUser(req *types.CreateUserReq) error {
	isExists, err := svc.IsEmailExists(req.Email)
	if err != nil {
		return err
	}

	if isExists {
		return errutil.ErrUserIsAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(fmt.Sprintf("error occurred: [%v] while hashing password: [%v]", err, req.Password))
		return err
	}

	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleID:    req.RoleID,
	}

	// create user in the db
	if err := svc.userRepo.CreateUser(user); err != nil {
		logger.Error(fmt.Sprintf("error occurred: [%v] while creating user: [%v]", err, user))
		return err
	}

	return nil
}

func (svc *UserServiceImpl) IsEmailExists(email string) (bool, error) {
	total, err := svc.userRepo.UserCountByEmail(email)
	if err != nil {
		logger.Error(fmt.Sprintf("error occurred: [%v] while counting users by email: [%v]", err, email))
		return false, err
	}

	return total > 0, nil
}

func (svc *UserServiceImpl) ReadUserByEmail(email string) (*models.User, error) {
	user, err := svc.userRepo.ReadUserByEmail(email)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error(fmt.Sprintf("error occurred: [%v] while reading user by email: [%v]", err, email))
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errutil.ErrRecordNotFound
	}

	return user, nil
}

func (svc *UserServiceImpl) StoreInCache(user *types.UserInfo) error {
	userCacheKey := config.Redis().MandatoryPrefix + config.Redis().UserPrefix + strconv.Itoa(user.ID)

	if err := svc.redisSvc.SetStruct(userCacheKey, user, config.Redis().UserCacheTTL); err != nil {
		logger.Error(fmt.Sprintf("could not cache user in redis, id: [%d], err: [%v]", user.ID, err))
	}
	return nil
}

func (svc *UserServiceImpl) ReadUser(id int) (*models.User, error) {
	return svc.userRepo.ReadUserById(id)
}

func (svc *UserServiceImpl) ReadPermissionsByRole(roleID int) ([]*models.Permission, error) {
	return svc.userRepo.ReadPermissionsByRole(roleID)
}
