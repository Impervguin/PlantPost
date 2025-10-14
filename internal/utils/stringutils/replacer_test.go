//go:build unit

package stringutils

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type StringUtilsTestSuite struct {
	suite.Suite
}

func (s *StringUtilsTestSuite) BeforeEach(t provider.T) {
	t.Epic("String Utils")
	t.Feature("String Manipulation")
}

func (s *StringUtilsTestSuite) TestReplaceFunc(t provider.T) {
	t.Tags("replace", "function")
	t.Description("Test ReplaceFunc functionality with various patterns and replacers")

	t.Run("Basic replacement with single match", func(t provider.T) {
		t.Parallel()

		input := "Hello hello world"
		pattern := "h%so"
		replacer := func(match string) string {
			return "beautiful"
		}

		result := ReplaceFunc(input, pattern, replacer)
		expected := "Hello beautiful world"

		t.WithNewStep("Verify basic replacement", func(pctx provider.StepCtx) {
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
		})
	})

	t.Run("Multiple matches replacement", func(t provider.T) {
		t.Parallel()

		input := "Hello hello, welcome to hello world"
		pattern := "h%so"
		replacer := func(match string) string {
			return "beautiful"
		}

		result := ReplaceFunc(input, pattern, replacer)
		expected := "Hello beautiful, welcome to beautiful world"

		t.WithNewStep("Verify multiple replacements", func(pctx provider.StepCtx) {
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
		})
	})

	t.Run("Custom pattern with prefix and suffix", func(t provider.T) {
		t.Parallel()

		input := "Hello {{name}}, your age is {{age}}"
		pattern := "{{%s}}"
		replacer := func(match string) string {
			switch match {
			case "name":
				return "Alice"
			case "age":
				return "25"
			default:
				return match
			}
		}

		result := ReplaceFunc(input, pattern, replacer)
		expected := "Hello Alice, your age is 25"

		t.WithNewStep("Verify custom pattern replacement", func(pctx provider.StepCtx) {
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
		})
	})

	t.Run("Empty matches handling", func(t provider.T) {
		t.Parallel()

		input := "Test [] empty [] match"
		pattern := "[%s]"
		replacer := func(match string) string {
			if match == "" {
				return "EMPTY"
			}
			return match
		}

		result := ReplaceFunc(input, pattern, replacer)
		expected := "Test EMPTY empty EMPTY match"

		t.WithNewStep("Verify empty matches handling", func(pctx provider.StepCtx) {
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
		})
	})

	t.Run("No matches found", func(t provider.T) {
		t.Parallel()

		input := "Hello world"
		pattern := "h%so"
		replacer := func(match string) string {
			return "replaced"
		}

		result := ReplaceFunc(input, pattern, replacer)
		expected := "Hello world"

		t.WithNewStep("Verify no replacement when no matches", func(pctx provider.StepCtx) {
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
		})
	})

	t.Run("Invalid pattern without %s", func(t provider.T) {
		t.Parallel()

		input := "Hello world"
		pattern := "invalid"
		replacer := func(match string) string {
			return "replaced"
		}

		result := ReplaceFunc(input, pattern, replacer)
		expected := "Hello world"

		t.WithNewStep("Verify invalid pattern returns original string", func(pctx provider.StepCtx) {
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
		})
	})

	t.Run("Pattern with multiple %s placeholders", func(t provider.T) {
		t.Parallel()

		input := "Hello hello world"
		pattern := "h%so%s" // Invalid pattern - should return original
		replacer := func(match string) string {
			return "replaced"
		}

		result := ReplaceFunc(input, pattern, replacer)
		expected := "Hello hello world"

		t.WithNewStep("Verify pattern with multiple placeholders returns original", func(pctx provider.StepCtx) {
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
		})
	})

	t.Run("Empty input string", func(t provider.T) {
		t.Parallel()

		input := ""
		pattern := "h%so"
		replacer := func(match string) string {
			return "replaced"
		}

		result := ReplaceFunc(input, pattern, replacer)
		expected := ""

		t.WithNewStep("Verify empty input handling", func(pctx provider.StepCtx) {
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
		})
	})

	t.Run("Replacer returns empty string", func(t provider.T) {
		t.Parallel()

		input := "Hello hello world"
		pattern := "h%so"
		replacer := func(match string) string {
			return ""
		}

		result := ReplaceFunc(input, pattern, replacer)
		expected := "Hello  world"

		t.WithNewStep("Verify empty replacement handling", func(pctx provider.StepCtx) {
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
		})
	})
}

func TestStringUtils(t *testing.T) {
	suite.RunSuite(t, new(StringUtilsTestSuite))
}
