package api

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

const (
	KeySize    = 32
	TokenTTL   = 60 * 60 * 8
	AuthCookie = "auth"
)

func ComputeHMAC(data, key, out []byte) ([]byte, error) {
	hm := hmac.New(sha256.New, key)
	_, err := hm.Write(data)
	if err != nil {
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

type AuthKey db.Key

func (sc *StateCtx) CreateKey() (*AuthKey, error) {
	data := make([]byte, KeySize)
	_, err := rand.Reader.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key data: %w", err)
	}

	key, err := sc.Queries.CreateKey(sc.Ctx, data)
	if err != nil {
		return nil, err
	}

	return (*AuthKey)(key), nil
}

func (sc *StateCtx) LoadKeys(required bool) error {
	keys, err := sc.Queries.FindKeys(sc.Ctx)
	if err != nil {
		return err
	}

	if required && len(keys) == 0 {
		return errors.New("no keys in db")
	}

	sc.KeysLock.Lock()
	defer sc.KeysLock.Unlock()

	clear(sc.Keys)
	sc.KeyInUse = 0

	keyIDs := make([]int32, 0, len(keys))
	for i, key := range keys {
		keyIDs = append(keyIDs, key.ID)

		sc.Keys[key.ID] = (*AuthKey)(key)
		if i == 0 {
			sc.KeyInUse = key.ID
		} else {
			sc.KeyInUse = max(sc.KeyInUse, key.ID)
		}
	}

	ctxlog.Info(
		sc.Ctx, "loaded keys from db",
		"key_count", len(keys),
		"key_ids", keyIDs,
		"key_in_use", sc.KeyInUse,
	)

	return nil
}

func (s *State) GetKey(id int32) *AuthKey {
	s.KeysLock.RLock()
	defer s.KeysLock.RUnlock()
	return s.Keys[id]
}

func (s *State) GetKeyInUse() *AuthKey {
	s.KeysLock.RLock()
	defer s.KeysLock.RUnlock()
	return s.Keys[s.KeyInUse]
}

type Token struct {
	KeyID     int32
	AccountID int32
	Expires   int64 // unix timestamp in seconds
}

func (t *Token) TTL() int64 {
	return t.Expires - time.Now().Unix()
}

func (k *AuthKey) CreateToken(accountID int32) *Token {
	return &Token{
		KeyID:     k.ID,
		AccountID: accountID,
		Expires:   time.Now().Unix() + TokenTTL,
	}
}

func (k *AuthKey) EncodeToken(token *Token) (string, error) {
	data, err := binary.Append(nil, binary.BigEndian, token)
	if err != nil {
		return "", fmt.Errorf("failed to encode token: %w", err)
	}

	data, err = ComputeHMAC(data, k.Data, data)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func DecodeToken(str string) (*Token, []byte, []byte, error) {
	data, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, nil, nil, ErrMalformedToken.WithInternal(
			fmt.Errorf("failed to decode base64: %w", err),
		)
	}

	token := Token{}
	size, err := binary.Decode(data, binary.BigEndian, &token)
	if err != nil {
		return nil, nil, nil, ErrMalformedToken.WithInternal(
			fmt.Errorf("failed to decode token: %w", err),
		)
	}

	return &token, data[:size], data[size:], nil
}

func (sc *StateCtx) CheckToken(str string) (*Token, *db.Account, error) {
	token, data, hm, err := DecodeToken(str)
	if err != nil {
		return nil, nil, err
	}

	key := sc.GetKey(token.KeyID)
	if key == nil {
		return nil, nil, ErrNoTokenKey
	}

	ok, err := CheckHMAC(data, key.Data, hm)
	if err != nil {
		return nil, nil, err
	} else if !ok {
		return nil, nil, ErrWrongTokenKey
	}

	if token.TTL() < 0 {
		return nil, nil, ErrExpiredToken
	}

	account, err := sc.FindAccount(&token.AccountID, nil)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return nil, nil, ErrUnauthorized.WithInternal(err)
		}
		return nil, nil, err
	}

	return token, account, nil
}
