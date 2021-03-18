package model

import (
	"github.com/go-openapi/strfmt"
	"github.com/goharbor/harbor/src/pkg/user/models"
	svrmodels "github.com/goharbor/harbor/src/server/v2.0/models"
)

// User ...
type User struct {
	*models.User
}

// ToSearchRespItem ...
func (u *User) ToSearchRespItem() *svrmodels.UserSearchRespItem {
	return &svrmodels.UserSearchRespItem{
		UserID:   int64(u.UserID),
		Username: u.Username,
	}
}

// ToUserProfile ...
func (u *User) ToUserProfile() *svrmodels.UserProfile {
	return &svrmodels.UserProfile{
		Email:    u.Email,
		Realname: u.Realname,
		Comment:  u.Comment,
	}
}

// ToUserCommon ...
func (u *User) ToUserCommon() *svrmodels.UserCommon {
	return &svrmodels.UserCommon{
		UserProfile: *u.ToUserProfile(),
		Username:    u.Username,
	}
}

// ToUserResp ...
func (u *User) ToUserResp() *svrmodels.UserResp {
	res := &svrmodels.UserResp{
		UserCommon:      *u.ToUserCommon(),
		UserID:          int64(u.UserID),
		SysadminFlag:    u.SysAdminFlag,
		AdminRoleInAuth: u.AdminRoleInAuth,
		CreationTime:    strfmt.DateTime(u.CreationTime),
		UpdateTime:      strfmt.DateTime(u.UpdateTime),
	}
	if u.OIDCUserMeta != nil {
		res.OidcUserMeta = &svrmodels.OIDCUserInfo{
			ID:           u.OIDCUserMeta.ID,
			UserID:       int64(u.OIDCUserMeta.UserID),
			Subiss:       u.OIDCUserMeta.SubIss,
			Secret:       u.OIDCUserMeta.PlainSecret,
			CreationTime: strfmt.DateTime(u.OIDCUserMeta.CreationTime),
			UpdateTime:   strfmt.DateTime(u.OIDCUserMeta.UpdateTime),
		}
	}
	return res
}
