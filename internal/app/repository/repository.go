package repository

import (
	"database/sql"
	"wedding_website/internal/app/models"
)

type RSVPRepository struct {
	db *sql.DB
}

func NewRSVPRepository(db *sql.DB) *RSVPRepository {
	return &RSVPRepository{db: db}
}

func (r *RSVPRepository) CreateRSVP(rsvp *models.RSVP) error {
	query := `
        INSERT INTO rsvp_responses (name, attendance, companion)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `

	return r.db.QueryRow(
		query,
		rsvp.Name,
		rsvp.Attendance,
		rsvp.Companion,
	).Scan(&rsvp.ID, &rsvp.CreatedAt)
}

func (r *RSVPRepository) GetAllRSVPs() ([]models.RSVP, error) {
	query := `SELECT id, name, attendance, companion, created_at FROM rsvp_responses ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []models.RSVP
	for rows.Next() {
		var rsvp models.RSVP
		err := rows.Scan(
			&rsvp.ID,
			&rsvp.Name,
			&rsvp.Attendance,
			&rsvp.Companion,
			&rsvp.CreatedAt,
		)
		if err != nil {
		}
		responses = append(responses, rsvp)
	}

	return responses, nil
}
