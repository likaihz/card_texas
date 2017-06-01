-- obj class table --
module(..., package.seeall)

-- class obj (base)
local obj = class("obj")

function obj:ctor(data)
	self.data = data
end

function obj:dtor()
	return self.data
end

-- interface --
function obj:get(k)
	local tbl = {img = 0, portrait = 0}
	if tbl[k] then
		return self:getImg()
	end
	return self.data[k]
end

function obj:set(k, v)
	self.data[k] = v
end

function obj:add(k, v)
	v = v + self:get(k)
	self:set(k, v)
end

function obj:getpth()
	local cls = self:get("cls")
	return "obj/"..cls.."/"
end

-- implementation --
function obj:getImg()
	local id = self:get("id")
	local pth = self:getpth()
	return pth..id..".png"
end
