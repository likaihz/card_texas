-- xxui.Label --
module(..., package.seeall)

-- class Label --
Label = class("Label", function()
    return xxui.Widget.new()
end)
-- args = {
--    node, num, img = {"img_NM", "img_HL"}, dir = "x"/"y",
--    gap = , d1, d2, h, func = i_th label's touchEvent
-- }
function Label:ctor(param)
    self.args = self:init(param)
    local img = param.img[1]
    local anch = cc.p(0, 0)
    anch[self.args.dir] = 0.5
    local num = self:getBtn("num")
    for i = 1, num do
        xxui.create {
            node = self, btn = {img, img},
            name = "Label_btn"..i, anch = anch,
            func = function()
                self:onClicked(i)
            end
        }
    end
    local g = function() end
    self.func = {on = g, off = g}
end

-- interface --
-- i: btn index or "num"
function Label:getBtn(i)
    if i == "num" then
        return self.args.num
    end
    return self:getleaf("Label_btn"..i)
end

-- key: "on" or "off"
function Label:setEvent(func, key)
    key = (key == "off") and key or "on"
    self.func[key] = func or function() end
end

function Label:addEvent(func, key)
    if not func then return end
    key = (key == "off") and key or "on"
    local f = self.func[key]
    self.func[key] = function(i)
        f(i); func(i)
    end
end

function Label:onClicked(i)
    if i == self.select then
        return
    end
    local tmp = self.select
    self.select = i
    self:turnOn(i)
    if tmp then
        self.func.off(tmp)
    end
    self.func.on(i)
end

function Label:setTouchEnabled(flag, i)
    if i then
        local w = self:getBtn(i)
        return w:setTouchEnabled(flag)
    end
    local num = self:getBtn("num")
    for j = 1, num do
        local w = self:getBtn(j)
        return w:setTouchEnabled(flag)
    end
end

-- implemetation --
function Label:turnOn(i)
    local args = self.args
    for j = 1, args.num do
        local idx = (j == i) and 2 or 1
        local img = args.img[idx]
        local w = self:getBtn(j)
        w:load({img, img})
        w:set{pos = self:pos(j)}
    end
end

function Label:init(param)
    self:addTo(param.node)
    local args = clone(param)
    args.dir = (param.dir == "y") and "y" or "x"
    args.gap = param.gap or 0
    args.d2 = param.d2 or param.d1
    local num, gap = args.num, args.gap
    local d1, d2 = args.d1, args.d2
    local w = gap * (num+1) + d1 * num + (d2-d1)
    local h = args.h
    local sz = (args.dir == "x") and {w, h} or {h, w}
    self:setsize(sz)
    return args
end

function Label:pos(i)
    local args = self.args
    local gap = args.gap
    local d1, d2 = args.d1, args.d2
    local s = (i == self.select) and d2/2 or d1/2
    local s1 = (i > self.select) and d2 - d1 or 0
    local len = gap * i + d1 * (i-1) + s + s1
    if args.dir == "x" then
        return cc.p(len, 0)
    else
        return cc.p(0, -len)
    end
end