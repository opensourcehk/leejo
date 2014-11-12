package oauth2

import (
	"github.com/RangelReale/osin"
	"github.com/gourd/session"
	"leejo/data"
	"log"
	"upper.io/db"
)

// storage struct to fulfill osin interface
type AuthStorage struct {
	S       session.Session
	ClientP ServiceProvider
	AuthP   ServiceProvider
	Db      db.Database
}

func (a *AuthStorage) Clone() osin.Storage {
	return a
}

func (a *AuthStorage) Close() {
}

// GetClient loads the client by id (client_id)
func (a *AuthStorage) GetClient(id string) (c osin.Client, err error) {
	cc, err := a.Db.Collection("leejo_api_client")
	res := cc.Find(db.Cond{"id": id})

	var cs []data.ApiClient
	err = res.All(&cs)
	log.Printf("GetClient: %s: %#v\n", id, cs)
	if err != nil {
		panic(err)
	}

	// if there is result, pass it out
	if len(cs) > 0 {
		c = cs[0].ToOsin()
		log.Printf("Client Obtained: %#v\n", c)
	}
	return
}

// SaveAuthorize saves authorize data.
func (a *AuthStorage) SaveAuthorize(d *osin.AuthorizeData) (err error) {
	log.Printf("SaveAuthorize: %#v\n", d)
	ac, err := a.Db.Collection("leejo_api_authdata")
	if err != nil {
		return
	}
	dd := (&data.ApiAuthData{}).FromOsin(d)
	log.Printf("SaveAuthorize: %#v\n", dd)
	_, err = ac.Append(dd)
	return
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (a *AuthStorage) LoadAuthorize(code string) (d *osin.AuthorizeData, err error) {
	log.Printf("LoadAuthorize: %s\n", code)
	ac, err := a.Db.Collection("leejo_api_authdata")
	if err != nil {
		return
	}

	ds := []data.ApiAuthData{}
	res := ac.Find(db.Cond{
		"code": code,
	})
	err = res.All(&ds)
	log.Printf("AuthData retrieved: %s: %#v\n", code, ds)
	if err != nil {
		return
	}

	// get authdata
	if len(ds) == 0 {
		return
	}
	d = ds[0].ToOsin()
	log.Printf("AuthData.ToOsin: %#v\n", d)

	// also load api client data
	c, err := a.GetClient(d.Client.GetId())
	if err != nil {
		return
	}
	d.Client = c

	// TODO: also load user data

	return
}

// RemoveAuthorize revokes or deletes the authorization code.
func (a *AuthStorage) RemoveAuthorize(code string) (err error) {
	log.Printf("RemoveAuthorize: %s\n", code)
	ac, err := a.Db.Collection("leejo_api_authdata")
	if err != nil {
		return
	}

	res := ac.Find(db.Cond{
		"code": code,
	})
	err = res.Remove()

	return
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (a *AuthStorage) SaveAccess(ad *osin.AccessData) (err error) {
	log.Printf("SaveAccess: %#v\n", ad)
	ac, err := a.Db.Collection("leejo_api_access")
	if err != nil {
		return
	}
	dd := (&data.ApiAccess{}).FromOsin(ad)
	log.Printf("SaveAuthorize adapted: %#v\n", dd)
	_, err = ac.Append(dd)
	return
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (a *AuthStorage) LoadAccess(token string) (d *osin.AccessData, err error) {
	log.Printf("LoadAccess: %s\n", token)
	ac, err := a.Db.Collection("leejo_api_access")
	if err != nil {
		return
	}

	ds := []data.ApiAccess{}
	res := ac.Find(db.Cond{
		"access_token": token,
	})
	err = res.All(&ds)
	log.Printf("Access retrieved: %s: %#v\n", token, ds)
	if err != nil {
		return
	}

	// get access
	if len(ds) == 0 {
		return
	}
	d = ds[0].ToOsin()
	log.Printf("Access.ToOsin: %#v\n", d)

	// also load api client data
	c, err := a.GetClient(d.Client.GetId())
	if err != nil {
		return
	}
	d.Client = c

	return
}

// RemoveAccess revokes or deletes an AccessData.
func (a *AuthStorage) RemoveAccess(token string) (err error) {
	log.Printf("RemoveAccess: %s\n", token)
	ac, err := a.Db.Collection("leejo_api_access")
	if err != nil {
		return
	}

	res := ac.Find(db.Cond{
		"access_token": token,
	})
	err = res.Remove()

	return
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (a *AuthStorage) LoadRefresh(token string) (d *osin.AccessData, err error) {
	log.Printf("LoadRefresh: %s\n", token)
	ac, err := a.Db.Collection("leejo_api_access")
	if err != nil {
		return
	}

	ds := []data.ApiAccess{}
	res := ac.Find(db.Cond{
		"refresh_token": token,
	})
	err = res.All(&ds)
	log.Printf("Access (Refresh) retrieved: %s: %#v\n", token, ds)
	if err != nil {
		return
	}

	// get access
	if len(ds) == 0 {
		return
	}
	d = ds[0].ToOsin()
	log.Printf("Access.ToOsin: %#v\n", d)

	// also load api client data
	c, err := a.GetClient(d.Client.GetId())
	if err != nil {
		return
	}
	d.Client = c

	return
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (a *AuthStorage) RemoveRefresh(token string) (err error) {
	log.Printf("RemoveRefresh: %s\n", token)
	ac, err := a.Db.Collection("leejo_api_access")
	if err != nil {
		return
	}

	res := ac.Find(db.Cond{
		"refresh_token": token,
	})
	err = res.Remove()

	return
}
