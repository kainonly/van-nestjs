package index

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/huandu/xstrings"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"github.com/weplanx/utils/csrf"
	"github.com/weplanx/utils/passlib"
	"net/http"
	"reflect"
	"time"
)

type Controller struct {
	IndexService *Service
	Csrf         *csrf.Csrf
	Values       *common.Values
}

// Ping
// @router / [GET]
func (x *Controller) Ping(ctx context.Context, c *app.RequestContext) {
	x.Csrf.SetToken(c)
	c.JSON(http.StatusOK, utils.H{
		"ip":   c.ClientIP(),
		"time": time.Now(),
	})
}

type LoginDto struct {
	Email    string `json:"email,required" vd:"email($)"`
	Password string `json:"password,required" vd:"len($)>=8"`
}

// Login
// @router /login [POST]
func (x *Controller) Login(ctx context.Context, c *app.RequestContext) {
	var dto LoginDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	var metadata model.LoginMetadata
	metadata.Email = dto.Email
	metadata.Channel = "email"
	ts, err := x.IndexService.Login(ctx, dto.Email, dto.Password, &metadata)
	if err != nil {
		c.Error(err)
		return
	}

	metadata.Ip = c.ClientIP()
	var data model.LoginData

	data.UserAgent = string(c.UserAgent())
	go x.IndexService.WriteLoginLog(ctx, metadata, data)

	c.SetCookie("access_token", ts, 0, "/", "", protocol.CookieSameSiteLaxMode, true, true)
	c.JSON(200, utils.H{
		"code":    0,
		"message": "ok",
	})
}

// Verify
// @router /verify [GET]
func (x *Controller) Verify(ctx context.Context, c *app.RequestContext) {
	ts := c.Cookie("access_token")
	if ts == nil {
		c.JSON(401, utils.H{
			"code":    0,
			"message": MsgAuthenticationExpired,
		})
		return
	}

	if _, err := x.IndexService.Verify(ctx, string(ts)); err != nil {
		c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteLaxMode, true, true)
		c.JSON(401, utils.H{
			"code":    0,
			"message": MsgAuthenticationExpired,
		})
		return
	}

	c.JSON(200, utils.H{
		"code":    0,
		"message": "ok",
	})
}

// GetRefreshCode
// @router /code [GET]
func (x *Controller) GetRefreshCode(ctx context.Context, c *app.RequestContext) {
	claims := common.GetClaims(c)
	code, err := x.IndexService.GetRefreshCode(ctx, claims.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code": code,
	})
}

type RefreshTokenDto struct {
	Code string `json:"code,required"`
}

// RefreshToken
// @router /refresh_token [POST]
func (x *Controller) RefreshToken(ctx context.Context, c *app.RequestContext) {
	var dto RefreshTokenDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	claims := common.GetClaims(c)
	ts, err := x.IndexService.RefreshToken(ctx, claims, dto.Code)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("access_token", ts, 0, "/", "", protocol.CookieSameSiteLaxMode, true, true)
	c.JSON(http.StatusOK, utils.H{
		"code":    0,
		"message": "ok",
	})
}

// Logout
// @router /logout [POST]
func (x *Controller) Logout(ctx context.Context, c *app.RequestContext) {
	claims := common.GetClaims(c)
	if err := x.IndexService.Logout(ctx, claims.UserId); err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteLaxMode, true, true)
	c.JSON(http.StatusOK, utils.H{
		"code":    0,
		"message": "ok",
	})
}

// GetUser
// @router /user [GET]
func (x *Controller) GetUser(ctx context.Context, c *app.RequestContext) {
	claims := common.GetClaims(c)
	data, err := x.IndexService.GetUser(ctx, claims.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

type SetUserDto struct {
	Set         string `json:"$set,required" vd:"in($, 'email', 'name', 'avatar', 'password', 'backup_email')"`
	Email       string `json:"email,omitempty" vd:"(Set)$!='Email' || email($);msg:'must be email'"`
	BackupEmail string `json:"backup_email,omitempty" vd:"(Set)$!='BackupEmail' || email($);msg:'must be email'"`
	Name        string `json:"name,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Password    string `json:"password,omitempty" vd:"(Set)$!='Password' || len($)>8;msg:'must be greater than 8 characters'"`
}

// SetUser
// @router /user [POST]
func (x *Controller) SetUser(ctx context.Context, c *app.RequestContext) {
	var dto SetUserDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	data := make(map[string]interface{})
	path := xstrings.ToCamelCase(dto.Set)
	value := reflect.ValueOf(dto).FieldByName(path).Interface()
	if dto.Set == "password" {
		data[dto.Set], _ = passlib.Hash(value.(string))
	} else {
		data[dto.Set] = value
	}

	claims := common.GetClaims(c)
	_, err := x.IndexService.SetUser(ctx, claims.UserId, data)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code":    0,
		"message": "ok",
	})
}

type OptionsDto struct {
	Type string `query:"type"`
}

// Options 返回通用配置
func (x *Controller) Options(ctx context.Context, c *app.RequestContext) {
	var dto OptionsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	switch dto.Type {
	case "upload":
		switch x.Values.Cloud {
		case "tencent":
			c.JSON(http.StatusOK, utils.H{
				"type": "cos",
				"url": fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`,
					x.Values.TencentCosBucket, x.Values.TencentCosRegion,
				),
				"limit": x.Values.TencentCosLimit,
			})
			return
		}
	case "office":
		switch x.Values.Office {
		case "feishu":
			c.JSON(http.StatusOK, utils.H{
				"url":      "https://open.feishu.cn/open-apis/authen/v1/index",
				"redirect": x.Values.RedirectUrl,
				"app_id":   x.Values.FeishuAppId,
			})
			return
		}
	}
	return
}
