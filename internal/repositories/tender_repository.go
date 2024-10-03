package repositories

import (
	"context"
	"database/sql"
	"tender/internal/models"
)

type TenderRepository struct {
	Db *sql.DB
}

// CreateTender inserts a new tender into the database.
func (r *TenderRepository) CreateTender(ctx context.Context, tender models.Tender) (int, error) {
	result, err := r.Db.ExecContext(ctx, `
        INSERT INTO tenders (
            type, tender_number, user_id, company_id, organization,
            total, commission, completed_date, date, status
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tender.Type, tender.TenderNumber, tender.UserID, tender.CompanyID, tender.Organization,
		tender.Total, tender.Commission, tender.CompletedDate, tender.Date, tender.Status,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// DeleteTender removes a tender from the database by ID.
func (r *TenderRepository) DeleteTender(ctx context.Context, id int) error {
	result, err := r.Db.ExecContext(ctx, "DELETE FROM tenders WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrTenderNotFound
	}

	return nil
}

// UpdateTender updates an existing tender in the database.
func (r *TenderRepository) UpdateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	query := "UPDATE tenders SET"
	params := []interface{}{}

	if tender.Type != "" {
		query += " type = ?,"
		params = append(params, tender.Type)
	}
	if tender.TenderNumber != nil {
		query += " tender_number = ?,"
		params = append(params, tender.TenderNumber)
	}
	if tender.UserID != 0 {
		query += " user_id = ?,"
		params = append(params, tender.UserID)
	}
	if tender.CompanyID != 0 {
		query += " company_id = ?,"
		params = append(params, tender.CompanyID)
	}
	if tender.Organization != "" {
		query += " organization = ?,"
		params = append(params, tender.Organization)
	}
	if tender.Total != 0 {
		query += " total = ?,"
		params = append(params, tender.Total)
	}
	if tender.Commission != 0 {
		query += " commission = ?,"
		params = append(params, tender.Commission)
	}
	if tender.CompletedDate == nil {
		query += " completed_date = ?,"
		params = append(params, tender.CompletedDate)
	}
	if !tender.Date.IsZero() {
		query += " date = ?,"
		params = append(params, tender.Date)
	}
	if tender.Status != 0 {
		query += " status = ?,"
		params = append(params, tender.Status)
	}

	if len(params) == 0 {
		// No fields to update
		return r.GetTenderByID(ctx, tender.ID)
	}

	// Trim the last comma from the query
	query = query[:len(query)-1]
	query += " WHERE id = ?"
	params = append(params, tender.ID)

	_, err := r.Db.ExecContext(ctx, query, params...)
	if err != nil {
		return models.Tender{}, err
	}

	// Retrieve the updated tender data
	return r.GetTenderByID(ctx, tender.ID)
}

// GetTenderByID retrieves a tender by ID from the database.
func (r *TenderRepository) GetTenderByID(ctx context.Context, id int) (models.Tender, error) {
	var tender models.Tender
	err := r.Db.QueryRowContext(ctx, `
        SELECT id, type, tender_number, user_id, company_id, organization,
               total, commission, completed_date, date, status
        FROM tenders WHERE id = ?`, id).Scan(
		&tender.ID, &tender.Type, &tender.TenderNumber, &tender.UserID, &tender.CompanyID, &tender.Organization,
		&tender.Total, &tender.Commission, &tender.CompletedDate, &tender.Date, &tender.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return tender, models.ErrTenderNotFound
		}
		return tender, err
	}

	return tender, nil
}

// GetAllTenders retrieves all tenders from the database.
func (r *TenderRepository) GetAllTenders(ctx context.Context) ([]models.Tender, error) {
	rows, err := r.Db.QueryContext(ctx, `
        SELECT id, type, tender_number, user_id, company_id, organization,
               total, commission, completed_date, date, status
        FROM tenders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(
			&tender.ID, &tender.Type, &tender.TenderNumber, &tender.UserID, &tender.CompanyID, &tender.Organization,
			&tender.Total, &tender.Commission, &tender.CompletedDate, &tender.Date, &tender.Status,
		)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenders, nil
}
