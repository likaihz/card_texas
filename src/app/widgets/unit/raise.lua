-- mm.raise --
module(..., package.seeall)

-- raise btns n countdown tips
raise = class("raise")

function raise:ctor(scene)
	self.scene = scene
    self:init()
end

-- interface --
function raise:setevent(func)
	for i, btn in pairs(self.btn) do
		btn:setTouchEvent(function()
			self:load(false)
			if func then func(i) end
		end)
	end
end

function raise:setcomplete(func)
	self.oncomplete = func
end

function raise:load(bool)
	self:show(bool)
	if not bool then
		return xx.unschedule(self.sch)
	end
	local w = self.tips:getchild("text")
	local txt = xx.translate("raise", "tips")
	local function oncount(i)
		w:load(txt..i)
	end
	self.sch = self:countdown(oncount, function()
		self:load(false)
		local f = self.oncomplete
		if f then f() end
	end)
end

function raise:show(bool)
	for _, w in pairs(self.btn) do
		w:setVisible(bool)
	end
	self.tips:setVisible(bool)
end

function raise:remove()
	for _, w in pairs(self.btn) do
		w:removeSelf()
	end
	self.tips:removeSelf()
end

-- create widgets --
function raise:init()
	self.btn = {}
	for i = 1, 3 do
		self.btn[i] = self:createbtn(i)
	end
	self.tips = self:createtips()
end

function raise:createbtn(i)
	local txt = xx.translate("raise")
	txt = i.." "..txt
	local x = display.cx + (i-2)*220
	local y = 250
	local btn = xxui.Txtbtn.new {
		node = self.scene, mode = "orange", text = txt,
		size = "small", anch = cc.p(0.5, 0.5), pos = cc.p(x, y)
	}
	btn:setVisible(false)
	return btn
end

function raise:createtips()
	local tips = xxui.create {
		node = self.scene, img = xxres.panel("bar_gray"),
		pos = "center"
	}
	xxui.create {
		node = tips, txt = "", name = "text",
		size = 40, pos = "center"
	}
	tips:setVisible(false)
	return tips
end

function raise:countdown(oncount, oncomplete)
	return xx.countdown(7, function(time)
		local i = math.round(time)
		if math.abs(time - i) < 0.03 then
			oncount(i)
		end
	end, oncomplete)
end