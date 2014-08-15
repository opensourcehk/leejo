package oauth2

import (
	"github.com/RangelReale/osin"
	"time"
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

// User interface
type User interface {
	GetId() int64
}

// implementation for User interface
type apiAuthUser struct {
	Id int64
}

func (u *apiAuthUser) GetId() int64 {
	return u.Id
}

// database adapted struct
type apiAuthData struct {
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

func (d *apiAuthData) FromOsin(od *osin.AuthorizeData) *apiAuthData {

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

func (d *apiAuthData) ToOsin() (od *osin.AuthorizeData) {
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
		UserData: &apiAuthUser{
			Id: d.UserId,
		},
	}
	return
}

// database adapted struct
type apiAccess struct {
	Id           int64  `db:"id,omitempty"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
	ClientId     string `db:"client_id"`
	UserId       int64  `db:"user_id"`
	Scope        string `db:"scope"`
	Created      int64  `db:"created_timestamp"`
	Expired      int64  `db:"expired_timestamp"`
}

func (d *apiAccess) FromOsin(od *osin.AccessData) *apiAccess {

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

func (d *apiAccess) ToOsin() (od *osin.AccessData) {
	od = &osin.AccessData{
		AccessToken:  d.AccessToken,
		RefreshToken: d.RefreshToken,
		Client: &osin.DefaultClient{
			Id: d.ClientId,
		},
		Scope:     d.Scope,
		ExpiresIn: int32(d.Expired - time.Now().Unix()),
		CreatedAt: time.Unix(d.Created, 0),
		UserData: &apiAuthUser{
			Id: d.UserId,
		},
	}
	return
}
