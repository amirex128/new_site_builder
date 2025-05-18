package mysql

import (
	"github.com/amirex128/new_site_builder/src/internal/domain"
	"gorm.io/gorm"
)

// TicketMediaRepo implements ITicketMediaRepository
type TicketMediaRepo struct {
	database *gorm.DB
}

func NewTicketMediaRepository(db *gorm.DB) *TicketMediaRepo {
	return &TicketMediaRepo{
		database: db,
	}
}

func (r *TicketMediaRepo) AddMediaToTicket(ticketID int64, mediaIDs []int64) error {
	// Begin a transaction
	tx := r.database.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Add new media associations
	for _, mediaID := range mediaIDs {
		if err := tx.Exec("INSERT INTO ticket_media (ticket_id, media_id) VALUES (?, ?)", ticketID, mediaID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *TicketMediaRepo) RemoveMediaFromTicket(ticketID int64, mediaIDs []int64) error {
	// Begin a transaction
	tx := r.database.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Remove media associations
	for _, mediaID := range mediaIDs {
		if err := tx.Exec("DELETE FROM ticket_media WHERE ticket_id = ? AND media_id = ?", ticketID, mediaID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *TicketMediaRepo) GetMediaByTicketID(ticketID int64) ([]domain.Media, error) {
	var media []domain.Media
	result := r.database.Raw("SELECT m.* FROM media m JOIN ticket_media tm ON m.id = tm.media_id WHERE tm.ticket_id = ?", ticketID).Scan(&media)
	if result.Error != nil {
		return nil, result.Error
	}
	return media, nil
}

// CustomerTicketMediaRepo implements ICustomerTicketMediaRepository
type CustomerTicketMediaRepo struct {
	database *gorm.DB
}

func NewCustomerTicketMediaRepository(db *gorm.DB) *CustomerTicketMediaRepo {
	return &CustomerTicketMediaRepo{
		database: db,
	}
}

func (r *CustomerTicketMediaRepo) AddMediaToCustomerTicket(customerTicketID int64, mediaIDs []int64) error {
	// Begin a transaction
	tx := r.database.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Add new media associations
	for _, mediaID := range mediaIDs {
		if err := tx.Exec("INSERT INTO customer_ticket_media (customer_ticket_id, media_id) VALUES (?, ?)", customerTicketID, mediaID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *CustomerTicketMediaRepo) RemoveMediaFromCustomerTicket(customerTicketID int64, mediaIDs []int64) error {
	// Begin a transaction
	tx := r.database.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Remove media associations
	for _, mediaID := range mediaIDs {
		if err := tx.Exec("DELETE FROM customer_ticket_media WHERE customer_ticket_id = ? AND media_id = ?", customerTicketID, mediaID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *CustomerTicketMediaRepo) GetMediaByCustomerTicketID(customerTicketID int64) ([]domain.Media, error) {
	var media []domain.Media
	result := r.database.Raw("SELECT m.* FROM media m JOIN customer_ticket_media ctm ON m.id = ctm.media_id WHERE ctm.customer_ticket_id = ?", customerTicketID).Scan(&media)
	if result.Error != nil {
		return nil, result.Error
	}
	return media, nil
}
