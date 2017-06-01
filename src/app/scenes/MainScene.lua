require "app.xx.init"
require "app.models.init"

GM = require("app.control.GameManager").new()

local MainScene = class("MainScene", function()
    return display.newScene("MainScene")
end)

function MainScene:ctor()
	self:init()
	-- test --
	mm.window.new(self, 2, false)
end

function MainScene:onEnter()
end

-- implementation --
function MainScene:socket()
    local ws = GM:getsocket()
    ws:onopen(function()
    	GM:send {opt = "connect", open = "on"}
    end)
end

function MainScene:login()
	GM:test(function(res)
    	self:socket()
    	GM:pushscene("hall")
	end)
end

-- create widgets --
function MainScene:init()
	xxui.create {
		node = self, img = xxres.scene("login")
	}
	mm.button.new {
		node = self, btn = "login", name = "login",
		anch = cc.p(0.5, 0.5), align = cc.p(0.5, 0.2),
		func = function() self:login() end
	}
end

return MainScene
