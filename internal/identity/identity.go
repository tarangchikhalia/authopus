// Package identity defines identity references. References do not imply authentication.
package identity

type Kind string

const (
	Agent       Kind = "agent"
	Integration Kind = "integration"
)

// Reference identifies a subject within an issuer namespace.
type Reference struct {
	Issuer  string
	Subject string
}
