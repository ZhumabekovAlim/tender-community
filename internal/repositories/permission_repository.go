package repositories

import (
	"context"
	"database/sql"
	"tender/internal/models"
)

type PermissionRepository struct {
	Db *sql.DB
}

// AddPermission inserts a new permission into the database.
func (r *PermissionRepository) AddPermission(ctx context.Context, permission models.Permission) error {
	_, err := r.Db.ExecContext(ctx, "INSERT INTO permissions (user_id, company_id, status) VALUES (?, ?, 1)", permission.UserID, permission.CompanyID)
	return err
}

// DeletePermission removes a permission from the database by ID.
func (r *PermissionRepository) DeletePermission(ctx context.Context, id int) error {
	result, err := r.Db.ExecContext(ctx, "DELETE FROM permissions WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrPermissionNotFound
	}

	return nil
}

// UpdatePermission updates an existing permission in the database.
func (r *PermissionRepository) UpdatePermission(ctx context.Context, permission models.Permission) error {
	result, err := r.Db.ExecContext(ctx, "UPDATE permissions SET user_id = ?, company_id = ?, status = ? WHERE id = ?", permission.UserID, permission.CompanyID, permission.Status, permission.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrPermissionNotFound
	}

	return nil
}

// GetPermission retrieves a permission by ID from the database.
func (r *PermissionRepository) GetPermissionsByUserID(ctx context.Context, userID int) ([]models.Permission, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT id, user_id, company_id, status FROM permissions WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission

	for rows.Next() {
		var permission models.Permission
		err := rows.Scan(&permission.ID, &permission.UserID, &permission.CompanyID)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
