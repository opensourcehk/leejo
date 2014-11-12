package data

import (
	"github.com/RangelReale/osin"
	"time"
)

// database adapted struct
type ApiClient struct {
	Id          string `db:"id"`
	Secret      string `db:"secret"`
	RedirectUri string `db:"redirect_uri"`
}

func (c *ApiClient) ToOsin() osin.Client {
	return &osin.DefaultClient{
		Id:          c.Id,
		Secret:      c.Secret,
		RedirectUri: c.RedirectUri,
	}
}

// User interface
type ApiUser interface {
	GetId() int64
}

// implementation for User interface
type ApiAuthUser struct {
	Id int64
}

func (u *ApiAuthUser) GetId() int64 {
	return u.Id
}

// database adapted struct
type ApiAuthData struct {
	Id          int64  `db:"id,omitempty"`
	Code        string `db:"code"`
	UserId      int64  `db:"user_id"`
	ClientId    string `db:"client_id"`
	Scope       string `db:"scope"`
	State       string `db:"state"`
	RedirectUri string `db:"redirect_uri"`
	Created     int64  `db:"created_timestamp"`
	Expired     int64  `db:"expired_timestamp"`
}

func (d *ApiAuthData) FromOsin(od *osin.AuthorizeData) *ApiAuthData {

	// attempt to case userdata
	if user, ok := od.UserData.(User); ok {
		d.UserId = user.GetId()
	}

	d.Code = od.Code
	d.ClientId = od.Client.GetId()
	d.Scope = od.Scope
	d.State = od.State
	d.RedirectUri = od.RedirectUri
	d.Created = od.CreatedAt.Unix()
	d.Expired = od.CreatedAt.Unix() + int64(od.ExpiresIn)
	return d
}

func (d *ApiAuthData) ToOsin() (od *osin.AuthorizeData) {
	od = &osin.AuthorizeData{
		Code: d.Code,
		Client: &osin.DefaultClient{
			Id: d.ClientId,
		},
		Scope:       d.Scope,
		State:       d.State,
		RedirectUri: d.RedirectUri,
		ExpiresIn:   int32(d.Expired - time.Now().Unix()),
		CreatedAt:   time.Unix(d.Created, 0),
		UserData: &ApiAuthUser{
			Id: d.UserId,
		},
	}
	return
}

// database adapted struct
type ApiAccess struct {
	Id           int64  `db:"id,omitempty"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
	ClientId     string `db:"client_id"`
	UserId       int64  `db:"user_id"`
	Scope        string `db:"scope"`
	Created      int64  `db:"created_timestamp"`
	Expired      int64  `db:"expired_timestamp"`
}

func (d *ApiAccess) FromOsin(od *osin.AccessData) *ApiAccess {

	// attempt to case userdata
	if user, ok := od.UserData.(User); ok {
		d.UserId = user.GetId()
	}

	d.AccessToken = od.AccessToken
	d.RefreshToken = od.RefreshToken
	d.ClientId = od.Client.GetId()
	d.Scope = od.Scope
	d.Created = od.CreatedAt.Unix()
	d.Expired = od.CreatedAt.Unix() + int64(od.ExpiresIn)
	return d
}

func (d *ApiAccess) ToOsin() (od *osin.AccessData) {
	od = &osin.AccessData{
		AccessToken:  d.AccessToken,
		RefreshToken: d.RefreshToken,
		Client: &osin.DefaultClient{
			Id: d.ClientId,
		},
		Scope:     d.Scope,
		ExpiresIn: int32(d.Expired - time.Now().Unix()),
		CreatedAt: time.Unix(d.Created, 0),
		UserData: &ApiAuthUser{
			Id: d.UserId,
		},
	}
	return
}
