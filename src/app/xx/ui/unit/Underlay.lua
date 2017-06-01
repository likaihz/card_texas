-- xxui.Underlay --
module(..., package.seeall)

Underlay = class("Underlay", function()
    return xxui.Widget.new()
end)
-- hidden: nodes need to be hidden
-- mode = 0: no return, 1: btn return, 2: screen return
-- mode = -1: btn invisible, -2: screen invisible
-- mode >= 0: delete, or setVisible(false)
-- opacity: number (0 ~ 1)
function Underlay:ctor(stage, mode, hidden, opacity)
    self.stage = stage; self:addTo(stage)
    self.mode_ = mode
    self.hidden = self:packup(hidden)
    xxui.turnswitch(nil, hidden)
    local w = self:createUI(math.abs(mode), opacity)
    if w then self:closedBy(w) end
    self:setTouchEnabled(true)
end

-- interface --
function Underlay:getpanel()
    return self:getchild("underlay@panel")
end

function Underlay:setTouchEnabled(flag)
    local panel = self:getpanel()
    panel:setTouchEnabled(flag)
end

function Underlay:disable()
    local panel = self:getpanel()
    panel:setTouchEvent(function() end)
end

function Underlay:setTouchEvent(func)
    local w = self:getpanel()
    w:setTouchEvent(func)
    if math.abs(self.mode_) == 2 then
        self:closedBy(w)
    end
end

function Underlay:setHidden(node)
    local tbl = self.hidden
    self.hidden = self:packup(node)
    xxui.turnswitch(tbl, self.hidden)
end

function Underlay:addHidden(node)
    local tbl = self:packup(node)
    xxui.turnswitch({}, tbl)
    table.merge(self.hidden, tbl)
end

function Underlay:switch(flag)
    if flag then
        xxui.turnswitch(self, self.hidden)
    else
        xxui.turnswitch(self.hidden, self)
    end
end

function Underlay:removeSelf()
    xxui.turnswitch(self.hidden)
    xxui.clean(self, 0)
    self:removeFromParent()
end

-- implemetation --
function Underlay:closedBy(w)
    w:addTouchEvent(function()
        if self.mode_ < 0 then
            self:switch(false)
        else
            self:removeSelf()
        end
    end)
end

function Underlay:packup(node)
    local t = type(node)
    if not node or t == "userdata" then
        return {node}
    elseif t == "table" then
        return node
    end
    xx.error(node)
end

-- create widgets --
-- mode: positive number
function Underlay:createUI(mode, opac)
    self:setLocalZOrder(1)
    local x, y = display.width, display.height
    -- create color node --
    local node = cc.DrawNode:create():addTo(self)
    node:drawPolygon({
        cc.p(0, 0), cc.p(x, 0), cc.p(x, y), cc.p(0, y)
    }, {fillColor = cc.c4f(0, 0, 0, opac or 0.7)})
    local panel = xxui.create {
        node = self, img = xxres.panel("underlay"),
        name = "underlay@panel", scale9 = cc.size(x, y)
    }
    panel:setTouchEnabled(true)
    self:setsize(panel)
    -- return the widget that hold close event
    if mode == 0 then return end
    if mode == 2 then return panel end
    if mode == 1 then
        return xxui.Button.new {
            node = panel, mold = "return",
            name = "rtn", anch = cc.p(0, 1), pos = cc.p(0, y)
        }
    end
    xx.error("underlay: invalid mode!", 0)
end