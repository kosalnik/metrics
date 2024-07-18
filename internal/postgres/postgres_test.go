package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kosalnik/metrics/internal/models"
	"github.com/kosalnik/metrics/internal/postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PostgresSuite struct {
	suite.Suite
	db      *sql.DB
	storage *postgres.DBStorage
}

func (s *PostgresSuite) SetupTest() {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		s.T().Skip("Unable to connect with DB. Env TEST_DSN is not specified")
		return
	}
	db, err := postgres.NewConn(dsn)
	s.Require().NoError(err)
	s.db = db
	s.storage, err = postgres.NewDBStorage(s.db)
	s.Require().NoError(err)
	s.Require().NoError(s.storage.CreateTablesIfNotExist(context.Background()))
}

func TestDBStorage(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}

func (s *PostgresSuite) TestPing() {
	require.NoError(s.T(), s.storage.Ping(context.Background()))
}

func (s *PostgresSuite) TestSetGauge() {
	var val = 123.456
	id := fmt.Sprintf("sg%d", time.Now().Unix())
	m, err := s.storage.SetGauge(context.Background(), id, val)
	s.Require().NoError(err)
	s.Require().Equal(val, m.Value)
	s.Require().Equal(id, m.ID)
	s.Require().Equal(models.MGauge, m.MType)

	got, err := s.storage.GetGauge(context.Background(), id)
	s.Require().NoError(err)
	s.Require().NotNil(got)
	s.Require().Equal(val, got.Value)
}

func (s *PostgresSuite) TestSetCounter() {
	var val int64 = 789
	id := fmt.Sprintf("sc%d", time.Now().Unix())
	m, err := s.storage.IncCounter(context.Background(), id, val)
	s.Require().NoError(err)
	s.Require().Equal(val, m.Delta)
	s.Require().Equal(id, m.ID)
	s.Require().Equal(models.MCounter, m.MType)

	m, err = s.storage.IncCounter(context.Background(), id, val)
	s.Require().NoError(err)
	s.Require().Equal(val+val, m.Delta)

	got, err := s.storage.GetCounter(context.Background(), id)
	s.Require().NoError(err)
	s.Require().NotNil(got)
	s.Require().Equal(val+val, got.Delta)
}

func (s *PostgresSuite) TestBulkActions() {
	list := []models.Metrics{
		{ID: "a", MType: models.MCounter, Delta: 3},
		{ID: "b", MType: models.MCounter, Delta: 4},
		{ID: "c", MType: models.MGauge, Value: 3.2},
		{ID: "d", MType: models.MGauge, Value: 3.3},
	}
	s.Require().NoError(s.storage.UpsertAll(context.Background(), list))
	got, err := s.storage.GetAll(context.Background())
	s.Require().NoError(err)
	ok := true
	for _, wantV := range list {
		found := false

		for _, gotV := range got {
			if wantV.ID == gotV.ID && wantV.MType == gotV.MType {
				found = true
				break
			}
		}
		if !found {
			ok = false
		}
	}
	s.Assert().True(ok, "GetAll is not contain all saved metrics.")
}
