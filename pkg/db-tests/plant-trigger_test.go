//go:build integration

package dbtests_test

import (
	"PlantSite/internal/infra/sqpgx"
	filestorage "PlantSite/internal/repositories/pgminio/file-storage"
	"PlantSite/internal/repositories/tests"
	"PlantSite/internal/testutils/pgtest"
	"context"
	"os"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

var (
	testPlantUUID = uuid.New()
)

type TestPlantTriggerSuite struct {
	suite.Suite
	dbContainer testcontainers.Container
	db          *sqpgx.SquirrelPgx
	fileRepo    *filestorage.PgMinioStorage
	prevDir     string
	photoId     uuid.UUID
}

func TestPlantTrigger(t *testing.T) {
	suite.Run(t, new(TestPlantTriggerSuite))
}

func (s *TestPlantTriggerSuite) SetupSuite() {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(s.T(), err)
	s.prevDir = prevDir

	// Change directory to test working directory
	err = os.Chdir(tests.GetTestWorkingDir())
	require.NoError(s.T(), err)

	// Setup PostgreSQL container
	dbContainer, dbCreds, err := pgtest.NewTestPostgres(ctx)
	require.NoError(s.T(), err)
	s.dbContainer = dbContainer

	// Run migrations
	err = pgtest.Migrate(ctx, &dbCreds)
	require.NoError(s.T(), err)

	// Create database connection
	dbConfig := &sqpgx.SqpgxConfig{
		User:                   dbCreds.User,
		Password:               dbCreds.Password,
		DbName:                 dbCreds.Database,
		Host:                   dbCreds.Host,
		Port:                   dbCreds.Port,
		MaxConnections:         10,
		MaxConnectionsLifetime: time.Minute,
	}
	s.db, err = sqpgx.NewSquirrelPgx(ctx, dbConfig)
	require.NoError(s.T(), err)
}

func (s *TestPlantTriggerSuite) TearDownSuite() {
	ctx := context.Background()
	if s.dbContainer != nil {
		s.dbContainer.Terminate(ctx)
	}
	os.Chdir(s.prevDir)
}

func (s *TestPlantTriggerSuite) TearDownTest() {
	_, err := s.db.Delete(context.Background(), squirrel.Delete("plant").Suffix("CASCADE"))
	require.NoError(s.T(), err)
	_, err = s.db.Delete(context.Background(), squirrel.Delete("file").Suffix("CASCADE"))
	require.NoError(s.T(), err)
	_, err = s.db.Delete(context.Background(), squirrel.Delete("plant_category").Suffix("CASCADE"))
	require.NoError(s.T(), err)
}

func (s *TestPlantTriggerSuite) SetupTest() {
	// create test category
	_, err := s.db.Insert(context.Background(), squirrel.Insert("plant_category").
		Columns("name", "attributes").
		Values("test_category", `{
			"int_field": {"type": "number", "min": 1, "max": 10},
			"float_field": {"type": "float", "min": 0.5},
			"select_field": {"type": "string", "options": ["opt1", "opt2"]},
			"string_field": {"type": "string"}
		}`))
	require.NoError(s.T(), err)

	// create test file for main photo
	s.photoId = uuid.New()
	_, err = s.db.Insert(context.Background(), squirrel.Insert("file").
		Columns("id", "name", "url").
		Values(s.photoId, "test_photo.jpg", "https://test.com/test_photo.jpg"))
	require.NoError(s.T(), err)

	// create test plant
	_, err = s.db.Insert(context.Background(), squirrel.Insert("plant").
		Columns("id", "name", "latin_name", "description", "category", "main_photo_id", "specification").
		Values(testPlantUUID, "Test Plant", "Testus plantus", "Test description", "test_category", s.photoId, `{
			"int_field": 5,
			"float_field": 1.5,
			"select_field": "opt1",
			"string_field": "any string value"
		}`))
	require.NoError(s.T(), err)
}

func (s *TestPlantTriggerSuite) TestPlantTriggerInsert() {
	var testCases = []struct {
		name          string
		category      string
		specification map[string]interface{}
		expectingErr  bool
	}{
		{
			name:     "Valid plant",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: false,
		},
		{
			name:          "Invalid category",
			category:      "invalid_category",
			specification: map[string]interface{}{},
			expectingErr:  true,
		},
		{
			name:     "missing required select field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "missing required int field",
			category: "test_category",
			specification: map[string]interface{}{
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "extra field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
				"extra_field":  "extra value",
			},
			expectingErr: true,
		},
		{
			name:     "invalid int field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    "invalid",
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "float integer field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5.5,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "below min int field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    0,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "above max int field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    11,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "invalid float field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  "invalid",
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "below min float field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  0.4,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "invalid select field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"select_field": 2,
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "invalid select option",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"select_field": "invalid",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "invalid string field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": 2,
			},
			expectingErr: true,
		},
		{
			name:          "nil specification",
			category:      "test_category",
			specification: nil,
			expectingErr:  true,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			plantID := uuid.New()
			_, err := s.db.Insert(context.Background(), squirrel.Insert("plant").
				Columns("id", "name", "latin_name", "description", "category", "main_photo_id", "specification").
				Values(plantID, "Test Plant name", "Testus plantus", "Test description", testCase.category, s.photoId, testCase.specification))
			if testCase.expectingErr {
				require.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)
				var name string
				row, err := s.db.QueryRow(context.Background(), squirrel.Select("name").From("plant").Where(squirrel.Eq{"id": plantID}))
				require.NoError(s.T(), err)
				err = row.Scan(&name)
				require.NoError(s.T(), err)
				require.Equal(s.T(), name, "Test Plant name")
			}
		})
	}
}

