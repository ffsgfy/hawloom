package api

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

type AuthKey db.Key

type Auth struct {
	Lock     sync.RWMutex // only for the keystore
	Keystore map[int32]*AuthKey

	KeyInUse atomic.Pointer[AuthKey]
}

func NewAuth() *Auth {
	auth := &Auth{
		Keystore: map[int32]*AuthKey{},
	}
	auth.KeyInUse.Store(&AuthKey{})
	return auth
}

func (a *Auth) GetKey(id int32) *AuthKey {
	a.Lock.RLock()
	defer a.Lock.RUnlock()
	return a.Keystore[id]
}

func (sc *StateCtx) CreateAuthKey() (*AuthKey, error) {
	data := make([]byte, sc.Config.Auth.KeySize.V)
	if _, err := rand.Reader.Read(data); err != nil {
		return nil, fmt.Errorf("failed to generate key data: %w", err)
	}

	key, err := sc.Queries.CreateKey(sc.Ctx, data)
	if err != nil {
		return nil, err
	}

	return (*AuthKey)(key), nil
}

func (sc *StateCtx) LoadAuthKeys(required bool) error {
	keys, err := sc.Queries.FindKeys(sc.Ctx)
	if err != nil {
		return err
	}

	if required && len(keys) == 0 {
		return errors.New("no keys in db")
	}

	sc.Auth.Lock.Lock()
	defer sc.Auth.Lock.Unlock()
	clear(sc.Auth.Keystore)

	keyIDs := make([]int32, 0, len(keys))
	keyInUse := sc.Auth.KeyInUse.Load()

	for i, key := range keys {
		keyIDs = append(keyIDs, key.ID)
		authKey := (*AuthKey)(key)

		sc.Auth.Keystore[authKey.ID] = authKey
		if i == 0 || authKey.ID > keyInUse.ID {
			keyInUse = authKey
		}
	}

	sc.Auth.KeyInUse.Store(keyInUse)

	ctxlog.Info(
		sc.Ctx, "loaded keys from db",
		"key_count", len(keys),
		"key_ids", keyIDs,
		"key_in_use", keyInUse.ID,
	)

	return nil
}

type AuthToken struct {
	AccountID int32
	KeyID     int32
	Expires   int64 // unix timestamp in seconds
}

func (t *AuthToken) TTL() int64 {
	return t.Expires - time.Now().Unix()
}

func (t *AuthToken) String() string {
	return fmt.Sprintf("<%d:%d:%d>", t.AccountID, t.KeyID, t.Expires)
}

func ComputeHMAC(data, key, out []byte) ([]byte, error) {
	hm := hmac.New(sha256.New, key)
	if _, err := hm.Write(data); err != nil {
		return nil, fmt.Errorf("failed to compute HMAC: %w", err)
	}

	return hm.Sum(out), nil
}

func CheckHMAC(data, key, hm []byte) (bool, error) {
	dataHM, err := ComputeHMAC(data, key, nil)
	if err != nil {
		return false, err
	}

	return hmac.Equal(dataHM, hm), nil
}

func CreateAuthToken(key *AuthKey, accountID int32, ttl int) *AuthToken {
	return &AuthToken{
		AccountID: accountID,
		KeyID:     key.ID,
		Expires:   time.Now().Unix() + int64(ttl),
	}
}

func EncodeAuthToken(key *AuthKey, token *AuthToken) (string, error) {
	data, err := binary.Append(nil, binary.BigEndian, token)
	if err != nil {
		return "", fmt.Errorf("failed to encode auth token: %w", err)
	}

	data, err = ComputeHMAC(data, key.Data, data)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func DecodeAuthToken(str string) (*AuthToken, []byte, []byte, error) {
	data, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, nil, nil, ErrMalformedToken.WithInternal(
			fmt.Errorf("failed to decode base64: %w", err),
		)
	}

	token := AuthToken{}
	size, err := binary.Decode(data, binary.BigEndian, &token)
	if err != nil {
		return nil, nil, nil, ErrMalformedToken.WithInternal(
			fmt.Errorf("failed to decode auth token: %w", err),
		)
	}

	return &token, data[:size], data[size:], nil
}

func CreateAuthCookie(key *AuthKey, token *AuthToken, name string) (*http.Cookie, error) {
	tokenStr, err := EncodeAuthToken(key, token)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     name,
		Value:    tokenStr,
		MaxAge:   int(token.TTL()),
		SameSite: http.SameSiteStrictMode,
	}, nil
}

type AuthState struct {
	Token *AuthToken
	Error error
}

func (a *AuthState) Valid() bool {
	return a.Error == nil && a.Token != nil
}

type authStateKeyType struct{}

var authStateKey = authStateKeyType{}

func WithAuthState(ctx context.Context, authState *AuthState) context.Context {
	return context.WithValue(ctx, authStateKey, authState)
}

func GetAuthState(ctx context.Context) *AuthState {
	if authState, ok := ctx.Value(authStateKey).(*AuthState); ok {
		return authState
	}
	return nil
}

func GetValidAuthToken(ctx context.Context) (*AuthToken, error) {
	if authState := GetAuthState(ctx); authState != nil {
		if authState.Valid() {
			return authState.Token, nil
		}
		if authState.Error != nil {
			return nil, authState.Error
		}
	}
	return nil, ErrUnauthorized
}

func (s *State) CheckAuthToken(tokenStr string) (*AuthToken, error) {
	token, data, hm, err := DecodeAuthToken(tokenStr)
	if err != nil {
		return nil, err
	}

	key := s.Auth.GetKey(token.KeyID)
	if key == nil {
		return token, ErrNoTokenKey
	}

	if ok, err := CheckHMAC(data, key.Data, hm); err != nil {
		return token, err
	} else if !ok {
		return token, ErrWrongTokenHash
	}

	if token.TTL() <= 0 {
		return token, ErrExpiredToken
	}

	return token, nil
}

func (s *State) CreateAuthState(tokenStr string) *AuthState {
	state := AuthState{}
	state.Token, state.Error = s.CheckAuthToken(tokenStr)
	return &state
}
