package tests

import (
	"bytes"
	"context"
<<<<<<< HEAD
	"crypto/rsa"
=======
>>>>>>> 75e84a18a119dcf4c0fc171fbe504bdb132894dd
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	configs "wallet-service/internal/config"
	"wallet-service/internal/domain"
<<<<<<< HEAD
	jwtclaims "wallet-service/internal/jwt_claims"
=======
>>>>>>> 75e84a18a119dcf4c0fc171fbe504bdb132894dd
	"wallet-service/internal/repository"
	"wallet-service/internal/repository/psql"
	"wallet-service/internal/service"
	"wallet-service/internal/transport/kafka/producer"
	"wallet-service/internal/transport/rest"

<<<<<<< HEAD
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
=======
>>>>>>> 75e84a18a119dcf4c0fc171fbe504bdb132894dd
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	cancelFunc  context.CancelFunc
	cfg         *configs.Config
<<<<<<< HEAD
	publicKey   *rsa.PublicKey
=======
	cancel      context.CancelFunc
>>>>>>> 75e84a18a119dcf4c0fc171fbe504bdb132894dd
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
<<<<<<< HEAD

	s.cancelFunc = cancel
=======
	s.cancel = cancel
>>>>>>> 75e84a18a119dcf4c0fc171fbe504bdb132894dd

	var err error

	s.cfg, err = configs.Init()
	s.Require().NoError(err)

	s.psql, err = psql.New(s.cfg)
	s.Require().NoError(err)

	err = s.psql.Up()
	s.Require().NoError(err)

	s.usersRepo = repository.NewUsersRepository(s.psql.Database())
	s.walletsRepo = repository.NewWalletRepository(s.psql.Database())

	s.kProducer = producer.New(s.cfg)

	s.services = service.New(s.usersRepo, s.walletsRepo)

	s.publicKey, err = jwtclaims.ReadPublicKey()
	s.Require().NoError(err)

	s.server = rest.New(s.services, s.usersRepo, s.publicKey)

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

<<<<<<< HEAD
func (s *IntegrationTestSuite) getToken(user domain.User) string {
	tokenTime := time.Now().Add(time.Hour * 24)

	claims := jwtclaims.Claims{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	privateKey, err := jwtclaims.ReadPrivateKey()
	s.Require().NoError(err)

	token, err := claims.GenerateToken(privateKey)
	s.Require().NoError(err)

	// logrus.Warnf("claims: %v\n", claims)
	// logrus.Warnf("token: %v\n", token)
	return token
}

func (s *IntegrationTestSuite) sendHTTPRequest(method, path string, statusCode int, entity, result any, user domain.User) {
	clientHTTP := http.Client{}

=======
func (s *IntegrationTestSuite) sendHTTPRequest(method, path string, statusCode int, entity, result any, user domain.User) {
	clientHTTP := http.Client{}
	
>>>>>>> 75e84a18a119dcf4c0fc171fbe504bdb132894dd
	entityJSON, err := json.Marshal(entity)
	s.Require().NoError(err)

	url := fmt.Sprintf("http://localhost:%s%s", s.cfg.HTTP.Port, path)

<<<<<<< HEAD
	req, err := http.NewRequestWithContext(context.Background(), method, url, bytes.NewReader(entityJSON))
	s.Require().NoError(err, "failed to create new request")

	token := s.getToken(user)
	req.Header.Set("Authorization", "Bearer "+token)

=======
	req, err := http.NewRequest(method, url, bytes.NewReader(entityJSON))
	s.Require().NoError(err, "failed to create new request")

>>>>>>> 75e84a18a119dcf4c0fc171fbe504bdb132894dd
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
<<<<<<< HEAD

		return
=======
>>>>>>> 75e84a18a119dcf4c0fc171fbe504bdb132894dd
	}

	if result == nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	s.Require().NoError(err)
}
