package parser

import (
	"PlantSite/internal/models/post"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	LatexLikePlantParserType = "latex"
)

type LatexLikePlantParser struct {
	plantGetter PlantGetter
}

func NewLatexLikePlantParser(plantGetter PlantGetter) *LatexLikePlantParser {
	return &LatexLikePlantParser{
		plantGetter: plantGetter,
	}
}

// Пример входного текста:
//
//	....\plant{id1}... \plant{id2}... \plant{id3}...
func (p *LatexLikePlantParser) Parse(text string) ([]uuid.UUID, string, error) {
	plantIDs := make([]uuid.UUID, 0)
	const tagStart = "\\plant{"
	var result strings.Builder
	lastPos := 0

	for {
		startIdx := strings.Index(text[lastPos:], tagStart)
		if startIdx == -1 {
			break
		}
		startIdx += lastPos

		result.WriteString(text[lastPos:startIdx])

		openBraceIdx := startIdx + len(tagStart)
		closeBraceIdx := strings.Index(text[openBraceIdx:], "}")
		if closeBraceIdx == -1 {
			return nil, "", fmt.Errorf("%w: } not found", post.ErrContentParsingError)
		}
		closeBraceIdx += openBraceIdx

		content := strings.TrimSpace(text[openBraceIdx:closeBraceIdx])

		var plantID uuid.UUID
		if id, err := uuid.Parse(content); err == nil {
			plantID = id
		} else {
			plant, err := p.plantGetter.GetPlantByName(content)
			if err != nil {
				return nil, "", fmt.Errorf("%w: %v", post.ErrContentParsingError, err)
			}
			plantID = plant.ID()
			content = plantID.String() // Change name to UUID
		}

		result.WriteString(tagStart)
		result.WriteString(content)
		result.WriteString("}")

		plantIDs = append(plantIDs, plantID)
		lastPos = closeBraceIdx + 1
	}

	result.WriteString(text[lastPos:])

	// Check plant existence
	if len(plantIDs) > 0 {
		_, err := p.plantGetter.GetPlants(plantIDs)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %v", post.ErrContentParsingError, err)
		}
	}

	return plantIDs, result.String(), nil
}

func (p *LatexLikePlantParser) Suffix() string {
	return "latex"
}
