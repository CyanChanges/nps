package customDev

import (
	"ehang.io/nps/server"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strconv"
)

func ApiWebServer() {
	app := fiber.New()

	setupRoutes(app)

	_ = app.Listen(":" + ApiPort)
}

// Set Routes
func setupRoutes(app *fiber.App) {
	// set handler for index page
	app.Get("/api/freePort", GetFreePort)        // set handler for index page
	app.Get("/api/randHttpProxy", RandHttpProxy) // set handler for index page
}

// 返回一个空闲可用端口, 注意防火墙开启端口
func GetFreePort(c *fiber.Ctx) error {
	availablePort := strconv.Itoa(FindFreePort())

	return c.SendString(availablePort)
}

// 返回代理
func RandHttpProxy(c *fiber.Ctx) error {
	u, err := url.Parse(c.BaseURL())
	if err != nil {
		return err
	}

	result := getProxy(u.Hostname())

	return c.JSON(result)
}

func getProxy(host string) (result map[string]interface{}) {
	// 通过 nps 服务端内置的列队来获取代理开放的端口, 种类 httpProxy
	//list, cnt := server.GetClientList(0, 100, "", "", "", 0)  // 客户端列表
	list, num := server.GetTunnel(0, 100, "httpProxy", 0, "") // 隧道列表

	var (
		m = map[string]interface{}{"proxies": []Proxy{}}
		p []*Proxy
	)

	if num <= 0 {
		// 还没有任何代理
		return nil
	}

	// 合成 json, 代理是私人代理, 所以需要提供帐号密码授权
	for _, item := range list {
		tmp := Proxy{fmt.Sprintf("http://%s:%d", host, item.Port), item.Client.Cnf.U, item.Client.Cnf.P}
		p = append(p, &tmp)
	}

	m["proxies"] = p

	return m
}
