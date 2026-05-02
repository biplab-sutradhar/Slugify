package services_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/biplab-sutradhar/slugify/api/internal/cache"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/idgen"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
)

// ---------- fakes ----------

type fakeLinkRepo struct {
	mu    sync.Mutex
	links map[string]models.Link // by id
}

func newFakeLinkRepo() *fakeLinkRepo {
	return &fakeLinkRepo{links: map[string]models.Link{}}
}

func (r *fakeLinkRepo) CreateLink(l models.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.links[l.ID] = l
	return nil
}
func (r *fakeLinkRepo) GetLinkByShortCode(code string) (models.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, l := range r.links {
		if l.ShortCode == code {
			return l, nil
		}
	}
	return models.Link{}, errors.New("not found")
}
func (r *fakeLinkRepo) GetLinkByIDForUser(id, userID string) (models.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if l, ok := r.links[id]; ok && l.UserID == userID {
		return l, nil
	}
	return models.Link{}, errors.New("not found")
}
func (r *fakeLinkRepo) ListLinksForUser(userID string, limit, offset int) ([]models.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []models.Link{}
	for _, l := range r.links {
		if l.UserID == userID {
			out = append(out, l)
		}
	}
	if offset >= len(out) {
		return []models.Link{}, nil
	}
	end := offset + limit
	if end > len(out) {
		end = len(out)
	}
	return out[offset:end], nil
}
func (r *fakeLinkRepo) UpdateLinkStatusForUser(id, userID string, isActive bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if l, ok := r.links[id]; ok && l.UserID == userID {
		l.IsActive = isActive
		r.links[id] = l
		return nil
	}
	return errors.New("not found")
}
func (r *fakeLinkRepo) DeleteLinkForUser(id, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if l, ok := r.links[id]; ok && l.UserID == userID {
		delete(r.links, id)
		return nil
	}
	return errors.New("not found")
}
func (r *fakeLinkRepo) IncrementClicks(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, l := range r.links {
		if l.ShortCode == code {
			l.Clicks++
			r.links[id] = l
			return nil
		}
	}
	return errors.New("not found")
}

type fakeCache struct {
	mu   sync.Mutex
	data map[string]string
}

func newFakeCache() *fakeCache { return &fakeCache{data: map[string]string{}} }
func (c *fakeCache) GetURL(_ context.Context, code string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data[code], nil
}
func (c *fakeCache) SetURL(_ context.Context, code, longURL string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[code] = longURL
	return nil
}
func (c *fakeCache) Close() error { return nil }

type stubTicket struct {
	codes []string
	idx   int
}

func newStubTicket(codes ...string) *stubTicket { return &stubTicket{codes: codes} }
func (t *stubTicket) GenerateID(_ context.Context) (string, error) {
	if t.idx >= len(t.codes) {
		return "", errors.New("exhausted")
	}
	c := t.codes[t.idx]
	t.idx++
	return c, nil
}
func (t *stubTicket) Close() error { return nil }

// compile-time interface checks
var (
	_ db.LinkRepository  = (*fakeLinkRepo)(nil)
	_ cache.Cache        = (*fakeCache)(nil)
	_ idgen.TicketServer = (*stubTicket)(nil)
)

// ---------- tests ----------

func newLinkSvc(codes ...string) (*services.LinkService, *fakeLinkRepo, *fakeCache) {
	repo := newFakeLinkRepo()
	cch := newFakeCache()
	svc := services.NewLinkService(repo, cch, newStubTicket(codes...), newFakeAPIKeyRepo(), "<https://slugify.test>")
	return svc, repo, cch
}

func TestSaveLink_RequiresLongURL(t *testing.T) {
	svc, _, _ := newLinkSvc("abc")
	if _, err := svc.SaveLink("user-1", ""); err == nil {
		t.Error("expected error for empty long_url")
	}
}

func TestSaveLink_RequiresUserID(t *testing.T) {
	svc, _, _ := newLinkSvc("abc")
	if _, err := svc.SaveLink("", "<https://example.com>"); err == nil {
		t.Error("expected error for empty user_id")
	}
}

func TestSaveLink_PersistsAndCaches(t *testing.T) {
	svc, repo, cch := newLinkSvc("abc")
	link, err := svc.SaveLink("user-1", "<https://example.com>")
	if err != nil {
		t.Fatalf("SaveLink err = %v", err)
	}
	if link.UserID != "user-1" {
		t.Errorf("UserID = %q, want %q", link.UserID, "user-1")
	}
	if link.ShortCode != "abc" {
		t.Errorf("ShortCode = %q, want %q", link.ShortCode, "abc")
	}
	if got, _ := repo.GetLinkByShortCode("abc"); got.ID != link.ID {
		t.Error("link not persisted in repo")
	}
	if cached, _ := cch.GetURL(context.Background(), "abc"); cached != "<https://example.com>" {
		t.Errorf("cache value = %q, want %q", cached, "<https://example.com>")
	}
}

func TestListLinks_OnlyReturnsCallerLinks(t *testing.T) {
	svc, repo, _ := newLinkSvc("a1", "a2", "b1")

	if _, err := svc.SaveLink("alice", "<https://alice-1.com>"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SaveLink("alice", "<https://alice-2.com>"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SaveLink("bob", "<https://bob-1.com>"); err != nil {
		t.Fatal(err)
	}

	if total := len(repo.links); total != 3 {
		t.Fatalf("repo has %d links, want 3", total)
	}

	got, err := svc.ListLinks("alice", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("alice's links count = %d, want 2", len(got))
	}
	for _, l := range got {
		if l.UserID != "alice" {
			t.Errorf("foreign link in alice's list: %+v", l)
		}
	}
}

func TestUpdateAndDelete_RejectForeignUser(t *testing.T) {
	svc, _, _ := newLinkSvc("a1")
	link, _ := svc.SaveLink("alice", "<https://example.com>")

	if err := svc.UpdateLinkStatus(link.ID, "bob", false); err == nil {
		t.Error("bob should not be able to update alice's link")
	}
	if err := svc.DeleteLink(link.ID, "bob"); err == nil {
		t.Error("bob should not be able to delete alice's link")
	}
}

func TestListLinks_NormalizesLimit(t *testing.T) {
	svc, _, _ := newLinkSvc("a", "b", "c", "d")
	for i := 0; i < 4; i++ {
		_, _ = svc.SaveLink("u", "<https://x.com>")
	}
	// limit 0 -> defaults to 20 (which covers all 4)
	got, err := svc.ListLinks("u", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 4 {
		t.Errorf("limit=0 returned %d, want 4 (defaulted to 20)", len(got))
	}
}
