local Battle = require "app.models.battle.Battle"

local room = class("room", function()
    return display.newScene("room")
end)

local TIME = 0.1 -- second

function room:ctor(data)
	self.data = data
	self:init()
    self.battle = Battle.new(self)
end

function room:onEnter()
	local b = self.battle
	GM:setsocket(function(msg)
		b:receive(msg)
	end)
	b:send("enter")
	self:load()
	self:loadtime()
	self.clock = xx.schedule(function()
		self:loadtime()
	end, 60)
end

function room:onExit()
	self.clock = xx.unschedule(self.clock)
	GM:setsocket()
end

-- interface --
local KEYS = {"roundnum", "turning", "roomnum"}

function room:get(k)
	return self.data[k]
end

function room:load(k)
	if k then
		return self:loadone(k)
	end
	for _, k in pairs(KEYS) do
		self:loadone(k)
	end
end

function room:countdown(time, func)
	local w = xxui.create {
		node = self, txt = "", name = "text",
		size = 60, pos = "center"
	}
	local i = time
	xxui.schedule(self, 1, time, function()
		w:load(i)
		i = i - 1
	end, function()
		w:removeSelf()
		if func then func() end
	end)
end

function room:removechild(name)
	local w = xxui.getchild(self, name)
	if w then
		w:removeSelf()
	end
end

-- implementation --
function room:loadone(k)
	local w = self.panel:getleaf(k)
	if not w then return end
	local v = self:get(k)
	if k == "roundnum" then
		local n = self.battle:get("roundcnt")
		v = n.."/"..v
	elseif k == "turning" then
		v = xx.translate(v)
	end
	w:load(v)
end

function room:loadtime()
	local w = self.panel:getleaf("time")
	w:load(os.date("%H:%M"))
end

-- create widgets --
function room:init()
	xxui.create {
		node = self, img = xxres.scene("room")
	}
	self:createinfo()
	self:createbtns()
end

function room:createinfo()
	local panel = xxui.create {
		node = self, img = xxres.panel("info"),
		scale9 = {340, 290}, anch = cc.p(0, 1),
		align = cc.p(0, 1)
	}
	self:newtime(panel)
	for i, name in ipairs(KEYS) do
		self:newtxt(panel, name, i)
	end
	self.panel = panel
end

function room:createbtns()
	self:newbtn("return", function()
		GM:popscene()
	end)
	self:newbtn("setting", function()
		print("setting...")
	end)
	local txt = xx.translate("ready", "button")
	local ready = xxui.Txtbtn.new {
		node = self, mode = "orange", text = txt,
		name = "ready", pos = "center"
	}
	ready:setTouchEvent(function()
		self.battle:send("ready")
	end)
end

function room:newbtn(name, func)
	local x = 0.04
	if name == "setting" then
		x = 1 - x
	end
	local btn = xxui.create {
		node = self, btn = xxres.button(name),
		anch = cc.p(0.5, 0.5), align = cc.p(x, 0.93),
		func = func
	}
	return btn
end

function room:newtxt(node, name, i)
	local color = cc.c3b(100, 120, 140)
	local y = i * 42
	local txt = xx.translate {name, ":"}
	xxui.create {
		node = node, txt = txt, size = 30,
		color = color, pos = cc.p(40, y)
	}
	xxui.create {
		node = node, txt = "", name = name,
		size = 30, color = color, pos = cc.p(120, y)
	}
end

function room:newtime(panel)
	local color = cc.c3b(100, 120, 140)
	local time = xxui.create {
		node = panel, txt = "", name = "time",
		size = 50, color = color, align = cc.p(0.44, 0.66)
	}
end

return room