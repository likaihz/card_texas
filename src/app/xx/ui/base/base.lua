-- base widgets --
module(..., package.seeall)

local Func = require "app.xx.ui.base.func"
local Anim = require "app.xx.ui.base.anim"
local Touch = require "app.xx.ui.base.touch"

-- base Widget Class
local Base = xx.inherit(nil, Func, Anim, Touch)

-- load: 纯虚函数
function Base:load() end

-- init: 初始化控件参数为默认值
local SIZE = {txt = 46, dft = 40}

function Base:init_(param)
    local args = clone(param or {})
    if self.type_ == "txt" then
        args.font = args.font or "dft"
        args.size = args.size or SIZE[args.font]
        args.color = args.color or "white"
    end
    args.anch = args.anch or cc.p(0, 0)
    args.pos = args.pos or cc.p(0, 0)
    self:load(args[self.type_], args.rect)
    self:set(args)
    self:setStateEvent()
end

-- externsion of cocos' widgets
local cocos = {
    ccui.Widget, ccui.Layout,
    ccui.ImageView, ccui.Button,
    cc.Node, cc.Sprite, cc.Sprite3D
}

for _, w in ipairs(cocos) do
    xx.inherit(w, Base)
end


-- class Widget --
Widget = class("Widget", function()
    return ccui.Widget:create()
end)

function Widget:ctor(param)
    self:init_(param)
end


-- class Img --
Img = class("Img", function()
    return ccui.ImageView:create()
end)

function Img:ctor(param)
    self.type_ = "img"
    self:init_(param)
end

function Img:load(content, rect)
    self:loadTexture(content)
    if rect then
        self:setTextureRect(rect)
    end
end

-- class Btn --
Btn = class("Btn", function()
    return ccui.Button:create()
end)

function Btn:ctor(param)
    self.type_ = "btn"
    self:init_(param)
end

function Btn:load(content)
    content = content or ""
    if type(content) == "string" then
        content = {content}
    elseif type(content) ~= "table" then
        xx.error(content)
    end
    content[2] = content[2] or ""
    content[3] = content[3] or content[1]
    -- 1:normal 2:selected 3:disabled
    self:loadTextures(unpack(content))
end

function Btn:enable(flag)
    self:setTouchEnabled(flag)
    self:setBright(flag)
end

-- class Txt --
Txt = class("Txt", function()
    return cc.Label:create()
end)

function Txt:ctor(param)
    self.type_ = "txt"
	self:init_(param)
    self.canTouch = true
end

function Txt:load(content)
	self:setString(content or "")
end

function Txt:setTouchEnabled(flag)
    self.canTouch = flag
end

function Txt:isTouchEnabled()
    return self.canTouch
end


-- create: create widget (Img, Btn, Txt)
-- args[img/txt/btn] == "", create a blank widget
function create(args)
    local node, name = args.node, args.name
    local w
    if node and name then
        w = node:getChildByName(name)
    end
    if not w then
        local tbl = {img = Img, txt = Txt, btn = Btn}
        for k, widget in pairs(tbl) do
            if args[k] then
                return widget.new(args)
            end
        end
    	xx.error("error input!", 0)
    end
    w:setVisible(true)
    w:load(args[w.type_], args.rect)
    return w
end

-- create base widgets using scale or scale9
function strech(param)
    local args = clone(param)
    local node = Widget.new()
    local tbl = {
        img = 0, btn = 0, txt = 0, flip = 0
    }
    for k, v in pairs(tbl) do
        tbl[k], args[k] = args[k], nil
    end
    tbl.node = node
    local bg = create(tbl)
    if args.scale then
        bg, size = bg:scale(args.scale)
        node:setsize(size)
        args.scale = nil
    end
    if args.scale9 then
        bg = bg:scaleNine(args.scale9)
        node:setsize(bg)
        args.scale9 = nil
    end
    if tbl.flip then
        bg:set { anch = cc.p(0.5, 0.5), align = cc.p(0.5, 0.5) }
    end
    node:set(args)
    return node, bg
end