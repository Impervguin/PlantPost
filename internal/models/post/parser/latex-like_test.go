//go:build unit

package parser_test

import (
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/post/parser"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockPlantGetter struct {
	mock.Mock
}

func (m *MockPlantGetter) GetPlants(uuids []uuid.UUID) ([]*plant.Plant, error) {
	args := m.Called(uuids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*plant.Plant), args.Error(1)
}

func (m *MockPlantGetter) GetPlantByName(name string) (*plant.Plant, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*plant.Plant), args.Error(1)
}

type LatexLikeParserTestSuite struct {
	suite.Suite
	testID1    uuid.UUID
	testID2    uuid.UUID
	testPlant  *plant.Plant
	testSpec   *plant.ConiferousSpecification
	mockGetter *MockPlantGetter
}

func (s *LatexLikeParserTestSuite) BeforeEach(t provider.T) {
	t.Epic("Post Parsing")
	t.Feature("LaTeX-like Plant Parser")

	s.testID1 = uuid.New()
	s.testID2 = uuid.New()

	var err error
	s.testSpec, err = plant.NewConiferousSpecification(
		1,
		1,
		10,
		plant.DryMoisture,
		plant.Light,
		plant.HeavySoil,
		9,
	)
	require.NoError(t, err)

	s.testPlant, err = plant.CreatePlant(
		s.testID2,
		"rose",
		"rose plant",
		"rose plant",
		uuid.New(),
		plant.PlantPhotos{},
		plant.ConiferousCategory,
		s.testSpec,
		time.Now(),
		time.Now(),
	)
	require.NoError(t, err)

	s.mockGetter = new(MockPlantGetter)
}

func (s *LatexLikeParserTestSuite) TestNoPlants(t provider.T) {
	t.Tags("positive", "parsing")
	t.Description("Test parsing text without plant references")

	text := "just some text without plants"
	expectedText := "just some text without plants"
	expectedPlantIDs := []uuid.UUID{}

	parser := parser.NewLatexLikePlantParser(s.mockGetter)
	plantIDs, resultText, err := parser.Parse(text)

	t.WithNewStep("Verify parsing results", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		require.Equal(t, expectedText, resultText)
		require.Equal(t, expectedPlantIDs, plantIDs)
	})
}

func (s *LatexLikeParserTestSuite) TestSinglePlantByUUID(t provider.T) {
	t.Tags("positive", "parsing")
	t.Description("Test parsing single plant reference by UUID")

	text := "text with \\plant{" + s.testID1.String() + "} plant"
	expectedText := "text with \\plant{" + s.testID1.String() + "} plant"
	expectedPlantIDs := []uuid.UUID{s.testID1}

	s.mockGetter.On("GetPlants", []uuid.UUID{s.testID1}).Return([]*plant.Plant{s.testPlant}, nil)

	parser := parser.NewLatexLikePlantParser(s.mockGetter)
	plantIDs, resultText, err := parser.Parse(text)

	t.WithNewStep("Verify parsing results", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		require.Equal(t, expectedText, resultText)
		require.Equal(t, expectedPlantIDs, plantIDs)
		s.mockGetter.AssertExpectations(t)
	})
}

func (s *LatexLikeParserTestSuite) TestMultiplePlantsByUUID(t provider.T) {
	t.Tags("positive", "parsing")
	t.Description("Test parsing multiple plant references by UUID")

	text := "\\plant{" + s.testID1.String() + "} and \\plant{" + s.testID2.String() + "}"
	expectedText := "\\plant{" + s.testID1.String() + "} and \\plant{" + s.testID2.String() + "}"
	expectedPlantIDs := []uuid.UUID{s.testID1, s.testID2}

	s.mockGetter.On("GetPlants", []uuid.UUID{s.testID1, s.testID2}).Return([]*plant.Plant{s.testPlant, s.testPlant}, nil)

	parser := parser.NewLatexLikePlantParser(s.mockGetter)
	plantIDs, resultText, err := parser.Parse(text)

	t.WithNewStep("Verify parsing results", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		require.Equal(t, expectedText, resultText)
		require.Equal(t, expectedPlantIDs, plantIDs)
		s.mockGetter.AssertExpectations(t)
	})
}

