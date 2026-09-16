package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// translatedTerminology is a bounded context written in
// English with one German translation. The tests below
// derive their failing variants from it.
const translatedTerminology = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: billing
scope:
  domain: d
language: en
ubiquitousLanguage:
  - term: Invoice
    definition: A request for payment issued to a customer.
    translations:
      - language: de
        term: Rechnung
        definition: Eine an den Kunden gerichtete Zahlungsaufforderung.
`

const duplicateTranslationLanguage = translatedTerminology + `      - language: de
        term: Faktura
        definition: Zweite deutsche Fassung.
`

const translationIntoOwnLanguage = translatedTerminology + `      - language: en
        term: Bill
        definition: The same concept, in the context's own language.
`

func TestDuplicateTranslationLanguage(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/duplicate-translation-language")

	t.Run("does not throw when every translation of a term has its own language", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, translatedTerminology)))
	})

	t.Run("throws when a term has two translations into the same language", func(t *testing.T) {
		diags := runRule(t, rule, buildModel(t, duplicateTranslationLanguage))
		require.Len(t, diags, 1)
		assert.Equal(t, `term "Invoice" has more than one translation into "de"`, diags[0].Message)
	})
}

func TestTranslationIntoOwnLanguage(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/translation-into-own-language")

	t.Run("does not throw when translations use other languages than the context", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, translatedTerminology)))
	})

	t.Run("throws when a term is translated into the language of its own bounded context", func(t *testing.T) {
		diags := runRule(t, rule, buildModel(t, translationIntoOwnLanguage))
		require.Len(t, diags, 1)
		assert.Equal(t, `term "Invoice" is translated into "en", which is the language of bounded context "billing"`, diags[0].Message)
	})
}
