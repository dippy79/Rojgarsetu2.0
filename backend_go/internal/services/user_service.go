package services

import (
	"context"
	"errors"
	"fmt"

	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
	"golang.org/x/crypto/bcrypt"

	"github.com/rojgarsetu/backend/internal/db"
)

type UserService struct {
	db *db.PostgresDB
}

func NewUserService(d *db.PostgresDB) *UserService {
	return &UserService{db: d}
}

var ErrCompanyNameExists = errors.New("a company with this name already exists")

func (s *UserService) CreateUser(ctx context.Context, req db.RegisterRequest) (*db.User, error) {
	uid := uuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 1. Start Transaction (P0/P1 Atomicity Fix)
	tx, err := s.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer rollback (it does nothing if already committed)
	defer tx.Rollback()

	qtx := s.db.Queries.WithTx(tx)

	// 2. Create Base User
	user, err := qtx.CreateUser(ctx, db.CreateUserParams{
		ID:           uid,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
		Phone:        req.Phone,
	})
	if err != nil {
		return nil, err
	}

	// 3. Create Role-Specific Profile
	switch req.Role {
	case "candidate":
		_, err = qtx.CreateCandidate(ctx, db.CreateCandidateParams{
			UserID:            uid,
			Phone:             sql.NullString{},
			ResumeUrl:         sql.NullString{},
			ResumeParsed:      pqtype.NullRawMessage{},
			Skills:            []string{},
			ExperienceYears:   sql.NullInt32{},
			CurrentCompany:    sql.NullString{},
			CurrentPosition:   sql.NullString{},
			Location:          sql.NullString{},
			LinkedinUrl:       sql.NullString{},
			GithubUrl:         sql.NullString{},
			PortfolioUrl:      sql.NullString{},
			Bio:               sql.NullString{},
			IsOpenToWork:      sql.NullBool{},
			ExpectedSalary:    sql.NullString{},
			PreferredJobType:  []string{},
			PreferredLocation: []string{},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create candidate profile: %w", err)
		}
	case "company":
		_, err = qtx.CreateCompany(ctx, db.CreateCompanyParams{
			UserID:       uid,
			Name:         req.Name,
			Industry:     sql.NullString{},
			CompanySize:  sql.NullString{},
			Website:      sql.NullString{},
			LogoUrl:      sql.NullString{},
			Description:  sql.NullString{},
			Headquarters: sql.NullString{},
			FoundedYear:  sql.NullInt32{},
		})
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
				return nil, ErrCompanyNameExists
			}
			return nil, fmt.Errorf("failed to create company profile: %w", err)
		}
	}

	// 4. Commit Transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*db.User, error) {
	user, err := s.db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	result, err := s.db.Queries.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update last login: %w", err)
	}
	return &result, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*db.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := s.db.Queries.GetUserByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	result, err := s.db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *UserService) ListUsers(ctx context.Context, page, limit int, role string) ([]db.ListUsersRow, int, error) {
	page = clampPage(page)
	limit = clampLimit(limit)
	offset := (page - 1) * limit
	rows, err := s.db.Queries.ListUsers(ctx, db.ListUsersParams{
		Role:   role,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return rows, len(rows), nil
}

func (s *UserService) EmailExists(ctx context.Context, email string) (bool, error) {
	_, err := s.db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *UserService) UpdateLastLogin(ctx context.Context, id uuid.UUID) (*db.User, error) {
	result, err := s.db.Queries.UpdateLastLogin(ctx, id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
