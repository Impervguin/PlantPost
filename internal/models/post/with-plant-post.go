package post

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func CheckContentWithPlant(content *Content) bool {
	return strings.HasPrefix(string(content.ContentType), string(ContentTypeWithPlant))
}

func PlantParserType(content *Content) string {
	if !CheckContentWithPlant(content) {
		return ""
	}
	return string(content.ContentType)[len(string(ContentTypeWithPlant))+1:]
}

func WithPlantContentType(parserType string) ContentFormat {
	return ContentFormat(fmt.Sprintf("%s_%s", ContentTypeWithPlant, parserType))
}

const (
	ContentTypeWithPlant ContentFormat = "with_plant"
)

type ContentPlantParser interface {
	// Returns plant IDs and text unified text
	Parse(string) ([]uuid.UUID, string, error)
	Suffix() string
}

type ContentWithPlant struct {
	Content
	plantIDs []uuid.UUID
	parser   ContentPlantParser
}

func NewContentWithPlant(text string, format ContentFormat, plantParser ContentPlantParser) (*ContentWithPlant, error) {
	if err := format.Validate(); err != nil {
		return nil, err
	}

	plantIDs, text, err := plantParser.Parse(text)
	if err != nil {
		return nil, fmt.Errorf("plant parser failed: %v", err)
	}

	content := &ContentWithPlant{
		Content: Content{
			Text:        text,
			ContentType: format,
		},
		plantIDs: plantIDs,
	}

	if err := content.Validate(); err != nil {
		return nil, err
	}
	return content, nil
}

func (c *ContentWithPlant) PlantIDs() []uuid.UUID {
	return c.plantIDs
}

func (c *ContentWithPlant) Parser() ContentPlantParser {
	return c.parser
}

func (c *ContentWithPlant) UpdateContent(text string, plantParser ContentPlantParser) error {
	plantIDs, text, err := plantParser.Parse(text)
	if err != nil {
		return fmt.Errorf("plant parser failed: %v", err)
	}
	c.plantIDs = plantIDs
	c.Text = text
	c.parser = plantParser
	return nil
}
