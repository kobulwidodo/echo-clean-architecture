package rest

import (
	"go-clean/docs/swagger"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/configreader"
	"go-clean/src/lib/log"
	"go-clean/src/utils/config"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var once = &sync.Once{}

type REST interface {
	Run()
}

type rest struct {
	http         *echo.Echo
	configreader configreader.Interface
	uc           *usecase.Usecase
	auth         auth.Interface
	conf         config.ApplicationMeta
	log          log.Interface
}

func Init(confReader configreader.Interface, uc *usecase.Usecase, auth auth.Interface, conf config.ApplicationMeta, log log.Interface) REST {
	r := &rest{}
	once.Do(func() {
		e := echo.New()

		r = &rest{
			configreader: confReader,
			http:         e,
			uc:           uc,
			auth:         auth,
			conf:         conf,
			log:          log,
		}

		r.http.Use(middleware.CORS())
		r.http.Use(middleware.Logger())

		r.http.Validator = &entity.CustomValidator{Validator: validator.New()}

		r.Register()
	})

	return r
}

func (r *rest) Run() {
	r.http.Logger.Fatal(r.http.Start(":8080"))
}

func (r *rest) RegisterSwagger() {
	swagger.SwaggerInfo.Title = r.conf.Title
	swagger.SwaggerInfo.Description = r.conf.Description
	swagger.SwaggerInfo.Version = r.conf.Version
	swagger.SwaggerInfo.Host = r.conf.Host
	swagger.SwaggerInfo.BasePath = r.conf.BasePath

	r.http.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (r *rest) Register() {
	r.RegisterSwagger()

	publicApi := r.http.Group("/public")
	publicApi.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"msg": "hello world",
		})
	})

	api := r.http.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.POST("/register", r.RegisterUser)
	auth.POST("/login", r.LoginUser)
	auth.GET("/me", r.Me, r.VerifyUser())
}
