package permit

import (
	"context"
	"errors"
	"net/http"

	"github.com/cloudadrd/go-common"
	"github.com/cloudadrd/go-common/code"
	"github.com/cloudadrd/go-common/log"

	"github.com/gin-gonic/gin"
)

const (
	_uid        = "uid"
	_user       = "user"
	_permission = "permission"
)

type Auth interface {
	Permission(ctx context.Context, uid string) ([]string, error)
	Session(ctx context.Context, req *http.Request) (string, interface{}, error)
}

type Config struct {
	Auth Auth
}

func (c *Config) validate() error {
	if c == nil {
		return errors.New("Permit: empty config ")
	}
	if c.Auth == nil {
		return errors.New("Permit: Auth is nil ")
	}
	return nil
}

// Permit is a Permit instance.
type Permit struct {
	a Auth
}

// New new a Permit service.
func New(c *Config) (s *Permit) {
	return &Permit{
		a: c.Auth,
	}
}

// Login check login
func (p *Permit) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := p.login(ctx)
		if err != nil {
			pkg.Json(ctx, nil, err)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (p *Permit) login(ctx *gin.Context) (string, error) {
	uid, user, err := p.a.Session(ctx, ctx.Request)
	if err != nil {
		return "", err
	}
	if user != nil {
		ctx.Set(_user, user)
	}
	if uid != "" {
		ctx.Set(_uid, uid)
	}
	return uid, nil
}

// Permit check permissions
func (p *Permit) Permit(permit string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := p.login(ctx)
		if err != nil {
			pkg.Json(ctx, nil, err)
			ctx.Abort()
			return
		}
		if err := p.permit(ctx, uid, permit); err != nil {
			pkg.Json(ctx, nil, err)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (p *Permit) permit(ctx *gin.Context, uid string, permit string) error {
	if permit == "" {
		return nil
	}
	permits, err := p.a.Permission(ctx, uid)
	if err != nil {
		log.Warnf("get (%v) permit list err", uid)
		return code.Denied
	}
	if len(permits) > 0 {
		ctx.Set(_permission, permits)
	}
	for _, p := range permits {
		if permit == p {
			return nil
		}
	}
	return code.Denied
}

// FromContextPermission get permission list form context
func FromContextPermission(ctx context.Context) ([]string, bool) {
	v, ok := ctx.Value(_permission).([]string)
	return v, ok
}

// FromContextUserInfo get user info form context
func FromContextUserInfo(ctx context.Context) interface{} {
	return ctx.Value(_user)
}

// FromContextUid get uid form context
func FromContextUid(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(_uid).(string)
	return uid, ok
}
