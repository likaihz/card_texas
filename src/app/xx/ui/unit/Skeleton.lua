module(..., package.seeall)

Skeleton = class("Skeleton", function()
	return xxui.Widget.new()
end)



function Skeleton:ctor(param)
	local args = clone(param)
	self:init(args)
	if args.ani then
		self:load(args.ani)
	end
end

-- interface --
function Skeleton:load(name)
	if not name then
		return
	end
	local s = self:getchild("skeleton@ani")
	if s then
		s:removeFromParent()
	end
	if name == "" then
		return
	end
	local pth = "animation/skeleton/"..name
	local replace = self:getchild("replace")
	replace:setVisible(false)
	if not cc.FileUtils:getInstance():isFileExist(pth..".png") then
		self._s = nil
		local obj = mm.create(name)
		local img = obj:get("img")
		img = img:gsub(".png", "_full.png")
		replace:load(img)
		replace:setVisible(true)
		print("skeleton "..name.." not found!!!")
		return
	end
	s = sp.SkeletonAnimation:create(pth..".json", pth..".atlas")
	xxui.set(s, {node = self, align = cc.p(0.5, 0), name = "skeleton@ani"})
	self._s = s
end

-- isloop:是否循环播放，默认true, false为播放一次
function Skeleton:setAnimation(name, isloop, track)
	if not self._s then return end
	track = track or 0
	name = name or "idle"
	isloop = isloop or isloop == nil
	self._s:setAnimation(track, name, isloop)
end

function Skeleton:setTimeScale(num)
	self._s:setTimeScale(num)	--adjust play speed*num
end

-- implementation --
function Skeleton:init(args)
    self:setsize(args.size)   -- full obj box size
    self:set(args)
    xxui.create {
    	node = self, img = "", name = "replace",
    	anch = cc.p(0.5, 0), align = cc.p(0.5, 0)
	}
end

------------------------------------------------
function Skeleton:test(stage)
	local w = self.new{
		node = stage, ani = "pet_thunder_1", 
		size = cc.size(360, 400), pos = "center"
	}
	w:setAnimation("idle")
	-- w:setTimeScale(2)
end