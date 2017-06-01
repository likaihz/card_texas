local hall = class("hall", function()
    return display.newScene("hall")
end)

function hall:ctor()
	self:init()
	self.topbar = mm.topbar.new(self)
end

-- interface --

-- create widgets --
function hall:init()
	xxui.create {
		node = self, img = xxres.scene("hall")
	}
	self:createBtn()
end

function hall:createBtn()
	xxui.create {
		node = self, btn = xxres.icon("niuniu"),
		pos = "center", func = function()
			GM:pushscene("front", "niuniu")
		end
	}
end

return hall