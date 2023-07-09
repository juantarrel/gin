package gin

import (
	"github.com/gin-gonic/gin"

	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/contracts/validation"
)

const HttpBinding = "goravel.http"
const RouteBinding = "goravel.route"

var App foundation.Application

var (
	ConfigFacade      config.Config
	CacheFacade       cache.Cache
	LogFacade         log.Log
	RateLimiterFacade http.RateLimiter
	ValidationFacade  validation.Validation
)

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(HttpBinding, func(app foundation.Application) (any, error) {
		return NewGinContext(&gin.Context{}), nil
	})
	app.Bind(RouteBinding, func(app foundation.Application) (any, error) {
		return NewGinRoute(app.MakeConfig()), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	ConfigFacade = app.MakeConfig()
	CacheFacade = app.MakeCache()
	LogFacade = app.MakeLog()
	ValidationFacade = app.MakeValidation()
}
