package tests

import (
	"testing"

	"github.com/camelhr/camelhr-api/internal/config"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type IntegrationBaseSuite struct {
	suite.Suite
	Config         config.Config
	DB             database.Database
	RedisClient    *redis.Client
	RedisContainer *RedisContainer
	PGContainer    *PostgreSQLContainer
}

func (s *IntegrationBaseSuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("skipping suite in short mode")
	}

	setupDone := false
	defer func() {
		if !setupDone {
			s.TearDownSuite()
		}
	}()

	s.Config = config.Config{
		AppSecret: "test_secret",
	}

	pgContainer, err := NewPostgresContainer()
	s.Require().NoError(err)
	s.PGContainer = pgContainer

	db, err := pgContainer.Connect()
	s.Require().NoError(err)
	s.DB = database.NewPostgresDatabase(db)

	redisContainer, err := NewRedisContainer()
	s.Require().NoError(err)
	s.RedisContainer = redisContainer

	redisClient, err := redisContainer.Connect()
	s.Require().NoError(err)
	s.RedisClient = redisClient

	err = RunMigrations(db.DB)
	s.Require().NoError(err)

	setupDone = true
}

func (s *IntegrationBaseSuite) TearDownSuite() {
	if err := s.PGContainer.Purge(); err != nil {
		s.T().Logf("error purging postgres container: %v", err)
	}

	if err := s.RedisContainer.Purge(); err != nil {
		s.T().Logf("error purging redis container: %v", err)
	}
}
