package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"lrp/internal/lrp"
	"lrp/internal/status"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

type Result struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Info interface{} `json:"info,omitempty"`
}

type Api struct {
	token  string
	secret string

	lrps    *lrp.Server
	engine  *gin.Engine
	monitor *status.Monitor
}

func NewApi(lrps *lrp.Server, engine *gin.Engine, monitor *status.Monitor, secret string) *Api {
	return &Api{
		lrps:    lrps,
		secret:  secret,
		engine:  engine,
		monitor: monitor,
	}
}

func (api *Api) SetRouter() {
	api.engine.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, token",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	api.engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, Result{
			Msg:  "请求方法不存在",
			Code: -3,
		})
	})

	Router := api.engine.Group("/v1")
	{
		Router.GET("status", api.GetServerInfo)
		Router.GET("dashbord", api.GetDashBoardInfo)

		Router.POST("login", api.Login)
		Router.POST("proxy/add", api.AddProxy)
		Router.POST("proxy/del", api.DelProxy)
	}
}

func (api *Api) Login(c *gin.Context) {
	token := c.PostForm("token")
	if v, err := strconv.Atoi(token); err != nil {
		c.JSON(401, Result{
			Msg:  "认证失败",
			Code: -1,
		})
		return
	} else {
		if api.lrps.CheckToken(uint32(v)) {
			api.token = api.hmacSha256(token)
			c.JSON(200, Result{
				Msg:  "认证成功",
				Code: 1,
				Info: api.token,
			})
		} else {
			c.JSON(401, Result{
				Msg:  "认证失败",
				Code: -2,
			})
		}
	}
}

func (api *Api) AddProxy(c *gin.Context) {
	mark, cid := c.PostForm("mark"), c.PostForm("clientId")
	destAddr, ListenPort := c.PostForm("destAddr"), c.PostForm("listenPort")
	if cid == "" {
		c.JSON(500, Result{
			Msg:  "添加失败",
			Code: -1,
			Info: "客户端id不能为空",
		})
		return
	}
	if destAddr == "" {
		c.JSON(500, Result{
			Msg:  "添加失败",
			Code: -2,
			Info: "目标地址不能为空",
		})
		return
	}
	if err := api.lrps.AddProxy(cid, destAddr, mark, ListenPort); err != nil {
		c.JSON(500, Result{
			Msg:  "添加失败",
			Code: -3,
			Info: err.Error(),
		})
	} else {
		c.JSON(200, Result{
			Msg:  "添加成功",
			Code: 1,
		})
	}
}

func (api *Api) DelProxy(c *gin.Context) {
	cid, pid := c.PostForm("cid"), c.PostForm("pid")
	if cid == "" || pid == "" {
		c.JSON(500, Result{
			Msg:  "删除失败",
			Code: -1,
			Info: "客户端id或代理id不能为空",
		})
	} else {
		if err := api.lrps.DelProxy(cid, pid); err != nil {
			c.JSON(500, Result{
				Msg:  "删除失败",
				Code: -2,
				Info: err.Error(),
			})
		} else {
			c.JSON(200, Result{
				Msg:  "删除成功",
				Code: 1,
			})
		}
	}
}

func (api *Api) GetDashBoardInfo(c *gin.Context) {
	c.JSON(200, Result{
		Msg:  "获取成功",
		Code: 1,
		Info: api.lrps.GetServerInfo(),
	})
}

func (api *Api) GetServerInfo(c *gin.Context) {
	c.JSON(200, Result{
		Msg:  "获取成功",
		Code: 1,
		Info: api.monitor.Info(),
	})
}

func (api *Api) hmacSha256(data string) string {
	h := hmac.New(sha256.New, []byte(api.secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
