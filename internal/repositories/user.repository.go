package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Darari17/be-tickitz-full/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetEmail(c context.Context, email string) (*models.User, error) {
	q := "select id, email, password, role, created_at, updated_at from users where email = $1"
	user := models.User{}

	if err := ur.db.QueryRow(c, q, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) InsertUser(c context.Context, user *models.User) error {
	tx, err := ur.db.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	user.ID = uuid.New()
	user.Role = models.RoleUser

	queryInsertUser := "insert into users (id, email, password, role, created_at) values ($1, $2, $3, $4, $5)"
	if _, err := tx.Exec(c, queryInsertUser, user.ID, user.Email, user.Password, user.Role, "now()"); err != nil {
		return err
	}

	queryInsertProfile := "insert into profile (user_id, created_at) values ($1, $2)"
	if _, err := tx.Exec(c, queryInsertProfile, user.ID, "now()"); err != nil {
		return err
	}

	if err = tx.Commit(c); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetProfile(c context.Context, userID uuid.UUID) (*models.Profile, error) {
	sql := `
		SELECT user_id, firstname, lastname, phone_number, avatar, point, created_at, updated_at
		FROM profile
		WHERE user_id = $1
	`
	var profile models.Profile
	err := ur.db.QueryRow(c, sql, userID).Scan(
		&profile.UserID, &profile.FirstName, &profile.LastName, &profile.PhoneNumber,
		&profile.Avatar, &profile.Point, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (ur *UserRepository) UpdateProfile(c context.Context, p *models.Profile) error {
	now := time.Now()

	setParts := []string{}
	args := []any{}
	argID := 1

	if p.FirstName != nil {
		setParts = append(setParts, fmt.Sprintf("firstname = $%d", argID))
		args = append(args, *p.FirstName)
		argID++
	}
	if p.LastName != nil {
		setParts = append(setParts, fmt.Sprintf("lastname = $%d", argID))
		args = append(args, *p.LastName)
		argID++
	}
	if p.PhoneNumber != nil {
		setParts = append(setParts, fmt.Sprintf("phone_number = $%d", argID))
		args = append(args, *p.PhoneNumber)
		argID++
	}

	if len(setParts) == 0 {
		return nil
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, now)
	argID++

	args = append(args, p.UserID)

	sql := fmt.Sprintf(`
		UPDATE profile
		SET %s
		WHERE user_id = $%d
	`, strings.Join(setParts, ", "), argID)

	_, err := ur.db.Exec(c, sql, args...)
	return err
}

func (ur *UserRepository) VerifyPassword(c context.Context, userID uuid.UUID, oldPassword string) (string, error) {
	var hashedPassword string
	sql := `SELECT password FROM users WHERE id = $1`

	err := ur.db.QueryRow(c, sql, userID).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func (ur *UserRepository) UpdatePassword(c context.Context, userID uuid.UUID, hashedPassword string) error {
	sql := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`
	_, err := ur.db.Exec(c, sql, hashedPassword, userID)
	return err
}

func (ur *UserRepository) UpdateAvatar(c context.Context, userID uuid.UUID, avatar string) error {
	sql := `UPDATE profile SET avatar = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := ur.db.Exec(c, sql, avatar, userID)
	return err
}
