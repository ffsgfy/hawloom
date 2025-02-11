package api

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"

	"github.com/ffsgfy/hawloom/internal/config"
)

type TestDBSuite struct {
	suite.Suite
	m *migrate.Migrate
	s *State
}

func (suite *TestDBSuite) SetupSuite() {
	config, err := config.Load("./config/ut.json")
	suite.Require().NoError(err)

	dbURI := config.DB.MakePostgresURI()
	suite.m, err = migrate.New("file://../../db/migrations/", dbURI)
	suite.Require().NoError(err)

	_, dirty, err := suite.m.Version()
	suite.Require().NoError(err)
	suite.Require().False(dirty)
	if err = suite.m.Up(); !errors.Is(err, migrate.ErrNoChange) {
		suite.Require().NoError(err)
	}

	suite.s, err = NewState(context.Background(), config)
	suite.Require().NoError(err)
}
