// Package store provides user-related data access.
package store

import (
	"context"
	"time"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/pkg/database"
)

// IUserStore defines the interface for user data operations.
type IUserStore interface {
	// User operations
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserWithPreloads(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error
	UpdateUserModel(ctx context.Context, user *model.User) error
	CountUsersWithEmail(ctx context.Context, email string, excludeID *uint) (int64, error)

	// Email verification code operations
	CreateVerificationCode(ctx context.Context, code *model.EmailVerificationCode) error
	GetVerificationCode(ctx context.Context, email string, codeType model.CodeType) (*model.EmailVerificationCode, error)
	GetLatestVerificationCode(ctx context.Context, email string, codeType model.CodeType) (*model.EmailVerificationCode, error)
	CountRecentCodes(ctx context.Context, email string, since time.Time) (int64, error)
	UpdateVerificationCode(ctx context.Context, code *model.EmailVerificationCode) error
	DeleteVerificationCode(ctx context.Context, code *model.EmailVerificationCode) error

	// OAuth identity operations
	ListOAuthIdentities(ctx context.Context, userID uint) ([]model.OAuthIdentity, error)
	GetOAuthIdentityByUserIDAndProvider(ctx context.Context, userID uint, provider string) (*model.OAuthIdentity, error)
	GetOAuthIdentityByProviderAndProviderUserID(ctx context.Context, provider string, providerUserID string) (*model.OAuthIdentity, error)
	CreateOAuthIdentity(ctx context.Context, identity *model.OAuthIdentity) error
	UpdateOAuthIdentityModel(ctx context.Context, identity *model.OAuthIdentity) error
	DeleteOAuthIdentity(ctx context.Context, userID uint, provider string) error
	UpdateOAuthIdentity(ctx context.Context, userID uint, provider string, updates map[string]interface{}) error
	UpdateAllOAuthIdentities(ctx context.Context, userID uint, updates map[string]interface{}) error
	CountOAuthIdentities(ctx context.Context, userID uint) (int64, error)

	// Session operations
	CreateSession(ctx context.Context, session *model.Session) error
	RevokeSession(ctx context.Context, tokenHash string, userID uint) error
	GetValidSession(ctx context.Context, tokenHash string) (*model.Session, error)

	// Login history operations
	CreateLoginHistory(ctx context.Context, history *model.LoginHistory) error

	// OAuth state operations (for CSRF and duplicate prevention)
	CreateOAuthState(ctx context.Context, state *model.OAuthState) error
	GetOAuthState(ctx context.Context, stateToken string) (*model.OAuthState, error)
	MarkOAuthStateProcessed(ctx context.Context, stateToken string) error
	DeleteExpiredOAuthStates(ctx context.Context) error
}

// userStore implements IUserStore interface.
type userStore struct {
	db database.Database
}

// NewUserStore creates a new user store.
func NewUserStore(db database.Database) IUserStore {
	return &userStore{db: db}
}

// User operations

func (s *userStore) CreateUser(ctx context.Context, user *model.User) error {
	return s.db.Conn(ctx).Create(user).Error
}

func (s *userStore) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := s.db.Conn(ctx).First(&user, id).Error
	return &user, err
}

