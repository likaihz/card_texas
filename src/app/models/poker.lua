-- module(..., package.seeall)

-- class poker --
local poker = class("poker", function()
	return xxui.Widget.new()
end)

function poker:ctor(scene)
	self.scene = scene
	self:init()
end

function poker:dtor()
	return self.data
end

-- interface --
-- k: "cls", "idx"
function poker:get(k)
	return self.data[k]
end

function poker:load(data)
	self.data = data
end

function poker:show(bool, mode, func)
	func = func or function() end
	if mode == "flip" then
		self:flip(0.7, "left", func)
	else
		self:loadimg(bool)
		func()
	end
end

function poker:move(ft, idx, func)
	local time = 0.3
	local pos = self:initpos(ft, idx)
	self:moveto(time, pos, func)
end

function poker:slide(idx, mode)
	local num = idx-1
	local time = 0.1
	local x = self.int * num
	self:moveby(time * num, cc.p(x, 0), function()
		self:show(true, mode)
	end)
end

function poker:cover(idx)
	self:show()
	self:scale(1.3)
	local x, y = self:getPosition()
	local dx = (3-idx) * 55
	local dy = 120
	self:setPosition(x+dx, y+dy)
end

-- implementation --
function poker:flip(time, spin, func)
	local t = time/2
    local f1 = cc.OrbitCamera:create(t, 1, 0, 0, 90, 0, 20)
    local f2 = cc.OrbitCamera:create(t, 1, 0, 270, 90, 160, 20)
    if spin == "left" then
        f1, f2 = f2:reverse(), f1:reverse()
    end
    local seq = cc.Sequence:create(
    	f1, cc.CallFunc:create(function()
            self:loadimg(true)
        end), f2, cc.CallFunc:create(func)
    )
    self:runAction(seq)
end

function poker:loadimg(bool)
	local img = self:getimg(bool)
	local w = self:getchild("img")
	w:load(img)
end

function poker:initpos(ft, idx)
	local bool = ft:is("me")
	self:scale(bool and 1.7 or 1.3)
	self.int = bool and 125 or 70
	local pos = ft:cardpos()
	pos.x = pos.x + (idx-3)*self.int
	return pos
end

-- create widgets --
function poker:init()
	local img = self:getimg()
	img = xxui.create {
		node = self, img = img, name = "img"
	}
	self:setsize(img)
	self:set {
		node = self.scene, pos = "center"
	}
end

local function getpth()
	return "obj/poker/poker_"
end

function poker:getimg(bool)
	local pth = getpth()
	if not self.data or not bool then
		return pth.."back.png"
	end
	local idx = self:get("idx")
	local cls = self:get("cls")
	local c = string.sub(cls, 1, 1)
	pth = table.concat {
		pth, c, idx, ".png"
	}
	return pth
end

return poker
