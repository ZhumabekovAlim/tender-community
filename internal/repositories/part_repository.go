package repositories

import (
	"context"
	"database/sql"
	"tender/internal/models"
)

type PartRepository struct {
	db *sql.DB
}

func NewPartRepository(db *sql.DB) *PartRepository {
	return &PartRepository{db: db}
}

func (r *PartRepository) GetAllParts(ctx context.Context) ([]models.Part, error) {
	// Query database to fetch all parts
	// Example query:
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, brand, price, quantity FROM parts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []models.Part
	for rows.Next() {
		var part models.Part
		if err := rows.Scan(&part.ID, &part.Name, &part.Brand, &part.Price, &part.Quantity); err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	return parts, nil
}

func (r *PartRepository) AddPart(ctx context.Context, part models.Part) error {
	// Insert new part into database
	// Example query:
	_, err := r.db.ExecContext(ctx, "INSERT INTO parts(name, brand, price, quantity) VALUES (?, ?, ?, ?)",
		part.Name, part.Brand, part.Price, part.Quantity)
	if err != nil {
		return err
	}
	return nil
}
