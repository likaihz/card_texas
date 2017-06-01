-- mm.tips --
module(..., package.seeall)

local instant
-- class tips --
tips = class("tips", function()
	return xxui.Widget.new()
end)

function tips:ctor(stage, txt, color)
	if instant and instant.delete then
		instant:delete()
	end
	instant = self
	self:init(stage)
	local c = color or "green"
	self:load(txt, c)
end

function tips:delete()
	self:removeSelf()
	instant = nil
end

function tips:load(txt, color)
	txt = txt or ""
	if txt ~= "" then
		txt = xx.translate(txt, "tips")
	end
	local w = self:getleaf("text")
	self:setcolor(color)
	w:load(txt)
	self:move()
end

local COLOR = {
	green = cc.c3b(22, 90, 76),
	orange = cc.c3b(255, 145, 50)
}

function tips:setcolor(color)
	if not color then return end
	local c = COLOR[color]
	local w = self:getleaf("text")
	w:set {
		shadow = {color = c, offset = cc.size(0, -3)},
		outline = {color = c, size = 2}
	}
end

function tips:move()
	local w = self:getleaf("text")
	w:setopacity(0)
	w:fadein(0.15)
	self:moveby(0.15, cc.p(0, 10), function()
		self:moveby(1.7, cc.p(0, 40), function()
			self:fadeout(1, true, function()
				instant = nil
			end)
		end)
	end)
end

-- implementation --
function tips:init(stage)
	self:addTo(stage)
	local panel = xxui.create {
		node = self, img = xxres.panel("tips"),
		scale9 = cc.size(display.width, 70)
	}
	panel:setopacity(0.4)
	self:setsize(panel)
	self:set {
		anch = cc.p(0.5, 0.5), align = cc.p(0.5, 0.6), name = "tips"
	}
	xxui.create {
		node = self, txt = "", name = "text",
		size = 46, pos = "center"
	}
end

-- test --
function tips:test(node)
	self.new(node, "提示框")
end