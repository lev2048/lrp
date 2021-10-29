package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"lrp/internal/lrp"
	"lrp/internal/status"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

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
		c.JSON(404, map[string]interface{}{
			"msg":  "请求方法不存在",
			"code": -3,
		})
	})

	Router := api.engine.Group("/v1")
	{
		Router.GET("status", api.Login)
		Router.GET("dashbord", api.GetDashBoardInfo)

		Router.POST("login", api.Login)
		Router.POST("proxy/add", api.AddProxy)
		Router.POST("proxy/del", api.DelProxy)
	}
}

func (api *Api) Login(c *gin.Context) {
	token := c.PostForm("token")
	if v, err := strconv.Atoi(token); err != nil {
		c.JSON(401, map[string]interface{}{
			"msg":  "认证失败",
			"code": -1,
		})
		return
	} else {
		if api.lrps.CheckToken(uint32(v)) {
			api.token = api.hmacSha256(token)
			c.JSON(200, map[string]interface{}{
				"msg":  "认证成功",
				"code": 1,
				"data": api.token,
			})
		} else {
			c.JSON(401, map[string]interface{}{
				"msg":  "认证失败",
				"code": -2,
			})
		}
	}
}

func (api *Api) AddProxy(c *gin.Context) {
	mark, cid := c.PostForm("mark"), c.PostForm("clientId")
	destAddr, ListenPort := c.PostForm("destAddr"), c.PostForm("listenPort")
	if cid == "" {
		c.JSON(500, map[string]interface{}{
			"msg":  "添加失败",
			"code": -1,
			"info": "客户端id不能为空",
		})
		return
	}
	if destAddr == "" {
		c.JSON(500, map[string]interface{}{
			"msg":  "添加失败",
			"code": -1,
			"info": "目标地址不能为空",
		})
		return
	}
	if err := api.lrps.AddProxy(cid, destAddr, mark, ListenPort); err != nil {
		c.JSON(500, map[string]interface{}{
			"msg":  "添加失败",
			"code": -2,
			"info": err.Error(),
		})
	} else {
		c.JSON(200, map[string]interface{}{
			"msg":  "添加成功",
			"code": 1,
		})
	}
}

func (api *Api) DelProxy(c *gin.Context) {
	cid, pid := c.PostForm("cid"), c.PostForm("pid")
	if cid == "" || pid == "" {
		c.JSON(500, map[string]interface{}{
			"msg":  "删除失败",
			"code": -1,
			"info": "客户端id或代理id不能为空",
		})
	} else {
		if err := api.lrps.DelProxy(cid, pid); err != nil {
			c.JSON(500, map[string]interface{}{
				"msg":  "删除失败",
				"code": -2,
				"info": err.Error(),
			})
		} else {
			c.JSON(200, map[string]interface{}{
				"msg":  "删除成功",
				"code": 1,
			})
		}
	}
}

func (api *Api) GetDashBoardInfo(c *gin.Context) {
	if info, err := json.Marshal(api.lrps.GetServerInfo()); err != nil {
		c.JSON(500, map[string]interface{}{
			"msg":  "获取失败",
			"code": -1,
			"info": err.Error(),
		})
	} else {
		c.JSON(200, map[string]interface{}{
			"msg":  "获取成功",
			"code": 1,
			"info": string(info),
		})
	}
}

func (api *Api) GetServerInfo(c *gin.Context) {
	if res, err := api.monitor.Info(); err == nil {
		c.JSON(500, map[string]interface{}{
			"msg":  "获取失败",
			"code": -1,
			"info": err.Error(),
		})
	} else {
		c.JSON(200, map[string]interface{}{
			"msg":  "获取成功",
			"code": 1,
			"info": res,
		})
	}
}

func (api *Api) hmacSha256(data string) string {
	h := hmac.New(sha256.New, []byte(api.secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
