-- xxui.Txtbtn --
module(..., package.seeall)

-- param@
-- node: root of button
-- text: text of button
-- mode: color of button
-- size: "big", "small", nil (normal)

Txtbtn = class("Txtbtn", xxui.Btn)

local OUTLINE = {
    blue = cc.c3b(40, 95, 150), green = cc.c3b(45, 115, 60),
    orange = cc.c3b(104, 106, 57), gray = cc.c3b(85, 85, 85),
}

-- param: {node, mode, size, text}
function Txtbtn:ctor(param)
    local args = clone(param)
    self.size = args.size or ""
    self.mode = args.mode
    args.btn = self:getimg(true)
    self:setTouchEvent()
    Txtbtn.super.ctor(self, args)
    self:init(args.text)
    self:enable(true)
end

-- interface --
function Txtbtn:enable(bool)
    self:setTouchEnabled(bool)
    self:loadimg(nil, bool)
    return self
end

function Txtbtn:loadtxt(txt)
    local w = self:getchild("Txtbtn_text")
    w:load(txt or "")
end

function Txtbtn:loadimg(mode, enable)
    self.mode = mode or self.mode
    local img = self:getimg(enable)
    self:load(img)
    -- change text's outline
    local c = enable and self.mode or "gray"
    c = OUTLINE[c]
    local txt = self:getchild("Txtbtn_text")
    txt:set {outline = {color = c}}
end

function Txtbtn:setTouchEvent(func)
    xxui.setTouchEvent(self, func)
    xxui.addTouchEvent(self, function()
        -- GM:playSound("btn/sound_9.mp3")
    end)
end

-- implementation --
function Txtbtn:getimg(enable)
    local mode = self.mode
    local size = self.size
    local img = enable and mode or "gray"
    if size ~= "" then
        img = img.."_"..size
    end
    return xxres.button(img)
end

function Txtbtn:init(txt)
    local tbl = {small = 40, big = 50}
    local sz = tbl[self.size] or 45
    xxui.create {
        node = self, txt = "", name = "Txtbtn_text",
        font = "txt", size = sz, pos = "center"
    }
    self:loadtxt(txt)
end
