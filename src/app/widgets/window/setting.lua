module(..., package.seeall)

setting = class("setting", mm.window)

function setting:ctor(stage)
	setting.super.ctor(self, stage, 2, true)
	self:init()
	self:setTitle("setting")
end

function setting:init()
	local panel = self:getWidget("panel")
	self:createbtn()
end

function setting:createbtn()
	local btn = self:setBtn("logout")
	btn["logout"]:setTouchEvent(function()
		print("logout...")
	end)
end