func (s *LatexLikeParserTestSuite) TestPlantByName(t provider.T) {
	t.Tags("positive", "parsing")
	t.Description("Test parsing plant reference by name")

	text := "\\plant{rose}"
	expectedText := "\\plant{" + s.testID2.String() + "}"
	expectedPlantIDs := []uuid.UUID{s.testID2}

	s.mockGetter.On("GetPlantByName", "rose").Return(s.testPlant, nil)
	s.mockGetter.On("GetPlants", []uuid.UUID{s.testID2}).Return([]*plant.Plant{s.testPlant}, nil)

	parser := parser.NewLatexLikePlantParser(s.mockGetter)
	plantIDs, resultText, err := parser.Parse(text)

	t.WithNewStep("Verify parsing results", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		require.Equal(t, expectedText, resultText)
		require.Equal(t, expectedPlantIDs, plantIDs)
		s.mockGetter.AssertExpectations(t)
	})
}

func (s *LatexLikeParserTestSuite) TestMixedPlantsByUUIDAndName(t provider.T) {
	t.Tags("positive", "parsing")
	t.Description("Test parsing mixed plant references by UUID and name")

	text := "\\plant{" + s.testID1.String() + "} and \\plant{rose}"
	expectedText := "\\plant{" + s.testID1.String() + "} and \\plant{" + s.testID2.String() + "}"
	expectedPlantIDs := []uuid.UUID{s.testID1, s.testID2}

	s.mockGetter.On("GetPlantByName", "rose").Return(s.testPlant, nil)
	s.mockGetter.On("GetPlants", []uuid.UUID{s.testID1, s.testID2}).Return([]*plant.Plant{s.testPlant, s.testPlant}, nil)

	parser := parser.NewLatexLikePlantParser(s.mockGetter)
	plantIDs, resultText, err := parser.Parse(text)

	t.WithNewStep("Verify parsing results", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		require.Equal(t, expectedText, resultText)
		require.Equal(t, expectedPlantIDs, plantIDs)
		s.mockGetter.AssertExpectations(t)
	})
}

func (s *LatexLikeParserTestSuite) TestInvalidPlantUUID(t provider.T) {
	t.Tags("negative", "parsing")
	t.Description("Test parsing with invalid plant UUID")

	text := "\\plant{invalid-uuid}"

	s.mockGetter.On("GetPlantByName", "invalid-uuid").Return(nil, fmt.Errorf("not found"))

	parser := parser.NewLatexLikePlantParser(s.mockGetter)

	plantIDs, resultText, err := parser.Parse(text)

	t.WithNewStep("Verify parsing error", func(ctx provider.StepCtx) {
		require.Error(t, err)
		require.Empty(t, plantIDs)
		require.Empty(t, resultText)
		s.mockGetter.AssertExpectations(t)
	})
}

func (s *LatexLikeParserTestSuite) TestUnclosedPlantTag(t provider.T) {
	t.Tags("negative", "parsing")
	t.Description("Test parsing with unclosed plant tag")

	text := "\\plant{rose"

	parser := parser.NewLatexLikePlantParser(s.mockGetter)
	plantIDs, resultText, err := parser.Parse(text)

	t.WithNewStep("Verify parsing error", func(ctx provider.StepCtx) {
		require.Error(t, err)
		require.Empty(t, plantIDs)
		require.Empty(t, resultText)
	})
}

func (s *LatexLikeParserTestSuite) TestSuffix(t provider.T) {
	t.Tags("functionality", "metadata")
	t.Description("Test parser suffix method")

	parser := parser.NewLatexLikePlantParser(nil)

	t.WithNewStep("Verify suffix value", func(ctx provider.StepCtx) {
		require.Equal(t, "latex", parser.Suffix())
	})
}

func TestLatexLikeParser(t *testing.T) {
	suite.RunSuite(t, new(LatexLikeParserTestSuite))
}
