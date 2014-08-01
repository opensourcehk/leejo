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
