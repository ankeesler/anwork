// Package auth provides an authentication mechanism for the ANWORK project.
//
// This authentication mechanism is for single-occupancy use, i.e., only
// one identity exists and it is tied to an RSA key/32-byte secret pair. It
// generates RSA public key encrypted tokens that can only be consumed by
// those who have access to the matching RSA private key.
//
// The package provides two main types: Server and Client. The Server object
// provides the ability to generate encrypted tokens (Token()) and validate
// decrypted tokens (Authenticate()). The Client provides the ability to
// validate encrypted tokens (Authenticate()).
package auth
