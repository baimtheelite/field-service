package repositories

import (
	"context"
	"errors"
	errorWrap "field-service/common/error"
	errConstants "field-service/constants/error"
	"field-service/domain/dto"
	"field-service/domain/models"
	"fmt"

	"gorm.io/gorm"
)

type FieldRepository struct {
	db *gorm.DB
}

type IFieldRepository interface {
	FindAllWithPagination(context.Context, *dto.FieldRequestParam) ([]models.Field, int64, error)
	FindAllWithoutPagination(context.Context) ([]models.Field, error)
	FindByUUID(context.Context, string) (*models.Field, error)
	Create(context.Context, *models.Field) (*models.Field, error)
	Update(context.Context, *models.Field) (*models.Field, error)
	Delete(context.Context, string) error
}

func NewFieldRepository(db *gorm.DB) IFieldRepository {
	return &FieldRepository{
		db: db,
	}
}

func (f *FieldRepository) FindAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) ([]models.Field, int64, error) {
	var (
		fields []models.Field
		sort   string
		total  int64
	)
	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", param.SortColumn, param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := f.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&fields).
		Error
	if err != nil {
		return nil, 0, errorWrap.WrapError(errConstants.ErrSQLError)
	}

	err = f.db.WithContext(ctx).
		Model(&fields).
		Count(&total).
		Error
	if err != nil {
		return nil, 0, errorWrap.WrapError(errConstants.ErrSQLError)
	}
	return fields, total, nil
}

func (f *FieldRepository) FindAllWithoutPagination(ctx context.Context) ([]models.Field, error) {
	var fields []models.Field
	err := f.db.WithContext(ctx).
		Find(&fields).
		Error
	if err != nil {
		return nil, errorWrap.WrapError(errConstants.ErrSQLError)
	}
	return fields, nil
}

func (f *FieldRepository) FindByUUID(ctx context.Context, uuid string) (*models.Field, error) {
	var field models.Field
	err := f.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&field).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorWrap.WrapError(errConstants.ErrFieldNotFound)
		}
		return nil, errorWrap.WrapError(errConstants.ErrSQLError)
	}
	return &field, nil
}
