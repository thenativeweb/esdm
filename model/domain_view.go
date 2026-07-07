package model

// DomainView is the typed view over an ESDM document
// whose kind is "domain". Domain documents carry no
// kind-specific fields beyond the common ones provided
// by DocumentViewBase.
type DomainView struct {
	DocumentViewBase
}
