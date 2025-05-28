package db

import (
	"github.com/vivasoft-ltd/go-ems/models"
)

func (repo *EventRepositoryImpl) UserCountByEmail(email string) (int, error) {
	var total int64

	if err := repo.client.Model(&models.User{}).Where("email = ?", email).Count(&total).Error; err != nil {
		return 0, err
	}
	return int(total), nil
}

func (repo *EventRepositoryImpl) CreateUser(user *models.User) error {
	return repo.client.Create(user).Error
}

func (repo *EventRepositoryImpl) ReadUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.client.Model(&models.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *EventRepositoryImpl) ReadUserById(id int) (*models.User, error) {
	var user models.User
	if err := repo.client.Model(&models.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *EventRepositoryImpl) ReadPermissionsByRole(roleID int) ([]*models.Permission, error) {
	var permissions []*models.Permission

	if err := repo.client.Model(&models.RolePermission{}).
		Select("permissions.*").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

//
//SELECT permissions.id, permissions.permission, r.name
//FROM role_permissions
//JOIN permissions on role_permissions.permission_id = permissions.id
//JOIN event_management.roles r on r.id = role_permissions.role_id
