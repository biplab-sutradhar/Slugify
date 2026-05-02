package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
)

// ---------- fakes ----------

type fakeUserRepo struct {
	byEmail map[string]models.User
	byID    map[string]models.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		byEmail: map[string]models.User{},
		byID:    map[string]models.User{},
	}
}

func (r *fakeUserRepo) CreateUser(_ context.Context, u models.User) error {
	if _, exists := r.byEmail[u.Email]; exists {
		return errors.New("dup")
	}
	r.byEmail[u.Email] = u
	r.byID[u.ID] = u
	return nil
}

func (r *fakeUserRepo) GetUserByEmail(_ context.Context, email string) (models.User, error) {
	u, ok := r.byEmail[email]
	if !ok {
		return models.User{}, errors.New("not found")
	}
	return u, nil
}

func (r *fakeUserRepo) GetUserByID(_ context.Context, id string) (models.User, error) {
	u, ok := r.byID[id]
	if !ok {
		return models.User{}, errors.New("not found")
	}
	return u, nil
}

type fakeAPIKeyRepo struct {
	byID  map[string]models.APIKey
	byKey map[string]models.APIKey
}

func newFakeAPIKeyRepo() *fakeAPIKeyRepo {
	return &fakeAPIKeyRepo{
		byID:  map[string]models.APIKey{},
		byKey: map[string]models.APIKey{},
	}
}

func (r *fakeAPIKeyRepo) CreateAPIKey(_ context.Context, k models.APIKey) error {
	r.byID[k.ID] = k
	r.byKey[k.Key] = k
	return nil
}
func (r *fakeAPIKeyRepo) GetAPIKeyByKey(_ context.Context, k string) (models.APIKey, error) {
	v, ok := r.byKey[k]
	if !ok {
		return models.APIKey{}, errors.New("not found")
	}
	return v, nil
}
func (r *fakeAPIKeyRepo) GetAPIKeys(_ context.Context) ([]models.APIKey, error) {
	out := make([]models.APIKey, 0, len(r.byID))
	for _, v := range r.byID {
		out = append(out, v)
	}
	return out, nil
}
func (r *fakeAPIKeyRepo) GetAPIKeysByUser(_ context.Context, uid string) ([]models.APIKey, error) {
	out := []models.APIKey{}
	for _, v := range r.byID {
		if v.UserID == uid {
			out = append(out, v)
		}
	}
	return out, nil
}
func (r *fakeAPIKeyRepo) DeleteAPIKey(_ context.Context, id string) error {
	delete(r.byID, id)
	return nil
}
func (r *fakeAPIKeyRepo) DeleteAPIKeyForUser(_ context.Context, id, uid string) error {
	if k, ok := r.byID[id]; ok && k.UserID == uid {
		delete(r.byID, id)
		delete(r.byKey, k.Key)
		return nil
	}
	return errors.New("not found")
}
func (r *fakeAPIKeyRepo) IncrementUsage(_ context.Context, id string) error {
	if k, ok := r.byID[id]; ok {
		k.Usage++
		r.byID[id] = k
	}
	return nil
}

// compile-time interface checks
var (
	_ db.UserRepository   = (*fakeUserRepo)(nil)
	_ db.APIKeyRepository = (*fakeAPIKeyRepo)(nil)
)

// ---------- tests ----------

func newAuthSvc() (*services.AuthService, *fakeUserRepo, *fakeAPIKeyRepo) {
	users := newFakeUserRepo()
	keys := newFakeAPIKeyRepo()
	return services.NewAuthService(users, keys, "test-secret"), users, keys
}

func TestRegister_Success(t *testing.T) {
	svc, users, keys := newAuthSvc()

	resp, err := svc.Register(context.Background(), models.RegisterRequest{
		Email:    "Alice@Example.com",
		Password: "password123",
		Name:     "Alice",
	})
	if err != nil {
		t.Fatalf("Register err = %v", err)
	}
	if resp.Token == "" {
		t.Error("expected JWT token")
	}
	if resp.User.Email != "alice@example.com" {
		t.Errorf("email not lowercased: got %q", resp.User.Email)
	}
	if resp.User.PasswordHash != "" {
		t.Error("password hash leaked into response")
	}
	if resp.ApiKey == "" {
		t.Error("expected auto-minted API key in response")
	}
	if len(users.byEmail) != 1 {
		t.Errorf("user not persisted: %d users", len(users.byEmail))
	}
	if len(keys.byID) != 1 {
		t.Errorf("api key not persisted: %d keys", len(keys.byID))
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, _, _ := newAuthSvc()
	req := models.RegisterRequest{Email: "x@y.com", Password: "password123"}

	if _, err := svc.Register(context.Background(), req); err != nil {
		t.Fatalf("first register: %v", err)
	}
	_, err := svc.Register(context.Background(), req)
	if !errors.Is(err, services.ErrEmailTaken) {
		t.Errorf("got %v, want ErrEmailTaken", err)
	}
}

func TestLogin_Success(t *testing.T) {
	svc, _, _ := newAuthSvc()
	req := models.RegisterRequest{Email: "bob@example.com", Password: "password123", Name: "Bob"}
	if _, err := svc.Register(context.Background(), req); err != nil {
		t.Fatal(err)
	}

	resp, err := svc.Login(context.Background(), models.LoginRequest{
		Email: "bob@example.com", Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login err = %v", err)
	}
	if resp.Token == "" {
		t.Error("missing token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, _, _ := newAuthSvc()
	_, _ = svc.Register(context.Background(), models.RegisterRequest{
		Email: "c@d.com", Password: "password123",
	})

	_, err := svc.Login(context.Background(), models.LoginRequest{
		Email: "c@d.com", Password: "wrong-password",
	})
	if !errors.Is(err, services.ErrInvalidCredentials) {
		t.Errorf("got %v, want ErrInvalidCredentials", err)
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	svc, _, _ := newAuthSvc()

	_, err := svc.Login(context.Background(), models.LoginRequest{
		Email: "ghost@nowhere.com", Password: "anything",
	})
	if !errors.Is(err, services.ErrInvalidCredentials) {
		t.Errorf("got %v, want ErrInvalidCredentials", err)
	}
}

func TestMintAPIKey_AssignsToUser(t *testing.T) {
	svc, _, keys := newAuthSvc()
	r, err := svc.Register(context.Background(), models.RegisterRequest{
		Email: "z@z.com", Password: "password123",
	})
	if err != nil {
		t.Fatal(err)
	}

	newKey, err := svc.MintAPIKey(context.Background(), r.User.ID, "Mobile")
	if err != nil {
		t.Fatalf("MintAPIKey err = %v", err)
	}
	got, ok := keys.byKey[newKey]
	if !ok {
		t.Fatal("minted key not persisted")
	}
	if got.UserID != r.User.ID {
		t.Errorf("UserID = %q, want %q", got.UserID, r.User.ID)
	}
	if got.Name != "Mobile" {
		t.Errorf("Name = %q, want %q", got.Name, "Mobile")
	}
}
