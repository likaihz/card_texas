module(..., package.seeall)

-- class topbar --
topbar = class("topbar", function()
    return xxui.Widget.new()
end)

function topbar:ctor(stage)
    self:addTo(stage)
    self:init()
    self:load()
end

-- interface --
-- key: "name", "uid", "roomcard"
function topbar:load(key)
    if key then
        return self:loadone(key)
    end
    local tbl = {name = 0, uid = 0, roomcard = 0}
    for k in pairs(tbl) do
        self:loadone(k)
    end
end

-- name: "hall", "niuniu"...
function topbar:settitle(name)
    local w = self:getleaf("title")
    w:load(xxres.icon(name))
end

-- implementation --
function topbar:setevent(name)
    print("pop window: ", name)
end

function topbar:loadone(key)
    local val = GM.player:get(key)
    if key == "uid" then
        val = "ID: "..val
    end
    local w = self:getleaf(key)
    if w then
        w:load(val)
    end
end

-- create widgets --
function topbar:init()
    local sz = cc.size(display.width, 130)
	self:setsize(sz)
    self:set {anch = cc.p(0.5, 1), align = cc.p(0.5, 1)}
    xxui.create {
        node = self, img = xxres.panel("bar"), scale9 = sz,
        anch = cc.p(0.5, 1), align = cc.p(0.5, 1)
    }
    self:createAvt()
    self:createBtn()
    self:createTitle()
end

function topbar:createAvt()
    local grid = xxui.create {
        node = self, img = xxres.grid("avt"), name = "grid_avt",
        anch = cc.p(0, 1), align = cc.p(0, 1) 
    }
    xxui.create {
        node = grid, txt = "rockhim", name = "name", size = 32,
        anch = cc.p(0, 1), align = cc.p(1, 1), pos = cc.p(10, -5)
    }
    xxui.create {
        node = grid, txt = "ID: 123456", name = "uid", size = 32,
        anch = cc.p(0, 1), align = cc.p(1, 1), pos = cc.p(10, -45)
    }
end

function topbar:createTitle()
    local panel = xxui.create {
        node = self, img = xxres.panel("title"),
        anch = cc.p(0.5, 1), align = cc.p(0.5, 1)
    }
    xxui.create {
        node = panel, img = xxres.icon("niuniu"), name = "title",
        anch = cc.p(0.5, 0.5), align = cc.p(0.5, 0.58)
    }
end

function topbar:createBtn()
    local tbl = {"xxx", "yyy", "setting"}
    local len = table.len(tbl)
    local int = 100
    for i, name in ipairs(tbl) do
        local x = (i-len) * int - 10 
        xxui.create {
            node = self, btn = xxres.button(name),
            anch = cc.p(1, 0.5), align = cc.p(1, 0.5), pos = cc.p(x, 0),
            func = function()
                self:setevent(name)
            end
        }
    end
end