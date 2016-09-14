package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/kolide/kolide-ose/kolide"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

func (svc service) NewUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	user, err := userFromPayload(p, svc.config.Auth.SaltKeySize, svc.config.Auth.BcryptCost)
	if err != nil {
		return nil, err
	}
	user, err = svc.ds.NewUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc service) User(ctx context.Context, id uint) (*kolide.User, error) {
	return svc.ds.UserByID(id)
}

func (svc service) ChangePassword(ctx context.Context, userID uint, old, new string) error {
	user, err := svc.User(ctx, userID)
	if err != nil {
		return err
	}
	if err := user.ValidatePassword(old); err != nil {
		return fmt.Errorf("current password validation failed: %v", err)
	}
	hashed, salt, err := hashPassword(new, svc.config.Auth.SaltKeySize, svc.config.Auth.BcryptCost)
	if err != nil {
		return err
	}
	user.Salt = salt
	user.Password = hashed
	return svc.saveUser(user)
}

func (svc service) UpdateAdminRole(ctx context.Context, userID uint, isAdmin bool) error {
	user, err := svc.User(ctx, userID)
	if err != nil {
		return err
	}
	user.Admin = isAdmin
	return svc.saveUser(user)
}

func (svc service) UpdateUserStatus(ctx context.Context, userID uint, password string, enabled bool) error {
	user, err := svc.User(ctx, userID)
	if err != nil {
		return err
	}
	user.Enabled = enabled
	return svc.saveUser(user)
}

// saves user in datastore.
// doesn't need to be exposed to the transport
// the service should expose actions for modifying a user instead
func (svc service) saveUser(user *kolide.User) error {
	return svc.ds.SaveUser(user)
}

func userFromPayload(p kolide.UserPayload, keySize, cost int) (*kolide.User, error) {
	hashed, salt, err := hashPassword(*p.Password, keySize, cost)
	if err != nil {
		return nil, err
	}

	return &kolide.User{
		Username: *p.Username,
		Email:    *p.Email,
		Admin:    falseIfNil(p.Admin),
		AdminForcedPasswordReset: falseIfNil(p.AdminForcedPasswordReset),
		Salt:     salt,
		Enabled:  true,
		Password: hashed,
	}, nil
}

func hashPassword(plaintext string, keySize, cost int) ([]byte, string, error) {
	salt, err := generateRandomText(keySize)
	if err != nil {
		return nil, "", err
	}

	withSalt := []byte(fmt.Sprintf("%s%s", plaintext, salt))
	hashed, err := bcrypt.GenerateFromPassword(withSalt, cost)
	if err != nil {
		return nil, "", err
	}

	return hashed, salt, nil

}

// generateRandomText return a string generated by filling in keySize bytes with
// random data and then base64 encoding those bytes
func generateRandomText(keySize int) (string, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// helper to convert a bool pointer false
func falseIfNil(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
