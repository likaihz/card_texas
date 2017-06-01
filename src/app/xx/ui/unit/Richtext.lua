-- xxui.Richtext --
module(..., package.seeall)

Richtext = class("Richtext", function()
    return xxui.Widget.new()
end)

local FONT_SIZE = 40 -- default font size

-- param = {gap, width ...}
function Richtext:ctor(param)
    -- self.args: used to create single txt
    self.args = self:check(param)
    -- init args --
    self.cls = "Richtext"
    self.charLen = string.len("日")
    self.gap = param.gap or 12
    self.size = cc.size(param.width, 0)
    self:set(param)
    self:load(param.txt)
end

-- interface --
function Richtext:get(p)
    return self[p]
end

-- txt: string or table of strings
function Richtext:load(txt)
    self.size.height = 0
    self.rem = self.size.width
    self.num = self.num or 0
    self.pos = cc.p(0, 0)
    if self.num > 0 then
        for i = 1, self.num do
            local w = self:getChildByTag(i)
            w:removeSelf()
        end
        self.num = 0
    end
    if txt then
        self:push(txt)
    end
end

function Richtext:insert(str)
    if type(str) ~= "string" then
        xx.error(str)
    end
    local tbl = string.split(str, "/")
    for i, seg in ipairs(tbl) do
        local regular, args = self:interpret(seg)
        for j, reg in ipairs(regular) do
            local txts = string.split(reg, "\n")
            for k, txt in ipairs(txts) do
                if k > 1 then
                    self:newline(args)
                end
                args.txt = txt
                self:addWidget(args)
            end
        end
    end
end

function Richtext:push(content)
    if type(content) == "string" then
        self:insert(content)
    elseif type(content) == "table" then
        for i, str in ipairs(content) do
            self:insert(str)
        end
    else
        xx.error(content)
    end
end

-- key: "normal" / "hl"
function Richtext:setcolor(key, val)
    self.color = self.color or {}
    self.color[key] = val
end

function Richtext:getContentSize()
    return self.size
end

-- implementation --
function Richtext:check(param)
    local err = "Richtext error: "
    if param.anch then
        xx.error(err .. "anch = cc.p(0, 1)!", 0)
    end
    if not param.width then
        xx.error(err .. "width is needed!", 0)
    end
    local args = clone(param.args)
    args.node = self
    self:setcolor("normal", args.normal)
    self:setcolor("hl", args.hl)
    args.normal, args.hl = nil, nil
    return args
end

-- input: str; output: args table (to create a widget)
function Richtext:interpret(str)
    local args = clone(self.args)
    args.color = self.color.normal
    if string.find(str, "#") then
        args.color = self.color.hl
        str = string.sub(str, 2)
    end
    return string.regularize(str), args
end

function Richtext:addWidget(args)
    if args.txt == "" then return end
    local size = args.size or FONT_SIZE
    if self.num == 0 then
        self.size.height = size
        self.pos.y = -size
    end
    local len = {}
    local txt = args.txt
    if string.find(txt, "%w+") then
        len.char, size = 1, size/2
    else
        len.char = self.charLen
    end
    len.pixel = string.len(txt)/len.char * size
    local function createWidget(args)
        self.num = self.num + 1
        args.tag = self.num
        args.pos = self.pos
        return xxui.create(args)
    end
    while len.pixel - self.rem >= 0 do
        len.str = math.floor(self.rem/size) * len.char
        args.txt = string.sub(txt, 1, len.str)
        createWidget(args)
        len.pixel = len.pixel - self.rem
        txt = string.sub(txt, len.str+1)
        self:newline(args)
    end
    if len.pixel == 0 then return end
    args.txt = txt
    local w = createWidget(args)
    self.pos.x = self.pos.x + w:getContentSize().width
    self.rem = self.rem - len.pixel
end

function Richtext:newline(args)
    local hei = self.gap + (args.size or FONT_SIZE)
    self.size.height = self.size.height + hei
    self.pos.x = 0
    self.pos.y = self.pos.y - hei
    self.rem = self.size.width
end