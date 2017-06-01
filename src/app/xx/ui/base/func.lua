module(..., package.seeall)

-- getsize: 获取控件的 contentSize
-- mode == 2 return wid, hei; or return cc.size()
function getsize(w, mode)
    local sz = w:getContentSize()
    if mode == 2 then
        return sz.width, sz.height
    else
        return sz
    end
end

-- setsize: 设置控件的 contentSize
-- mode: node or xxui.size
function setsize(w, mode)
    local sz
    if type(mode) == "userdata" then
        sz = mode:getContentSize()
    else
        local wid, hei = xxui.size(mode)
        sz = cc.size(wid, hei)
    end
    w:setContentSize(sz)
end

-- setcolor: 设置控件颜色
-- p1: = string/c3b/c4b or {color, opacity}
function setcolor(w, p1, p2)
    local c, o = clone(p1), p2
    if type(c) == "table" and c.color then
        c, o = c.color, c.opacity
    end
    local c4b = xxui.color(c, o)
    if w.setTextColor then
        w:setTextColor(c4b)
    else
        w:setOpacity(c4b.a)
        c4b.a = nil
        w:setColor(c4b)
    end
end

-- setopacity: 设置控件透明度
local function setopac(w, num, all)
    w:setOpacity(num)
    if not all then return end
    local children = w:getChildren()
    for _, child in pairs(children) do
        setopac(child, num, true)
    end
end

function setopacity(w, num, all)
    num = xx.restrict(num, 0, 1)
    num = num * 255
    setopac(w, num, all)
end

-- setlight: 设置控件亮度
-- num: number from 0 to 1
function setlight(w, num)
    if type(num) ~= "number" then
        xx.error(num)
    end
    if num < 0 then
        xx.error("out of range!", 0)
    end
    if num > 1 then num = 1 end
    local c = 255 * num
    w:setColor(cc.c3b(c, c, c))
end

-- clean: 清除或隐藏控件上的子节点
-- delete: boolean (true: 清除, false: 隐藏)
function clean(w, delete)
    local children = w:getChildren()
    for k, v in pairs(children) do
        if delete then
            v:removeSelf()
        else
            v:setVisible(false)
        end
    end
end

-- getchild: 获取名为 name 的子节点
-- name: string or number
function getchild(w, name)
    local tp = type(name)
    if tp == "string" then
        return w:getChildByName(name)
    end
    if tp == "number" then
        return w:getChildByTag(name)
    end
    xx.error(name) 
end

-- getleaf: 获取名为 name 的叶节点（重名时随机获取）
-- name: string or number
function getleaf(w, name)
    local c = getchild(w, name)
    if c then return c end
    local children = w:getChildren()
    for _, child in pairs(children) do
        local c = getleaf(child, name)
        if c then return c end
    end
end

-- turnswitch: 控件显示开关，显示 a 隐藏 b
-- a: node or table of nodes
function turnswitch(a, b)
    local tbl = {a or {}, b or {}}
    for i = 1, 2 do
        if type(tbl[i]) == "table" then
            for _, v in pairs(tbl[i]) do
                v:setVisible(i == 1)
            end
        else
            tbl[i]:setVisible(i == 1)
        end
    end
end


-- implementation of align --
local function getlen(w, dir)
    if tolua.type(w) == "cc.Scene" then
        return display[dir]
    elseif type(w) == "number" then
        return w
    elseif type(w) == "userdata" then
        return getsize(w)[dir]
    else
        xx.error(w)
    end
end

-- align: 等距排列函数，返回位置坐标值
-- mode: "x" or "y", else: number
function align(w, mode, num, index, gap, margin)
    local dir = {x = "width", y = "height"}
    if not dir[mode] then
        xx.error("error input mode!", 0)
    end
    local len = getlen(w, dir[mode])
    gap = gap or 0
    local m = margin or gap / 2
    local d = (len - 2 * m - (num-1) * gap) / num
    local pos = m + d/2 + (index - 1) * (d + gap)
    return (mode == "x") and pos or len - pos
end


local MaxRate = 10
-- scale: 控件整体缩放
function scale(w, size)
    if not size then return w end
    if type(size) ~= "number" and type(size) ~= "table" then
        xx.error(size)
    end
    local function cal(x, len)
        local MaxRate = 10
        return (x < MaxRate) and x or x/len
    end
    local wid, hei = getsize(w, 2)
    local x, y
    if type(size) == "number" then
        x, y = size, size
    else
        x, y = xxui.size(size)
    end
    x, y = cal(x, wid), cal(y, hei)
    w:setScale(x, y)
    local size = cc.size(x * wid, y * hei)
    return w, size
end


-- scaleNine: 控件九宫格缩放 (size, cap are tables)
-- cap为拉伸区域占控件的width，height比例 (0 ~ 1)，居中计算
local scalenine = {} -- implementation: scaleNine

