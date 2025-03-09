package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type DBConfig struct {
	Host     Value[string] `json:"host"`
	Port     Value[uint16] `json:"port"`
	User     Value[string] `json:"user"`
	Password Value[string] `json:"password"`
	Database Value[string] `json:"database"`
}

func (c *DBConfig) MakePostgresURI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User.V, c.Password.V, c.Host.V, c.Port.V, c.Database.V,
	)
}

type HTTPConfig struct {
	BindAddress Value[string] `json:"bind_address"`
	BindPort    Value[uint16] `json:"bind_port"`
}

type AuthConfig struct {
	KeySize  Value[int]    `json:"key_size"`
	TokenTTL Value[int]    `json:"token_ttl"`
	Cookie   Value[string] `json:"cookie"`
}

type AccountConfig struct {
	NameMinLength     Value[int] `json:"name_min_length"`
	NameMaxLength     Value[int] `json:"name_max_length"`
	PasswordMinLength Value[int] `json:"password_min_length"`
	PasswordMaxLength Value[int] `json:"password_max_length"`
}

type DocConfig struct {
	TitleMinLength Value[int] `json:"title_min_length"`
	TitleMaxLength Value[int] `json:"title_max_length"`
}

type VordConfig struct {
	MinDuration       Value[int32]   `json:"min_duration"`
	DurationExtension Value[float64] `json:"duration_extension"`
	AutocommitPeriod  Value[int32]   `json:"autocommit_period"`
}

type MiscConfig struct {
	LogLevel Value[int] `json:"log_level"`
}

type Config struct {
	DB      DBConfig      `json:"db"`
	HTTP    HTTPConfig    `json:"http"`
	Auth    AuthConfig    `json:"auth"`
	Account AccountConfig `json:"account"`
	Doc     DocConfig     `json:"doc"`
	Vord    VordConfig    `json:"vord"`
	Misc    MiscConfig    `json:"misc"`
}

func Parse(data []byte) (*Config, error) {
	var result Config
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func Load(path string) (*Config, error) {
	if len(path) == 0 {
		return nil, errors.New("empty config path")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return Parse(data)
}
