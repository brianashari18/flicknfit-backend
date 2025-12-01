package services

import (
	"context"
	"errors"
	"flicknfit_backend/config"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// FirebaseService handles Firebase authentication operations
type FirebaseService struct {
	authClient *auth.Client
}

// NewFirebaseService creates a new Firebase service instance
func NewFirebaseService(cfg *config.Config) (*FirebaseService, error) {
	ctx := context.Background()

	// Initialize Firebase app with service account credentials
	opt := option.WithCredentialsFile(cfg.FirebasePrivateKeyPath)
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: cfg.FirebaseProjectID,
	}, opt)
	if err != nil {
		return nil, errors.New("failed to initialize Firebase app: " + err.Error())
	}

	// Get Firebase Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, errors.New("failed to get Firebase Auth client: " + err.Error())
	}

	return &FirebaseService{
		authClient: authClient,
	}, nil
}

// VerifyIDToken verifies a Firebase ID token and returns the decoded token
func (s *FirebaseService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := s.authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, errors.New("invalid Firebase ID token: " + err.Error())
	}

	return token, nil
}

// GetUserByUID retrieves user information from Firebase by UID
func (s *FirebaseService) GetUserByUID(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := s.authClient.GetUser(ctx, uid)
	if err != nil {
		return nil, errors.New("failed to get user from Firebase: " + err.Error())
	}

	return user, nil
}