func (s *userStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := s.db.Conn(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (s *userStore) GetUserWithPreloads(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := s.db.Conn(ctx).Preload("OAuthIdentities").First(&user, id).Error
	return &user, err
}

func (s *userStore) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.db.Conn(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Updates(updates).
		Error
}

func (s *userStore) UpdateUserModel(ctx context.Context, user *model.User) error {
	return s.db.Conn(ctx).Save(user).Error
}

func (s *userStore) CountUsersWithEmail(ctx context.Context, email string, excludeID *uint) (int64, error) {
	query := s.db.Conn(ctx).Model(&model.User{}).Where("email = ?", email)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Email verification code operations

func (s *userStore) CreateVerificationCode(ctx context.Context, code *model.EmailVerificationCode) error {
	return s.db.Conn(ctx).Create(code).Error
}

func (s *userStore) GetVerificationCode(ctx context.Context, email string, codeType model.CodeType) (*model.EmailVerificationCode, error) {
	var code model.EmailVerificationCode
	err := s.db.Conn(ctx).
		Where("email = ? AND type = ?", email, codeType).
		Order("created_at DESC").
		First(&code).
		Error
	return &code, err
}

func (s *userStore) GetLatestVerificationCode(ctx context.Context, email string, codeType model.CodeType) (*model.EmailVerificationCode, error) {
	var code model.EmailVerificationCode
	err := s.db.Conn(ctx).
		Where("email = ? AND type = ? AND used = false AND expires_at > ?", email, codeType, time.Now()).
		Order("created_at DESC").
		First(&code).
		Error
	return &code, err
}

func (s *userStore) CountRecentCodes(ctx context.Context, email string, since time.Time) (int64, error) {
	var count int64
	err := s.db.Conn(ctx).Model(&model.EmailVerificationCode{}).
		Where("email = ? AND created_at > ?", email, since).
		Count(&count).
		Error
	return count, err
}

func (s *userStore) UpdateVerificationCode(ctx context.Context, code *model.EmailVerificationCode) error {
	return s.db.Conn(ctx).Save(code).Error
}

func (s *userStore) DeleteVerificationCode(ctx context.Context, code *model.EmailVerificationCode) error {
	return s.db.Conn(ctx).Delete(code).Error
}

// OAuth identity operations

func (s *userStore) ListOAuthIdentities(ctx context.Context, userID uint) ([]model.OAuthIdentity, error) {
	var identities []model.OAuthIdentity
	err := s.db.Conn(ctx).Where("user_id = ?", userID).Find(&identities).Error
	return identities, err
}

func (s *userStore) GetOAuthIdentityByUserIDAndProvider(ctx context.Context, userID uint, provider string) (*model.OAuthIdentity, error) {
	var identity model.OAuthIdentity
	err := s.db.Conn(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		First(&identity).
		Error
	return &identity, err
}

func (s *userStore) GetOAuthIdentityByProviderAndProviderUserID(ctx context.Context, provider string, providerUserID string) (*model.OAuthIdentity, error) {
	var identity model.OAuthIdentity
	err := s.db.Conn(ctx).
		Where("provider = ? AND provider_user_id = ?", provider, providerUserID).
		First(&identity).
		Error
	return &identity, err
}

func (s *userStore) CreateOAuthIdentity(ctx context.Context, identity *model.OAuthIdentity) error {
	return s.db.Conn(ctx).Create(identity).Error
}

func (s *userStore) UpdateOAuthIdentityModel(ctx context.Context, identity *model.OAuthIdentity) error {
	return s.db.Conn(ctx).Save(identity).Error
}

func (s *userStore) DeleteOAuthIdentity(ctx context.Context, userID uint, provider string) error {
	return s.db.Conn(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		Delete(&model.OAuthIdentity{}).
		Error
}

func (s *userStore) UpdateOAuthIdentity(ctx context.Context, userID uint, provider string, updates map[string]interface{}) error {
	return s.db.Conn(ctx).Model(&model.OAuthIdentity{}).
		Where("user_id = ? AND provider = ?", userID, provider).
		Updates(updates).
		Error
}

func (s *userStore) UpdateAllOAuthIdentities(ctx context.Context, userID uint, updates map[string]interface{}) error {
	return s.db.Conn(ctx).Model(&model.OAuthIdentity{}).
		Where("user_id = ?", userID).
		Updates(updates).
		Error
}

func (s *userStore) CountOAuthIdentities(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := s.db.Conn(ctx).Model(&model.OAuthIdentity{}).
		Where("user_id = ?", userID).
		Count(&count).
		Error
	return count, err
}

// Session operations

func (s *userStore) CreateSession(ctx context.Context, session *model.Session) error {
	return s.db.Conn(ctx).Create(session).Error
}

func (s *userStore) RevokeSession(ctx context.Context, tokenHash string, userID uint) error {
	return s.db.Conn(ctx).Model(&model.Session{}).
		Where("token_hash = ? AND user_id = ?", tokenHash, userID).
		Update("revoked", true).
		Error
}

func (s *userStore) GetValidSession(ctx context.Context, tokenHash string) (*model.Session, error) {
	var session model.Session
	err := s.db.Conn(ctx).
		Where("token_hash = ? AND revoked = false AND expires_at > ?", tokenHash, time.Now()).
		First(&session).
		Error
	return &session, err
}

// Login history operations

func (s *userStore) CreateLoginHistory(ctx context.Context, history *model.LoginHistory) error {
	return s.db.Conn(ctx).Create(history).Error
}

// OAuth state operations

func (s *userStore) CreateOAuthState(ctx context.Context, state *model.OAuthState) error {
	return s.db.Conn(ctx).Create(state).Error
}

func (s *userStore) GetOAuthState(ctx context.Context, stateToken string) (*model.OAuthState, error) {
	var state model.OAuthState
	err := s.db.Conn(ctx).Where("state = ?", stateToken).First(&state).Error
	return &state, err
}

func (s *userStore) MarkOAuthStateProcessed(ctx context.Context, stateToken string) error {
	return s.db.Conn(ctx).Model(&model.OAuthState{}).
		Where("state = ?", stateToken).
		Update("processed", true).
		Error
}

func (s *userStore) DeleteExpiredOAuthStates(ctx context.Context) error {
	return s.db.Conn(ctx).Where("expires_at < ?", time.Now()).Delete(&model.OAuthState{}).Error
}
