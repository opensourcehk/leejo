package oauth2

import (
	"github.com/RangelReale/osin"
	"log"
	"time"
	"upper.io/db"
)

// database adapted struct
type apiClient struct {
	Id          string `db:"id"`
	Secret      string `db:"secret"`
	RedirectUri string `db:"redirect_uri"`
}

func (c *apiClient) ToOsin() osin.Client {
	return &osin.DefaultClient{
		Id:          c.Id,
		Secret:      c.Secret,
		RedirectUri: c.RedirectUri,
	}
}

// database adapted struct
type apiAuthData struct {
	Id          int    `db:"id,omitempty"`
	Code        string `db:"code"`
	UserId      int    `db:"user_id"`
	ClientId    string `db:"client_id"`
	Scope       string `db:"scope"`
	State       string `db:"state"`
	RedirectUri string `db:"redirect_uri"`
	Created     int64  `db:"created_timestamp"`
	Expired     int64  `db:"expired_timestamp"`
}

func (d *apiAuthData) FromOsin(od *osin.AuthorizeData) *apiAuthData {
	d.Code = od.Code
	d.ClientId = od.Client.GetId()
	d.Scope = od.Scope
	d.State = od.State
	d.RedirectUri = od.RedirectUri
	d.Created = od.CreatedAt.Unix()
	d.Expired = od.CreatedAt.Unix() + int64(od.ExpiresIn)
	return d
}

func (d *apiAuthData) ToOsin() (od *osin.AuthorizeData) {
	od = &osin.AuthorizeData{
		Code: d.Code,
		Client: &osin.DefaultClient{
			Id: d.ClientId,
		},
		Scope: d.Scope,
		State: d.State,
		RedirectUri: d.RedirectUri,
		ExpiresIn: int32(d.Expired - time.Now().Unix()),
		CreatedAt: time.Unix(d.Created, 0),
	}
	return
}

// storage struct to fulfill osin interface
type AuthStorage struct {
	Db db.Database
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

	var cs []apiClient
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
	dd := (&apiAuthData{}).FromOsin(d)
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

	ds := []apiAuthData{}
	res := ac.Find(db.Cond{
		"code":      code,
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
		"code":      code,
	})
	err = res.Remove()

	return
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (a *AuthStorage) SaveAccess(ad *osin.AccessData) (err error) {
	log.Printf("SaveAccess: %#v\n", ad)
	return
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (a *AuthStorage) LoadAccess(token string) (d *osin.AccessData, err error) {
	log.Printf("LoadAccess: %s\n", token)
	return
}

// RemoveAccess revokes or deletes an AccessData.
func (a *AuthStorage) RemoveAccess(token string) (err error) {
	log.Printf("RemoveAccess: %s\n", token)
	return
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (a *AuthStorage) LoadRefresh(token string) (d *osin.AccessData, err error) {
	log.Printf("LoadRefresh: %s\n", token)
	return
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (a *AuthStorage) RemoveRefresh(token string) (err error) {
	log.Printf("RemoveRefresh: %s\n", token)
	return
}
