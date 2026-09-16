package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/schema"
)

func TestIsLanguageTag(t *testing.T) {
	t.Run("uses the same pattern as the core schema's languageTag definition", func(t *testing.T) {
		var parsed struct {
			Defs struct {
				LanguageTag struct {
					Pattern string `yaml:"pattern"`
				} `yaml:"languageTag"`
			} `yaml:"$defs"`
		}
		require.NoError(t, yaml.Unmarshal(schema.Core(), &parsed))
		assert.Equal(t, parsed.Defs.LanguageTag.Pattern, model.LanguageTagPattern.String())
	})

	t.Run("accepts BCP 47 tags and rejects other strings", func(t *testing.T) {
		cases := []struct {
			value string
			isTag bool
		}{
			{"en", true},
			{"de", true},
			{"de-AT", true},
			{"zh-Hant", true},
			{"ast", true},
			{"German", false},
			{"DE", false},
			{"d", false},
			{"de-", false},
			{"", false},
		}
		for _, c := range cases {
			t.Run(c.value, func(t *testing.T) {
				assert.Equal(t, c.isTag, model.IsLanguageTag(c.value))
			})
		}
	})
}
