package repository

import (
	"math"
	"time"

	"github.com/BimaAdi/fiberGormBoilerplate/core"
	"github.com/BimaAdi/fiberGormBoilerplate/models"
	"gorm.io/gorm"
)

func GetPaginatedUser(tx *gorm.DB, page int, pageSize int, search *string) ([]models.User, int64, int64, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	query := tx.Where("deleted_at IS NULL")
	countQuery := tx.Where("deleted_at IS NULL")

	users := []models.User{}

	if search != nil {
		query = query.Where("email like '%" + *search + "%'")
		countQuery = countQuery.Where("email like '%" + *search + "%'")
	}

	if err := query.
		Order("created_at desc").
		Limit(limit).Offset(offset).
		Find(&users).Error; err != nil {
		return users, 0, 0, err
	}

	var numData int64
	if err := countQuery.Model(&models.User{}).Count(&numData).Error; err != nil {
		return users, 0, 0, err
	}

	numPage := math.Ceil(float64(numData) / float64(pageSize))
	return users, numData, int64(numPage), nil
}

func GetUserById(tx *gorm.DB, id string) (models.User, error) {
	user := models.User{}
	if err := tx.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func CreateUser(tx *gorm.DB, username string, email string, password string, isActive bool, isSuperuser bool, createdAt time.Time, updatedAt *time.Time) (models.User, error) {
	hashedPassword, err := core.HashPassword(password)
	if err != nil {
		return models.User{}, err
	}

	newUser := models.User{
		Email:       email,
		Username:    username,
		Password:    hashedPassword,
		IsActive:    isActive,
		IsSuperuser: isSuperuser,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		DeletedAt:   nil,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		return newUser, err
	}
	return newUser, nil
}

func UpdateUser(tx *gorm.DB, updatedUser models.User, email string, username string, password *string, isActive bool, isSuperUser bool) (models.User, error) {
	// Hashed Password
	if password != nil {
		rawPassword := password
		hashedPassword, err := core.HashPassword(*rawPassword)
		if err != nil {
			return models.User{}, err
		}
		updatedUser.Password = hashedPassword
	}

	// Update data
	updatedUser.Email = email
	updatedUser.Username = username
	updatedUser.IsActive = isActive
	updatedUser.IsSuperuser = isSuperUser
	now := time.Now()
	updatedUser.UpdatedAt = &now
	if err := tx.Save(&updatedUser).Error; err != nil {
		return updatedUser, err
	}
	return updatedUser, nil
}

func DeleteUser(tx *gorm.DB, user models.User) (models.User, error) {
	now := time.Now()
	user.DeletedAt = &now
	if err := tx.Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func GetUserByUsername(tx *gorm.DB, username string) (models.User, error) {
	user := models.User{}
	if err := tx.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func GetUserByUsernameOrEmail(tx *gorm.DB, usernameOrEmail string) (models.User, error) {
	user := models.User{}
	if err := tx.Where("(username = ? OR email = ? ) AND deleted_at IS NULL", usernameOrEmail, usernameOrEmail).
		First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
