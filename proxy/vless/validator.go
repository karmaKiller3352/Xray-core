package vless

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/karmaKiller3352/Xray-core/common/errors"
	"github.com/karmaKiller3352/Xray-core/common/log"
	"github.com/karmaKiller3352/Xray-core/common/protocol"
	"github.com/karmaKiller3352/Xray-core/common/uuid"
)

type Validator interface {
	Get(id uuid.UUID) *protocol.MemoryUser
	Add(u *protocol.MemoryUser) error
	Del(email string) error
	GetByEmail(email string) *protocol.MemoryUser
	GetAll() []*protocol.MemoryUser
	GetCount() int64
}

// MemoryValidator stores valid VLESS users.
type MemoryValidator struct {
	// Considering email's usage here, map + sync.Mutex/RWMutex may have better performance.
	email sync.Map
	users sync.Map
}

// Add a VLESS user, Email must be empty or unique.
func (v *MemoryValidator) Add(u *protocol.MemoryUser) error {
	if u.Email != "" {
		_, loaded := v.email.LoadOrStore(strings.ToLower(u.Email), u)
		if loaded {
			return errors.New("User ", u.Email, " already exists.")
		}
	}
	v.users.Store(u.Account.(*MemoryAccount).ID.UUID(), u)
	return nil
}

// Del a VLESS user with a non-empty Email.
func (v *MemoryValidator) Del(e string) error {
	if e == "" {
		return errors.New("Email must not be empty.")
	}
	le := strings.ToLower(e)
	u, _ := v.email.Load(le)
	if u == nil {
		return errors.New("User ", e, " not found.")
	}
	v.email.Delete(le)
	v.users.Delete(u.(*protocol.MemoryUser).Account.(*MemoryAccount).ID.UUID())
	return nil
}

// Get a VLESS user with UUID, nil if user doesn't exist.
func (v *MemoryValidator) Get(id uuid.UUID) *protocol.MemoryUser {
	u, _ := v.users.Load(id)
	if u != nil {
		return u.(*protocol.MemoryUser)
	}
	return nil
}

// Get a VLESS user with email, nil if user doesn't exist.
func (v *MemoryValidator) GetByEmail(email string) *protocol.MemoryUser {
	email = strings.ToLower(email)
	u, _ := v.email.Load(email)
	if u != nil {
		return u.(*protocol.MemoryUser)
	}
	return nil
}

// Get all users
func (v *MemoryValidator) GetAll() []*protocol.MemoryUser {
	u := make([]*protocol.MemoryUser, 0, 100)
	v.email.Range(func(key, value interface{}) bool {
		u = append(u, value.(*protocol.MemoryUser))
		return true
	})
	return u
}

// Get users count
func (v *MemoryValidator) GetCount() int64 {
	var c int64 = 0
	v.email.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

// APIValidator validates VLESS users via external HTTP API.
type APIValidator struct {
	apiURL      string
	defaultUser *protocol.MemoryUser
}

func NewAPIValidator(apiURL string, defaultUser *protocol.MemoryUser) *APIValidator {
	var level uint32 = 0
	var flow, encryption string
	if defaultUser != nil {
		level = defaultUser.Level
		if acc, ok := defaultUser.Account.(*MemoryAccount); ok {
			flow = acc.Flow
			encryption = acc.Encryption
		}
	}
	return &APIValidator{
		apiURL: apiURL,
		defaultUser: &protocol.MemoryUser{
			Level: level,
			Account: &MemoryAccount{
				Flow:       flow,
				Encryption: encryption,
			},
		},
	}
}

func (v *APIValidator) Add(u *protocol.MemoryUser) error             { return nil }
func (v *APIValidator) Del(email string) error                       { return nil }
func (v *APIValidator) GetByEmail(email string) *protocol.MemoryUser { return nil }
func (v *APIValidator) GetAll() []*protocol.MemoryUser               { return nil }
func (v *APIValidator) GetCount() int64                              { return 0 }
func (v *APIValidator) Get(id uuid.UUID) *protocol.MemoryUser {
	// Формируем правильный URL
	url := v.apiURL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	endpoint := fmt.Sprintf("%s%s", url, id.String())

	client := &http.Client{Timeout: 3 * time.Second}
	log.Record(&log.GeneralMessage{Content: fmt.Sprintf("Попытка аутентификации через API: %s", endpoint)})
	resp, err := client.Get(endpoint)
	if err != nil {
		log.Record(&log.GeneralMessage{Content: fmt.Sprintf("Ошибка запроса к API: %v", err)})
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 || resp.StatusCode == 204 {
		log.Record(&log.GeneralMessage{Content: fmt.Sprintf("API ответил %d, uuid %s аутентифицирован", resp.StatusCode, id.String())})
		user := &protocol.MemoryUser{
			Level: v.defaultUser.Level,
			Account: &MemoryAccount{
				ID:         protocol.NewID(id),
				Flow:       v.defaultUser.Account.(*MemoryAccount).Flow,
				Encryption: v.defaultUser.Account.(*MemoryAccount).Encryption,
			},
		}
		return user
	}
	log.Record(&log.GeneralMessage{Content: fmt.Sprintf("API ответил %d, uuid %s не аутентифицирован", resp.StatusCode, id.String())})
	return nil
}
