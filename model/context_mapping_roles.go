package model

// MappingRoles returns the two role fields of an asymmetric
// context-mapping type, in the order the schema lists them,
// and false for the symmetric types, which use `participants`
// instead of roles. It is the one place that knows which
// field names belong to which mapping type, for the resolver
// that checks term pairs and the renderers that name the
// endpoints.
func MappingRoles(mappingType string) ([2]string, bool) {
	roles, ok := mappingRoles[mappingType]
	return roles, ok
}

var mappingRoles = map[string][2]string{
	"customer-supplier":     {"customer", "supplier"},
	"conformist":            {"conformist", "upstream"},
	"anti-corruption-layer": {"downstream", "upstream"},
	"open-host-service":     {"host", "consumer"},
	"published-language":    {"publisher", "consumer"},
}
