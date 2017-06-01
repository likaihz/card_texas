-- class GameManager --
-- 管理玩家（player）数据，关卡（chapter）数据
-- 实现数据在服务器的存取接口
local GameManager = class("GameManager")

function GameManager:ctor()
	-- xx.fcheck("config")
	-- xx.fcheck("save")
	self.scenes = {}
	self.ready = {data = false, res = false}
	self.account = {}
end

-- interface of scene --
-- 进入新场景
function GameManager:pushscene(scene, ...)
	if type(scene) == "string" then
	    scene = require("app.scenes."..scene).new(...)
	    cc.Director:getInstance():pushScene(scene)
	end
    table.insert(self.scenes, scene)
    return scene
end

-- 弹出当前场景
function GameManager:popscene()
    cc.Director:getInstance():popScene()
    return table.remove(self.scenes)
end

-- 返回主场景
function GameManager:mainscene()
	for i = #self.scenes, 2, -1 do
		self:popscene()
	end
	return self.scenes[1]
end

-- 获取当前场景
function GameManager:getscene()
	return self.scenes[#self.scenes]
end

-- 资源预加载
-- data: {img = {"ui", "scene"}, ani = {"battle"}}
-- update: func called every png loaded
function GameManager:loadres(data, update)
	local pth, num = "app.data.res.", 0
	for cls, tbl in pairs(data) do
		if cls == "ani" then
			pth = pth.."animation."
		end
		for i, dir in ipairs(tbl) do
			tbl[i] = require(pth..dir)
			num = num + tbl[i].num
		end
	end
	for _, tbl in pairs(data) do
		for _, list in ipairs(tbl) do
			for _, data in ipairs(list) do
			    xx.cache(data, function()
			    	if update then update(num) end
			    end)
			end
		end
	end
end

-- 玩家数据加载
-- data: table
function GameManager:loadusr(data)
	local Player = require("app.control.Player")
	self.player = Player.new(data)
end

-- interface of file --
-- 读取配置表
-- key: config file's name
function GameManager:getconfig(key)
	return clone(self.config[key])
end

-- 存储配置表
function GameManager:setconfig(cfg)
	self.config = cfg or {}
end

-- 设置并获取设备ID
function GameManager:getdevice()
	local tbl = xx.io("save/device")
	if not tbl then
		local id = tostring(os.time())
		tbl = {id = id}
		xx.io("save/device", tbl)
	end
	return tbl.id
end

-- interface of web --
-- 获取 websocket 连接
function GameManager:getsocket()
	if not self.ws then
		local url = self:geturl("ws", "receive")
		self.ws = xx.websocket.new(url)
	end
	return self.ws
end

function GameManager:setsocket(func)
	local ws = self:getsocket()
	ws:onmessage(func or function() end)
end

-- 通过 websocket 连接发送内容
-- msg: table
function GameManager:send(msg)
	if not self.ws then
		xx.error("socket hasn't been created!", 0)
	end
	msg.uid = self.player:get("uid")
	self.ws:send(msg)
end

-- 向服务端发送 http 请求
-- req: table
function GameManager:request(req, func)
	local url = self:geturl("http", req.func)
	xx.request(url, req, function(res)
		local ok = self:checkresponse(res)
		if not ok then
			return
		end
		if func then func(res) end
	end)
end

-- 向服务端请求登录操作
function GameManager:login(req, func)
end

-- 向服务端请求测试登录
function GameManager:test(func)
	local req = {func = "test"}
	self:request(req, function(res)
		local data = res.user
		if not data then
			print("database error!!")
			return
		end
		self:setconfig(res.cfg)
		self:loadusr(data)
		if func then func(res) end
	end)

end

-- implementation --
function GameManager:getaccount(req)
	local acc = self.account
	req.account = acc.account
	req.password = acc.password
	req.device = acc.device
	return req
end

function GameManager:setaccount(req)
	req.device = self:getdevice()
	local acc = {}
	acc.account = req.account
	acc.password = req.password
	acc.device = req.device
	self.account = acc
end

function GameManager:checkresponse(res)
	local s = res.status
	if s == "fail" then
		print("network failed!")
		return false
	end
	return true
end

local NET = "115.159.208.44"
local LOC = "127.0.0.1"

local PORT = {
	login = "8001", test = "8001",
}

function GameManager:geturl(mode, key)
	local url = LOC
	local port = PORT[key] or "8000"
	return table.concat {
		mode, "://", url, ":", port, "/", key
	}
end

return GameManager