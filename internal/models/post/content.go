package post

import (
	"fmt"
	"strings"
)

type ContentFormat string

const (
	ContentTypePlainText ContentFormat = "plain_text"
)

func (c *ContentFormat) Validate() error {
	if *c == ContentTypePlainText {
		return nil
	} else if strings.HasPrefix(string(*c), string(ContentTypeWithPlant)) {
		return nil
	}
	return fmt.Errorf("invalid content format: %s", *c)
}

type Content struct {
	Text        string
	ContentType ContentFormat
}

func NewContent(text string, format ContentFormat) (*Content, error) {
	if err := format.Validate(); err != nil {
		return nil, err
	}
	content := &Content{
		Text:        text,
		ContentType: format,
	}
	if err := content.Validate(); err != nil {
		return nil, err
	}
	return content, nil
}

func (c *Content) Validate() error {
	if err := c.ContentType.Validate(); err != nil {
		return err
	}
	if c.Text == "" {
		return fmt.Errorf("content text cannot be empty")
	}

	// Проверки для будущих типов, например html, markdown, etc.
	if c.ContentType == ContentTypePlainText {
		return nil
	}
	if CheckContentWithPlant(c) {
		return nil
	}
	return fmt.Errorf("unsupported content type: %s", c.ContentType)
}
