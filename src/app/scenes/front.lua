local front = class("front", function()
    return display.newScene("front")
end)

-- kind: 表征不同玩法
function front:ctor(kind)
	self.kind = kind
	print("in front: ", kind)
    self:init()
    self.topbar = mm.topbar.new(self)
end

function front:onEnter()
	GM:setsocket(function(msg)
		self:receive(msg)
	end)
end

function front:onExit()
	GM:setsocket()
end

-- interface --
function front:receive(msg)
	local ok, data = self:check(msg)
	if not ok then return end
	self:removeinput()
	GM:pushscene("room", data)
end

function front:check(msg) 
	if msg.opt ~= "front" then
		return
	end
	local s = msg.status
	if s == "ok" then
		return true, msg.data
	end
	self:roomnumtips(s)
end

-- implementation --
function front:makeroom(opt, num)
	local msg = {opt = opt, roomnum = num}
	msg.kind = self.kind
	local config = {
		roundnum = "one", turning = "fixed",
		mode = "crazy"
	}
	msg.config = config
	GM:send(msg)
end

function front:createinput()
	self.input = mm.roomnum.new(self, function(num)
		self:makeroom("access", num)
	end)
end

function front:removeinput()
	if self.input then
		self.input:remove()
	end
	self.input = nil
end

function front:roomnumtips(status)
	local txt = {"room", status, "!"}
	mm.tips.new(self.input, txt, "orange")
end

-- create widgets --
function front:init()
	xxui.create {
		node = self, img = xxres.scene("front")
	}
	self:createbtns()
end

function front:createbtns()
	xxui.create {
		node = self, btn = xxres.button("back"),
		anch = cc.p(0, 0.5), align = cc.p(0, 0.5), pos = cc.p(10, 0),
		func = function()
			-- GM:send {opt = "connect", open = "off"}
			GM:popscene()
		end
	}
	local btns = {"create", "access"}
	for i, opt in ipairs(btns) do
		btns[i] = self:roombtn(opt)
	end
	btns[1]:setTouchEvent(function()
		self:makeroom("create")
	end)
	btns[2]:setTouchEvent(function()
		self:createinput()
		-- self:makeroom("access", 100000)  -- for test
	end)
end

function front:roombtn(opt)
	local name = "room_"..opt
	local tbl = {create = -1, access = 1}
	local x = tbl[opt] * 0.2
	local btn = xxui.create {
		node = self, btn = xxres.button(name),
		anch = cc.p(0.5, 0.5), align = cc.p(0.5+x, 0.5),
		pos = cc.p(100, 0)
	}
	xxui.create {
		node = btn, img = xxres.grid("gold"),
		zorder = -1, anch = cc.p(0.67, 0.15)
	}
	return btn
end

return front