function scaleNine(w, size, cap)
    if not size or size == 1 then
        return w
    end
    local sizeX, sizeY = xxui.size(size)
    local capX, capY
    if cap then
        capX, capY = xxui.size(cap)
    end
    local sz = scalenine.init(w)
    local x, wid, width = scalenine.calc(sz.width, sizeX, capX)
    local y, hei, height = scalenine.calc(sz.height, sizeY, capY)
    w:setCapInsets(cc.rect(x, y, wid, hei))
    w:setContentSize(cc.size(width, height))
    return w
end

-- implementation of scaleNine --
function scalenine.init(w)
    w:setScale9Enabled(false)
    w:setScale9Enabled(true)
    return w:getContentSize()
end

function scalenine.calc(length, size, cap)
    if cap and (cap > 1 or cap <= 0) then
        xx.error("wrong cap!!", 0)
    end
    local len = length
    length = (size < MaxRate) and size * length or size
    cap = cap or (length >= len and 0.2 or 1 - length/len)
    return len * (0.5 - cap/2),  len * cap, length
end


-- setpos: 设置锚点、对齐及坐标
-- param = {anch, align, pos}
function setpos(w, param)
    if not param then return end
    local anch = clone(param.anch)
    local align = clone(param.align)
    local pos = clone(param.pos)
    if pos == "center" then
        align = cc.p(0.5, 0.5)
        anch, pos = cc.p(0.5, 0.5), nil
    end
    if anch then
        if w.setAnchor then
            w:setAnchor(anch)
        else
            w:setAnchorPoint(anch)
        end
    end
    if not align and not pos then return end
    local x = pos and pos.x or 0
    local y = pos and pos.y or 0
    if align then
        local parent = w:getParent()
        local wid, hei
        if tolua.type(parent) == "cc.Scene" then
            wid, hei = display.width, display.height
        else
            wid, hei = getsize(parent, 2)
        end
        x = x + align.x * wid
        y = y + align.y * hei
    end
    w:setPosition(cc.p(x, y))
    return w
end

-- settings of text
local FONT = {dft = 0, txt = 0}
for k, v in pairs(FONT) do
    FONT[k] = "ui/font/"..k..".ttf"
end

local function setText(w, param)
    local args = clone(param)
    local config = w:getTTFConfig()
    if args.font then
        config.fontFilePath = FONT[args.font]
    end
    if args.size then
        config.fontSize = args.size
    end
    w:setTTFConfig(config)
    if args.area then
        w:setDimensions(xxui.size(args.area))
    end
    local function f(tbl, k)
        if type(tbl) ~= "table" then
            xx.error(k..": a table needed!", 0)
        end
        local c4b = xxui.color(tbl.color, tbl.opacity)
        local key = {shadow = "offset", outline = "size"}
        local val = {shadow = cc.size(2, -2), outline = 2}
        return c4b, tbl[key[k]] or val[k]
    end
    -- args.shadow:
    -- {color = string/c4b/c3b, opacity = num, offset = cc.size()}
    if args.shadow then
        w:enableShadow(f(args.shadow, "shadow"))
    end
    -- args.outline:
    -- {color = string/c4b/c3b, opacity = num, size = num}
    if args.outline then
        w:enableOutline(f(args.outline, "outline"))
    end
    -- label:setMaxLineWidth(pixels) --> 文本的最大行宽，可用来强制换行
    -- label:disableEffect() --> 取消所有效果
end


-- set: 设置控件的基本属性
-- 坑：控件空创建时不可提前设置 scale/scale9 参数
function set(w, param)
    if not param then return w end
    local args = clone(param)
    if args.node then w:addTo(args.node) end
    if args.name then w:setName(args.name) end
    if args.tag then w:setTag(args.tag) end
    if args.zorder then
        w:setLocalZOrder(args.zorder)
    end
    if args.flip then  -- args.flip: "x" or "y"
        w:setFlippedX(args.flip ~= "y")
        w:setFlippedY(args.flip == "y")
    end
    if args.color then
        setcolor(w, args.color)
    end
    scale(w, args.scale)
    scaleNine(w, args.scale9)
    setpos(w, args)
    if "txt" == w.type_ then
        setText(w, args)
    end
    if args.func then
        w:setTouchEvent(args.func)
    end
    return w
end

-- setStateEvent: func is a function table {enter =, exit =}
function setStateEvent(w, func)
    func = func or {}
    w.statefunc = func
    w:registerScriptHandler(function(event)
        if event == "enter" then
            if func.enter then func.enter() end
            w:onEnter()
        elseif event == "exit" then
            if func.exit then func.exit() end
            w:onExit()
        end
    end)
end

-- addStateEvent: func is a function table {enter =, exit =}
function addStateEvent(w, func)
    if not func then return end
    func = xx.funcAppend(w.statefunc, func)
    setStateEvent(w, func)
end
