package api

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

	ctxlog.Info(sc.Ctx, "created new key", "key_id", key.ID)

	return (*AuthKey)(key), nil
}

func (sc *StateCtx) LoadAuthKeys(createFirst bool) error {
	keys, err := sc.Queries.FindKeys(sc.Ctx)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		if !createFirst {
			return ErrNoKeys
		}
		if key, err := sc.CreateAuthKey(); err != nil {
			return fmt.Errorf("failed to create first key: %w", err)
		} else {
			keys = append(keys, (*db.Key)(key))
		}
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
	AccountName string `json:"n"`
	AccountID   int32  `json:"i"`
	KeyID       int32  `json:"k"`
	CreatedAt   int64  `json:"t"` // unix timestamp in seconds
}

func (t *AuthToken) TTL(ttl int64) int64 {
	return t.CreatedAt + ttl - time.Now().Unix()
}

func (t *AuthToken) String() string {
	return fmt.Sprintf("<%d:%d:%d>", t.AccountID, t.KeyID, t.CreatedAt)
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

func CreateAuthToken(key *AuthKey, accountName string, accountID int32) *AuthToken {
	return &AuthToken{
		AccountName: accountName,
		AccountID:   accountID,
		KeyID:       key.ID,
		CreatedAt:   time.Now().Unix(),
	}
}

func EncodeAuthToken(key *AuthKey, token *AuthToken) (string, error) {
	tokenData, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("failed to encode auth token: %w", err)
	}

	hmacData, err := ComputeHMAC(tokenData, key.Data, nil)
	if err != nil {
		return "", err
	}

	tokenPart := base64.RawURLEncoding.EncodeToString(tokenData)
	hmacPart := base64.RawURLEncoding.EncodeToString(hmacData)
	return fmt.Sprintf("%s.%s", tokenPart, hmacPart), nil
}

func DecodeAuthToken(str string) (*AuthToken, []byte, []byte, error) {
	splitAt := strings.IndexByte(str, '.')
	if splitAt < 0 {
		return nil, nil, nil, ErrMalformedToken.WithInternal(
			errors.New("no part separator found"),
		)
	}

	tokenPart := str[:splitAt]
	hmacPart := str[splitAt+1:]

	tokenData, err := base64.RawURLEncoding.DecodeString(tokenPart)
	if err != nil {
		return nil, nil, nil, ErrMalformedToken.WithInternal(
			fmt.Errorf("failed to decode token part: %w", err),
		)
	}

	hmacData, err := base64.RawURLEncoding.DecodeString(hmacPart)
	if err != nil {
		return nil, nil, nil, ErrMalformedToken.WithInternal(
			fmt.Errorf("failed to decode hmac part: %w", err),
		)
	}

	token := AuthToken{}
	if err := json.Unmarshal(tokenData, &token); err != nil {
		return nil, nil, nil, ErrMalformedToken.WithInternal(
			fmt.Errorf("failed to decode auth token: %w", err),
		)
	}

	return &token, tokenData, hmacData, nil
}

func (sc *StateCtx) CreateAuthCookie(accountName string, accountID int32) (*http.Cookie, error) {
	key := sc.Auth.KeyInUse.Load()
	token := CreateAuthToken(key, accountName, accountID)
	tokenStr, err := EncodeAuthToken(key, token)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:  sc.Config.Auth.Cookie.V,
		Value: tokenStr,
		Path:  "/",
		Expires: time.Unix(token.CreatedAt, 0).Add(
			time.Second * time.Duration(sc.Config.Auth.TokenTTL.V),
		),
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

	if token.TTL(s.Config.Auth.TokenTTL.V) <= 0 {
		return token, ErrExpiredToken
	}

	return token, nil
}

func (s *State) CreateAuthState(tokenStr string) *AuthState {
	state := AuthState{}
	state.Token, state.Error = s.CheckAuthToken(tokenStr)
	return &state
}
