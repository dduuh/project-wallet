package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	configs "wallet-service/internal/config"
	"wallet-service/internal/domain"
	jwtclaims "wallet-service/internal/jwt_claims"
	"wallet-service/internal/repository"
	"wallet-service/internal/repository/psql"
	"wallet-service/internal/service"
	"wallet-service/internal/transport/kafka/producer"
	"wallet-service/internal/transport/rest"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg         *configs.Config
	cancel      context.CancelFunc
	psql        *psql.PostgresDB
	usersRepo   *repository.UsersRepository
	walletsRepo *repository.WalletDB
	services    *service.Service
	server      *rest.Server
	kProducer   *producer.Producer
}

func (s *IntegrationTestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.cancel = cancel

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
		err := s.server.Run(ctx, s.cfg, s.server.InitRoutes())
		s.Require().NoError(err)
	}()

	time.Sleep(time.Millisecond * 50)
}

func TestIntegrationSetupSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) generateToken(user domain.User) string {
	claims := jwtclaims.New()

	privateKey, err := jwtclaims.ReadPrivateKey()
	s.Require().NoError(err)

	token, err := claims.GenerateToken(privateKey)
	s.Require().NoError(err)

	return token
}

func (s *IntegrationTestSuite) sendHTTPRequest(method, path string, statusCode int, entity, result any, user domain.User) {
	clientHTTP := http.Client{}

	entityJSON, err := json.Marshal(entity)
	s.Require().NoError(err)

	url := fmt.Sprintf("http://localhost:%s%s", s.cfg.HTTP.Port, path)

	req, err := http.NewRequest(method, url, bytes.NewReader(entityJSON))
	s.Require().NoError(err, "failed to create new request")

	token := s.generateToken(user)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := clientHTTP.Do(req)
	s.Require().NoError(err)

	defer func() {
		err := resp.Body.Close()
		s.Require().NoError(err)
	}()

	if statusCode != resp.StatusCode {
		respBody, err := io.ReadAll(resp.Body)
		s.Require().NoError(err)

		s.T().Logf("response body: %s", string(respBody))

		s.Require().Equal(statusCode, resp.StatusCode, "unexpected status code")
	}

	if result == nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	s.Require().NoError(err)
}
