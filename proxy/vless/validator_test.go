package vless

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karmaKiller3352/Xray-core/common/protocol"
	"github.com/karmaKiller3352/Xray-core/common/uuid"
)

func TestAPIValidator_Success(t *testing.T) {
	// Мокаем API, который всегда возвращает 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	defaultUser := &protocol.MemoryUser{
		Level: 42,
		Account: &MemoryAccount{
			Flow:       "xtls-rprx-vision",
			Encryption: "none",
		},
	}
	validator := NewAPIValidator(server.URL+"/", defaultUser)

	id := uuid.New()
	user := validator.Get(id)
	if user == nil {
		t.Fatal("Ожидался успешный пользователь, но вернулось nil")
	}
	if user.Level != 42 {
		t.Errorf("Level должен быть 42, а не %d", user.Level)
	}
	acc := user.Account.(*MemoryAccount)
	if acc.Flow != "xtls-rprx-vision" {
		t.Errorf("Flow должен быть xtls-rprx-vision, а не %s", acc.Flow)
	}
	if acc.Encryption != "none" {
		t.Errorf("Encryption должен быть none, а не %s", acc.Encryption)
	}
}

func TestAPIValidator_Fail(t *testing.T) {
	// Мокаем API, который всегда возвращает 403
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
	}))
	defer server.Close()

	defaultUser := &protocol.MemoryUser{
		Level: 1,
		Account: &MemoryAccount{
			Flow:       "",
			Encryption: "none",
		},
	}
	validator := NewAPIValidator(server.URL+"/", defaultUser)

	id := uuid.New()
	user := validator.Get(id)
	if user != nil {
		t.Fatal("Ожидалось nil при неуспешной аутентификации, но вернулся пользователь")
	}
}
