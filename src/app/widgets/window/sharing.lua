module(..., package.seeall)

sharing = class("sharing", mm.window)

function sharing:ctor(stage)
	sharing.super.ctor(self, stage, 2, true)
	self:setTitle("sharing")
end

function sharing:init()
	local panel = self:getWidget("panel")
	
end