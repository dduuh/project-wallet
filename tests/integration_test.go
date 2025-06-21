package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
	configs "wallet-service/internal/config"
	"wallet-service/internal/repository"
	"wallet-service/internal/repository/psql"
	"wallet-service/internal/service"
	"wallet-service/internal/transport/kafka/producer"
	"wallet-service/internal/transport/rest"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg         *configs.Config
	psql        *psql.PostgresDB
	usersRepo   *repository.UsersRepository
	walletsRepo *repository.WalletDB
	services    *service.Service
	server      *rest.Server
	kProducer   *producer.Producer
}

func (s *IntegrationTestSuite) SetupSuite() {
	var err error

	s.cfg, err = configs.Init()
	s.Require().NoError(err)

	s.psql, err = psql.New(s.cfg)
	s.Require().NoError(err)

	err = s.psql.Up()
	s.Require().NoError(err)

	s.kProducer = producer.New(s.cfg)

	s.services = service.New(s.usersRepo, s.walletsRepo)

	s.server = rest.New(s.services, s.usersRepo)

	//nolint:testifylint
	go func() {
		err := s.server.Run(s.cfg, s.server.InitRoutes())
		s.Require().NoError(err)
	}()
}

func TestIntegrationSetupSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(IntegrationTestSuite))
}
