package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"tender/internal/models"
	"time"
)

type TenderRepository struct {
	Db *sql.DB
}

// CreateTender inserts a new tender into the database.
func (r *TenderRepository) CreateTender(ctx context.Context, tender models.Tender) (int, error) {
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")

	result, err := r.Db.ExecContext(ctx, `
        INSERT INTO tenders (
            type, tender_number, user_id, company_id, organization,
            total, commission, completed_date, date, status
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tender.Type, tender.TenderNumber, tender.UserID, tender.CompanyID, tender.Organization,
		tender.Total, tender.Commission, tender.CompletedDate, formattedTime, tender.Status,
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
        SELECT tenders.id, type, tender_number, user_id, company_id, organization,
               total, commission, completed_date, date,  u.name, c.name
        FROM tenders
        JOIN tender.users u ON u.id = tenders.user_id
		JOIN tender.companies c ON c.id = tenders.company_id
		 WHERE tenders.id = ?
		ORDER BY tenders.date DESC`, id).Scan(
		&tender.ID, &tender.Type, &tender.TenderNumber, &tender.UserID, &tender.CompanyID, &tender.Organization,
		&tender.Total, &tender.Commission, &tender.CompletedDate, &tender.Date, &tender.UserName, &tender.CompanyName,
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
        SELECT tenders.id, type, tender_number, user_id, company_id, organization,
               total, commission, completed_date, date, status, u.name, c.name
        FROM tenders
        JOIN tender.users u ON u.id = tenders.user_id
		JOIN tender.companies c ON c.id = tenders.company_id
		ORDER BY tenders.date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(
			&tender.ID, &tender.Type, &tender.TenderNumber, &tender.UserID, &tender.CompanyID, &tender.Organization,
			&tender.Total, &tender.Commission, &tender.CompletedDate, &tender.Date, &tender.Status, &tender.UserName, &tender.CompanyName,
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

// GetTendersDebtByCompanyId retrieves debt sum by company id from the database.
func (r *TenderRepository) GetTotalNetByCompany(ctx context.Context) ([]struct {
	CompanyID int
	TotalNet  float64
}, error) {
	rows, err := r.Db.QueryContext(ctx, `
		SELECT company_id, SUM(total - commission) AS total_net
		FROM tenders
		GROUP BY company_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		CompanyID int
		TotalNet  float64
	}
	for rows.Next() {
		var result struct {
			CompanyID int
			TotalNet  float64
		}
		err := rows.Scan(&result.CompanyID, &result.TotalNet)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *TenderRepository) GetTendersByUserID(ctx context.Context, userID int) ([]models.Tender, error) {
	rows, err := r.Db.QueryContext(ctx, `
        SELECT tenders.id, type, tender_number, user_id, company_id, organization,
               total, commission, completed_date, date, status, u.name, c.name
        FROM tenders
        JOIN tender.users u ON u.id = tenders.user_id
        JOIN tender.companies c ON c.id = tenders.company_id
        WHERE tenders.user_id = ?
        ORDER BY tenders.date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(
			&tender.ID, &tender.Type, &tender.TenderNumber, &tender.UserID, &tender.CompanyID, &tender.Organization,
			&tender.Total, &tender.Commission, &tender.CompletedDate, &tender.Date, &tender.Status, &tender.UserName, &tender.CompanyName,
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

func (r *TenderRepository) GetTendersByCompanyID(ctx context.Context, companyID int) ([]models.Tender, error) {
	rows, err := r.Db.QueryContext(ctx, `
        SELECT tenders.id, type, tender_number, user_id, company_id, organization,
               total, commission, completed_date, date, status, u.name, c.name
        FROM tenders
        JOIN tender.users u ON u.id = tenders.user_id
        JOIN tender.companies c ON c.id = tenders.company_id
        WHERE tenders.company_id = ?
        ORDER BY tenders.date DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(
			&tender.ID, &tender.Type, &tender.TenderNumber, &tender.UserID, &tender.CompanyID, &tender.Organization,
			&tender.Total, &tender.Commission, &tender.CompletedDate, &tender.Date, &tender.Status, &tender.UserName, &tender.CompanyName,
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

func (r *TenderRepository) GetAllTendersSum(ctx context.Context) (*models.TenderDebt, error) {
	queryGOIK := `
        SELECT COALESCE(SUM(total-commission), 0) AS total_sum
        FROM tenders
        WHERE status = 3 AND status = 2
        AND type = 'ГОИК';
    `

	// Execute the query for Zakup
	var totalGOIK float64
	err := r.Db.QueryRowContext(ctx, queryGOIK).Scan(&totalGOIK)
	if err != nil {
		return nil, err
	}

	queryGOPP := `
        SELECT COALESCE(SUM(total-commission), 0) AS total_sum
        FROM tenders
        WHERE status = 3 AND status = 2
        AND type = 'ГОПП';
    `

	// Execute the query for Zakup
	var totalGOPP float64
	err = r.Db.QueryRowContext(ctx, queryGOPP).Scan(&totalGOPP)
	if err != nil {
		return nil, err
	}
	// Return the result in a struct
	return &models.TenderDebt{
		GOIK: totalGOIK,
		GOPP: totalGOPP,
	}, nil
}

func (r *TenderRepository) GetTenderCountsByUserID(ctx context.Context, userID int) (*models.TenderCount, error) {
	// Query to count total tenders by user ID
	queryTotal := `
        SELECT COUNT(*) 
        FROM tenders 
        WHERE user_id = ?;
    `

	// Query to count tenders by user ID and status
	queryStatus := `
        SELECT COUNT(*) 
        FROM tenders 
        WHERE user_id = ? AND status = ?;
    `

	counts := &models.TenderCount{}

	// Execute total tenders query
	err := r.Db.QueryRowContext(ctx, queryTotal, userID).Scan(&counts.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to count total tenders: %w", err)
	}

	// Execute queries for each status
	statusCounts := []*int{&counts.Status0, &counts.Status1, &counts.Status2, &counts.Status3}
	for i := 0; i < 4; i++ {
		err = r.Db.QueryRowContext(ctx, queryStatus, userID, i).Scan(statusCounts[i])
		if err != nil {
			return nil, fmt.Errorf("failed to count status %d tenders: %w", i, err)
		}
	}

	return counts, nil
}

func (r *TenderRepository) GetAllTendersByDateRange(ctx context.Context, startDate, endDate string, userId int) ([]models.Tender, error) {
	var query string
	var rows *sql.Rows
	var err error

	if userId == 1 {
		query = `
			SELECT t.id, t.type, t.tender_number, t.user_id, t.company_id, t.organization, 
			       t.total, t.commission, t.completed_date, t.date, t.status, 
			       u.name as username, c.name as companyname
			FROM tenders t 
			JOIN users u ON u.id = t.user_id
			JOIN companies c ON c.id = t.company_id
			WHERE t.completed_date BETWEEN ? AND ?
			ORDER BY t.completed_date DESC
		`
		rows, err = r.Db.QueryContext(ctx, query, startDate, endDate)
	} else {
		query = `
			SELECT t.id, t.type, t.tender_number, t.user_id, t.company_id, t.organization, 
			       t.total, t.commission, t.completed_date, t.date, t.status, 
			       u.name as username, c.name as companyname
			FROM tenders t 
			JOIN users u ON u.id = t.user_id
			JOIN companies c ON c.id = t.company_id
			WHERE t.user_id = ? AND t.completed_date BETWEEN ? AND ?
			ORDER BY t.completed_date DESC
		`
		rows, err = r.Db.QueryContext(ctx, query, userId, startDate, endDate)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query tenders: %w", err)
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Type, &tender.TenderNumber, &tender.UserID, &tender.CompanyID,
			&tender.Organization, &tender.Total, &tender.Commission, &tender.CompletedDate, &tender.Date,
			&tender.Status, &tender.UserName, &tender.CompanyName); err != nil {
			return nil, fmt.Errorf("failed to scan tender: %w", err)
		}
		tenders = append(tenders, tender)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenders, nil
}

func (r *TenderRepository) GetAllTendersByDateRangeCompany(ctx context.Context, startDate, endDate string, userId, companyId int) ([]models.Tender, error) {
	var query string
	var rows *sql.Rows
	var err error

	if userId == 1 {
		query = `
			SELECT t.id, t.type, t.tender_number, t.user_id, t.company_id, t.organization, 
			       t.total, t.commission, t.completed_date, t.date, t.status, 
			       u.name as username, c.name as companyname
			FROM tenders t 
			JOIN users u ON u.id = t.user_id
			JOIN companies c ON c.id = t.company_id
			WHERE t.company_id = ? AND t.completed_date BETWEEN ? AND ?
			ORDER BY t.completed_date DESC
		`
		rows, err = r.Db.QueryContext(ctx, query, companyId, startDate, endDate)
	} else {
		query = `
			SELECT t.id, t.type, t.tender_number, t.user_id, t.company_id, t.organization, 
			       t.total, t.commission, t.completed_date, t.date, t.status, 
			       u.name as username, c.name as companyname
			FROM tenders t 
			JOIN users u ON u.id = t.user_id
			JOIN companies c ON c.id = t.company_id
			WHERE t.user_id = ? AND t.company_id = ? AND t.completed_date BETWEEN ? AND ?
			ORDER BY t.completed_date DESC
		`
		rows, err = r.Db.QueryContext(ctx, query, userId, companyId, startDate, endDate)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query tenders: %w", err)
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Type, &tender.TenderNumber, &tender.UserID, &tender.CompanyID,
			&tender.Organization, &tender.Total, &tender.Commission, &tender.CompletedDate, &tender.Date,
			&tender.Status, &tender.UserName, &tender.CompanyName); err != nil {
			return nil, fmt.Errorf("failed to scan tender: %w", err)
		}
		tenders = append(tenders, tender)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenders, nil
}
