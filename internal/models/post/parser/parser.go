package parser

import (
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/post"
	"fmt"

	"github.com/google/uuid"
)

type PlantGetter interface {
	GetPlants(uuids []uuid.UUID) ([]*plant.Plant, error)
	GetPlantByName(name string) (*plant.Plant, error)
}

func GetParser(content *post.Content, plantGetter PlantGetter) (post.ContentPlantParser, error) {
	parser := post.PlantParserType(content)
	switch parser {
	case "latex":
		return NewLatexLikePlantParser(plantGetter), nil
	default:
		return nil, fmt.Errorf("unsupported plant parser: %s", parser)
	}

}
