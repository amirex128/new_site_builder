package repository

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

// CommentRepo implements ICommentRepository
type CommentRepo struct {
	database *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepo {
	return &CommentRepo{
		database: db,
	}
}

func (r *CommentRepo) GetAllByTicketID(ticketID int64) ([]domain.Comment, error) {
	var comments []domain.Comment
	result := r.database.Where("ticket_id = ?", ticketID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (r *CommentRepo) GetByID(id int64) (domain.Comment, error) {
	var comment domain.Comment
	result := r.database.First(&comment, id)
	if result.Error != nil {
		return comment, result.Error
	}
	return comment, nil
}

func (r *CommentRepo) Create(comment domain.Comment) error {
	result := r.database.Create(&comment)
	return result.Error
}

func (r *CommentRepo) Update(comment domain.Comment) error {
	result := r.database.Save(&comment)
	return result.Error
}

func (r *CommentRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Comment{}, id)
	return result.Error
}

// CustomerCommentRepo implements ICustomerCommentRepository
type CustomerCommentRepo struct {
	database *gorm.DB
}

func NewCustomerCommentRepository(db *gorm.DB) *CustomerCommentRepo {
	return &CustomerCommentRepo{
		database: db,
	}
}

func (r *CustomerCommentRepo) GetAllByCustomerTicketID(customerTicketID int64) ([]domain.CustomerComment, error) {
	var comments []domain.CustomerComment
	result := r.database.Where("customer_ticket_id = ?", customerTicketID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (r *CustomerCommentRepo) GetByID(id int64) (domain.CustomerComment, error) {
	var comment domain.CustomerComment
	result := r.database.First(&comment, id)
	if result.Error != nil {
		return comment, result.Error
	}
	return comment, nil
}

func (r *CustomerCommentRepo) Create(comment domain.CustomerComment) error {
	result := r.database.Create(&comment)
	return result.Error
}

func (r *CustomerCommentRepo) Update(comment domain.CustomerComment) error {
	result := r.database.Save(&comment)
	return result.Error
}

func (r *CustomerCommentRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.CustomerComment{}, id)
	return result.Error
}
