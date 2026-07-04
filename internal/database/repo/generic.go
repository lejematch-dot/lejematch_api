package repo

import "gorm.io/gorm"

type Repository[T any] interface {
	Create(entity *T) error
	FindByID(id uint) (*T, error)
	FindAll() ([]*T, error)
	Update(entity *T) error
	Delete(entity *T) error
}

type GenericRepo[T any] struct {
	db *gorm.DB
}

func NewGenericRepo[T any](db *gorm.DB) *GenericRepo[T] {
	return &GenericRepo[T]{db: db}
}

func (r *GenericRepo[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *GenericRepo[T]) FindByID(id int) (*T, error) {
	var entity T
	err := r.db.First(&entity, uint(id)).Error
	return &entity, err
}

func (r *GenericRepo[T]) FindAll() ([]*T, error) {
	var entities []*T
	err := r.db.Find(&entities).Error
	return entities, err
}

func (r *GenericRepo[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *GenericRepo[T]) Delete(entity *T) error {
	return r.db.Delete(entity).Error
}

func (r *GenericRepo[T]) UpdateFields(id int, fields map[string]interface{}) error {
	return r.db.Model(new(T)).Where("id = ?", id).Updates(fields).Error
}
