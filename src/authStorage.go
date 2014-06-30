package main

import (
	"github.com/RangelReale/osin"
)

type AuthStorage struct {
}

// GetClient loads the client by id (client_id)
func (a *AuthStorage) GetClient(id string) (c *osin.Client, err error) {
	return
}

// SaveAuthorize saves authorize data.
func (a *AuthStorage) SaveAuthorize(d *osin.AuthorizeData) (err error) {
	return
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (a *AuthStorage) LoadAuthorize(code string) (d *osin.AuthorizeData, err error) {
	return
}

// RemoveAuthorize revokes or deletes the authorization code.
func (a *AuthStorage) RemoveAuthorize(code string) (err error) {
	return
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (a *AuthStorage) SaveAccess(*osin.AccessData) (err error) {
	return
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (a *AuthStorage) LoadAccess(token string) (d *osin.AccessData, err error) {
	return
}

// RemoveAccess revokes or deletes an AccessData.
func (a *AuthStorage) RemoveAccess(token string) (err error) {
	return
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (a *AuthStorage) LoadRefresh(token string) (d *osin.AccessData, err error) {
	return
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (a *AuthStorage) RemoveRefresh(token string) (err error) {
	return
}
