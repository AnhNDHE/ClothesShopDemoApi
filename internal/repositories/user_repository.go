package repositories

import (
	"clothes-shop-api/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, email, password, role string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (email, password, role, is_active)
		VALUES ($1, $2, $3, false)
		RETURNING id, email, password, role, created_at, updated_at, created_by, updated_by, is_active, is_deleted
	`

	var user models.User
	err = r.DB.QueryRow(ctx, query, email, string(hashedPassword), role).Scan(
		&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.CreatedBy, &user.UpdatedBy, &user.IsActive, &user.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, role, created_at, updated_at, created_by, updated_by, is_active, is_deleted
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.CreatedBy, &user.UpdatedBy, &user.IsActive, &user.IsDeleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (r *UserRepository) ActivateUser(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET is_active = true, updated_at = now()
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query, userID)
	return err
}
