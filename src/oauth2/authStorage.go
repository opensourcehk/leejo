package oauth2

import (
	"github.com/RangelReale/osin"
	"github.com/gourd/service"
	"github.com/gourd/session"
	"leejo/data"
	"log"
)

// storage struct to fulfill osin interface
type AuthStorage struct {
	S  session.Session
	P  ServiceProvider
}

func (a *AuthStorage) SetSession(s session.Session) {
	a.S = s
}

func (a *AuthStorage) SetProvider(p ServiceProvider) {
	a.P = p
}

// clone the storage
func (a *AuthStorage) Clone() (c osin.Storage) {
	c = &AuthStorage{
		S:  a.S,
		P:  a.P,
	}
	return
}

func (a *AuthStorage) Close() {
}

// GetClient loads the client by id (client_id)
func (a *AuthStorage) GetClient(id string) (c osin.Client, err error) {
	cs := a.P.ClientService(a.S)

	// define entity list type
	el := cs.AllocEntityList()
	err = cs.Retrieve(id, nil, el)
	log.Printf("GetClient: %s: %#v\n", id, cs)
	if err != nil {
		panic(err)
	}

	// if there is result, pass it out
	if cs.Len(el) > 0 {
		// this assignment is clumcy and structure dependent
		// need to change
		cl := el.(*[]data.ApiClient)
		c = (*cl)[0].ToOsin()
		log.Printf("Client Obtained: %#v\n", c)
	}

	return
}

// SaveAuthorize saves authorize data.
func (a *AuthStorage) SaveAuthorize(d *osin.AuthorizeData) (err error) {
	log.Printf("SaveAuthorize: %#v\n", d)
	s := a.P.AuthService(a.S)
	dd := (&data.ApiAuthData{}).FromOsin(d)
	log.Printf("SaveAuthorize: %#v\n", dd)
	c := &service.BasicContext{}
	err = s.Create(c, dd)
	return
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (a *AuthStorage) LoadAuthorize(code string) (d *osin.AuthorizeData, err error) {

	log.Printf("LoadAuthorize: %s\n", code)
	s := a.P.AuthService(a.S)

	// allocate memory for variables
	el := s.AllocEntityList()

	// find all of entities with the code
	cond := service.NewConds().Add("code", code)
	s.Search(cond, el)

	// get authdata
	if s.Len(el) == 0 {
		return
	}

	l := el.(*[]data.ApiAuthData)
	d = (*l)[0].ToOsin()

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
	s := a.P.AuthService(a.S)

	// allocate memory for variables
	el := s.AllocEntityList()
	var al *[]data.ApiAuthData
	al = el.(*[]data.ApiAuthData)

	// search for AuthData
	cond := service.NewConds().Add("code", code)
	s.Search(cond, el)
	if s.Len(el) == 0 {
		return
	}

	// delete the AuthData
	err = s.Delete((*al)[0].Id, nil)
	return
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (a *AuthStorage) SaveAccess(ad *osin.AccessData) (err error) {
	log.Printf("SaveAccess: %#v\n", ad)
	s := a.P.AccessService(a.S)
	dd := (&data.ApiAccess{}).FromOsin(ad)
	log.Printf("SaveAccess: %#v\n", dd)
	c := &service.BasicContext{}
	err = s.Create(c, dd)
	return
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (a *AuthStorage) LoadAccess(token string) (d *osin.AccessData, err error) {

	log.Printf("LoadAccess: %s\n", token)
	s := a.P.AccessService(a.S)

	// allocate memory for variables
	el := s.AllocEntityList()

	// find all of entities with the code
	cond := service.NewConds().Add("access_token", token)
	s.Search(cond, el)

	// get access data
	if s.Len(el) == 0 {
		return
	}

	l := el.(*[]data.ApiAccess)
	d = (*l)[0].ToOsin()

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

	log.Printf("RemoveAuthorize: %s\n", token)
	s := a.P.AccessService(a.S)

	// allocate memory for variables
	el := s.AllocEntityList()
	var al *[]data.ApiAccess
	al = el.(*[]data.ApiAccess)

	// search for Access
	cond := service.NewConds().Add("access_token", token)
	s.Search(cond, el)
	if s.Len(el) == 0 {
		return
	}

	// delete the Access
	err = s.Delete((*al)[0].Id, nil)
	return
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (a *AuthStorage) LoadRefresh(token string) (d *osin.AccessData, err error) {

	log.Printf("LoadRefresh: %s\n", token)
	s := a.P.AccessService(a.S)

	// allocate memory for variables
	el := s.AllocEntityList()

	// find all of entities with the code
	cond := service.NewConds().Add("refresh_token", token)
	s.Search(cond, el)

	// get access data
	if s.Len(el) == 0 {
		return
	}

	l := el.(*[]data.ApiAccess)
	d = (*l)[0].ToOsin()

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
	s := a.P.AccessService(a.S)

	// allocate memory for variables
	el := s.AllocEntityList()
	var al *[]data.ApiAccess
	al = el.(*[]data.ApiAccess)

	// search for Access
	cond := service.NewConds().Add("refresh_token", token)
	s.Search(cond, el)
	if s.Len(el) == 0 {
		return
	}

	// delete the Access
	err = s.Delete((*al)[0].Id, nil)
	return
}
