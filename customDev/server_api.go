package customDev

import (
	"ehang.io/nps/lib/file"
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
	app.Get("/api/adslExpiry", AdslExpiry)
	app.Get("/api/freePort", GetFreePort)
	app.Get("/api/randHttpProxy/:amount?", RandHttpProxy)
}

// 返回一个空闲可用端口, 注意防火墙开启端口
func GetFreePort(c *fiber.Ctx) error {
	availablePort := strconv.Itoa(FindFreePort())

	return c.SendString(availablePort)
}

// 返回VPS拨号间隔
func AdslExpiry(c *fiber.Ctx) error {
	return c.JSON(serverPppoeExpiry)
}

// 返回代理
func RandHttpProxy(c *fiber.Ctx) error {
	result := getProxy(c)

	return c.JSON(result)
}

func getProxy(c *fiber.Ctx) (result map[string]interface{}) {
	needAmount, _ := strconv.Atoi(c.Params("amount")) // 客户端需要代理数量

	// 通过 nps 服务端内置的列队来获取代理开放的端口, 种类 httpProxy
	//list, cnt := server.GetClientList(0, 100, "", "", "", 0)  // 客户端列表
	list, num := server.GetTunnel(0, 100, "httpProxy", 0, "") // 隧道列表

	if num <= 0 {
		// 还没有任何代理
		return nil
	}

	if needAmount <= 0 {
		needAmount = 1
	}

	var (
		chooseList []*file.Tunnel
		m          = map[string]interface{}{"proxies": []Proxy{}}
		p          []*Proxy
	)

	// 随机获取N个代理
	chooseList = RandChooseByNums(list, needAmount)

	u, _ := url.Parse(c.BaseURL()) // 服务器地址

	// 合成 json, 代理是私人代理, 所以需要提供帐号密码授权
	for _, item := range chooseList {
		tmp := Proxy{fmt.Sprintf("http://%s:%d", u.Host, item.Port), item.Client.Cnf.U, item.Client.Cnf.P}
		p = append(p, &tmp)
	}

	m["proxies"] = p

	return m
}
