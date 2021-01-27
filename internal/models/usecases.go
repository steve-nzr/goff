package models

import (
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/pkg/abstract"
)

// UseCaseResponse holds the response from a game use-case.
type UseCaseResponse struct {
	// Direct response to the client that called the method.
	ResponseToCaller abstract.Serializable

	// Responses to other clients (like around players...).
	ResponsesToOthers map[customtypes.ID]abstract.Serializable
}
