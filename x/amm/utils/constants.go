package utils

import "reflect"

// ExtraAccountTypes is a map of extra account types that can be overridden.
// This is defined as a global variable so it can be modified in the chain's app.go and used here without
// having to import the chain. Specifically, this is used for compatibility with Osmosis' Cosmos SDK fork
var ExtraAccountTypes map[reflect.Type]struct{}
