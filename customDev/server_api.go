package customDev

import (
	"ehang.io/nps/lib/file"
	"ehang.io/nps/server"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strconv"
	"strings"
)

func ApiWebServer() {
	app := fiber.New()

	setupRoutes(app)

	_ = app.Listen(":" + ApiPort)
}

// Set Routes
func setupRoutes(app *fiber.App) {
	// set handler for index page
	app.Get("/api/freePort", GetFreePort)
	app.Get("/api/randHttpProxy/:amount?", RandHttpProxy)
	app.Get("/api/delClient", DelClient)
}

// 删除代理,以IP为搜索条件
func DelClient(c *fiber.Ctx) (err error) {
	list, num := server.GetClientList(0, 10000, "", "", "", 0)

	if num <= 0 {
		return c.SendString("not in system")
	}

	var clientId int

	// 从客户端列表里面找到对应客户端ID
	for _, item := range list {
		if strings.Contains(item.Addr, c.IP()) {
			clientId = item.Id
		}
	}

	if err := file.GetDb().DelClient(clientId); err != nil {
		return c.SendString("delete error")
	}
	server.DelTunnelAndHostByClientId(clientId, false)
	server.DelClientConnect(clientId)

	return c.SendString("delete success")
}

// 返回一个空闲可用端口, 注意防火墙开启端口
func GetFreePort(c *fiber.Ctx) error {
	availablePort := strconv.Itoa(FindFreePort())

	return c.SendString(availablePort)
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
		tmp := Proxy{fmt.Sprintf("http://%s:%d", u.Hostname(), item.Port), item.Client.Cnf.U, item.Client.Cnf.P}
		p = append(p, &tmp)
	}

	m["proxies"] = p

	return m
}
