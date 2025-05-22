package parser_test

import (
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/post/parser"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestLatexLikeParser(t *testing.T) {
	testID1 := uuid.New()
	testID2 := uuid.New()

	testSpec, err := plant.NewConiferousSpecification(
		1,
		1,
		10,
		plant.DryMoisture,
		plant.Light,
		plant.HeavySoil,
		9,
	)
	require.NoError(t, err)

	testPlant, err := plant.CreatePlant(
		testID2,
		"rose",
		"rose plant",
		"rose plant",
		uuid.New(),
		plant.PlantPhotos{},
		plant.ConiferousCategory,
		testSpec,
		time.Now(),
		time.Now(),
	)
	require.NoError(t, err)

	tests := []struct {
		name             string
		text             string
		expectedText     string
		expectedPlantIDs []uuid.UUID
		mockSetup        func(*MockPlantGetter)
		expectError      bool
	}{
		{
			name:             "no plants",
			text:             "just some text without plants",
			expectedText:     "just some text without plants",
			expectedPlantIDs: []uuid.UUID{},
			mockSetup:        func(m *MockPlantGetter) {},
		},
		{
			name:             "single plant by UUID",
			text:             "text with \\plant{" + testID1.String() + "} plant",
			expectedText:     "text with \\plant{" + testID1.String() + "} plant",
			expectedPlantIDs: []uuid.UUID{testID1},
			mockSetup: func(m *MockPlantGetter) {
				m.On("GetPlants", []uuid.UUID{testID1}).Return([]*plant.Plant{testPlant}, nil)
			},
		},
		{
			name:             "multiple plants by UUID",
			text:             "\\plant{" + testID1.String() + "} and \\plant{" + testID2.String() + "}",
			expectedText:     "\\plant{" + testID1.String() + "} and \\plant{" + testID2.String() + "}",
			expectedPlantIDs: []uuid.UUID{testID1, testID2},
			mockSetup: func(m *MockPlantGetter) {
				m.On("GetPlants", []uuid.UUID{testID1, testID2}).Return([]*plant.Plant{testPlant, testPlant}, nil)
			},
		},
		{
			name:             "plant by name",
			text:             "\\plant{rose}",
			expectedText:     "\\plant{" + testID2.String() + "}",
			expectedPlantIDs: []uuid.UUID{testID2},
			mockSetup: func(m *MockPlantGetter) {
				m.On("GetPlantByName", "rose").Return(testPlant, nil)
				m.On("GetPlants", []uuid.UUID{testID2}).Return([]*plant.Plant{testPlant}, nil)
			},
		},
		{
			name:             "mixed plants by UUID and name",
			text:             "\\plant{" + testID1.String() + "} and \\plant{rose}",
			expectedText:     "\\plant{" + testID1.String() + "} and \\plant{" + testID2.String() + "}",
			expectedPlantIDs: []uuid.UUID{testID1, testID2},
			mockSetup: func(m *MockPlantGetter) {
				m.On("GetPlantByName", "rose").Return(testPlant, nil)
				m.On("GetPlants", []uuid.UUID{testID1, testID2}).Return([]*plant.Plant{testPlant, testPlant}, nil)
			},
		},
		{
			name:             "invalid plant UUID",
			text:             "\\plant{invalid-uuid}",
			expectedPlantIDs: []uuid.UUID{},
			mockSetup: func(m *MockPlantGetter) {
				m.On("GetPlantByName", "invalid-uuid").Return(nil, fmt.Errorf("not found"))
			},
			expectError: true,
		},
		{
			name:        "unclosed plant tag",
			text:        "\\plant{rose",
			expectError: true,
			mockSetup:   func(m *MockPlantGetter) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := new(MockPlantGetter)
			tt.mockSetup(mockGetter)

			parser := parser.NewLatexLikePlantParser(mockGetter)
			plantIDs, resultText, err := parser.Parse(tt.text)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedText, resultText)
			require.Equal(t, tt.expectedPlantIDs, plantIDs)
			mockGetter.AssertExpectations(t)
		})
	}
}

func TestLatexLikeParser_Suffix(t *testing.T) {
	parser := parser.NewLatexLikePlantParser(nil)
	require.Equal(t, "latex", parser.Suffix())
}