func (s *TestPlantTriggerSuite) TestPlantTriggerUpdate() {
	var testCases = []struct {
		name          string
		category      string
		specification map[string]interface{}
		expectingErr  bool
	}{
		{
			name:     "Valid plant",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: false,
		},
		{
			name:          "Invalid category",
			category:      "invalid_category",
			specification: map[string]interface{}{},
			expectingErr:  true,
		},
		{
			name:     "missing required select field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "missing required int field",
			category: "test_category",
			specification: map[string]interface{}{
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "extra field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
				"extra_field":  "extra value",
			},
			expectingErr: true,
		},
		{
			name:     "invalid int field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    "invalid",
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "float integer field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5.5,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "below min int field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    0,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "above max int field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    11,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "invalid float field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  "invalid",
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "below min float field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  0.4,
				"select_field": "opt1",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "invalid select field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"select_field": 2,
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "invalid select option",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"select_field": "invalid",
				"string_field": "any string value",
			},
			expectingErr: true,
		},
		{
			name:     "invalid string field",
			category: "test_category",
			specification: map[string]interface{}{
				"int_field":    5,
				"float_field":  1.5,
				"select_field": "opt1",
				"string_field": 2,
			},
			expectingErr: true,
		},
		{
			name:          "nil specification",
			category:      "test_category",
			specification: nil,
			expectingErr:  true,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			_, err := s.db.Update(context.Background(), squirrel.Update("plant").
				SetMap(map[string]interface{}{
					"name":          "Test Plant updated",
					"latin_name":    "Testus plantus",
					"description":   "Test description",
					"category":      testCase.category,
					"main_photo_id": s.photoId,
					"specification": testCase.specification,
				}).
				Where(squirrel.Eq{"id": testPlantUUID}))
			if testCase.expectingErr {
				require.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)
				var name string
				row, err := s.db.QueryRow(context.Background(), squirrel.Select("name").From("plant").Where(squirrel.Eq{"id": testPlantUUID}))
				require.NoError(s.T(), err)
				err = row.Scan(&name)
				require.NoError(s.T(), err)
				require.Equal(s.T(), name, "Test Plant updated")
			}
		})

	}
}
