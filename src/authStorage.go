package main

import (
	"github.com/RangelReale/osin"
	"upper.io/db"
)

// database adapted struct
type apiClient struct {
	Id          string `db:"id"`
	Secret      string `db:"secret"`
	RedirectUri string `db:"redirect_uri"`
}

func (c *apiClient) ToOsin() (oc *osin.Client) {
	oc = &osin.Client{
		Id:          c.Id,
		Secret:      c.Secret,
		RedirectUri: c.RedirectUri,
	}
	return
}

// database adapted struct
type apiAuthData struct {
	Code     string `db:"code"`
	UserId   int    `db:"user_id"`
	ClientId string `db:"client_id"`
	Created  int    `db:"created_timestamp"`
	Expired  int    `db:"expired_timestamp"`
}

func (d *apiAuthData) FromOsin(od *osin.AuthorizeData) *apiAuthData {
	d.Code = od.Code
	d.ClientId = od.Client.Id
	return d
}

// storage struct to fulfill osin interface
type AuthStorage struct {
	Db db.Database
}

// GetClient loads the client by id (client_id)
func (a *AuthStorage) GetClient(id string) (c *osin.Client, err error) {
	cc, err := a.Db.Collection("leejo_api_client")
	res := cc.Find(db.Cond{"id": id})

	var cs []apiClient
	err = res.All(&cs)
	if err != nil {
		panic(err)
	}

	// if there is result, pass it out
	if len(cs) > 0 {
		c = cs[0].ToOsin()
	}
	return
}

// SaveAuthorize saves authorize data.
func (a *AuthStorage) SaveAuthorize(d *osin.AuthorizeData) (err error) {
	ac, err := a.Db.Collection("leejo_api_authdata")
	if err != nil {
		return
	}
	dd := (&apiAuthData{}).FromOsin(d)
	_, err = ac.Append(dd)
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
