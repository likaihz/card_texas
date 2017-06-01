module(..., package.seeall)

function screen()
    return cc.size(display.width, display.height)
end

-- imgpath: 返回一个函数，该函数在路径 path 中获取图片
function imgpath(path, ext)
    local ext = ext or ".png"
    return function(name)
        return table.concat {path, name, ext}
    end
end

-- pos: 输入位置参数预处理 (return x, y)
-- input: cc.p(10, 20), {x = 10, y = 20}
-- input: {10, 20} or (10, 20)
function pos(p1, p2)
	if type(p1) == "number" and p2 then
        return p1, p2
    end
    if type(p1) ~= "table" then xx.error(p1) end
    local x = p1.x or p1[1]
    local y = p1.y or p1[2]
    return x, y
end

-- size: 输入大小参数预处理 (return wid, hei)
-- input: cc.size(10, 20), {width = 10, height = 20}
-- input: {10, 20} or (10, 20)
function size(p1, p2)
    if type(p1) == "number" and p2 then
        return p1, p2
    end
    if type(p1) ~= "table" then xx.error(p1) end
    local wid = p1.width or p1[1]
    local hei = p1.height or p1[2]
    return wid, hei
end

-- implementation of color --
local COLOR = {
    white = cc.c3b(255, 255, 255),
    black = cc.c3b(0, 0, 0),
    blue = cc.c3b(54, 93, 153),
    gray = cc.c3b(75, 75, 75),
    --
    BLACK1 = cc.c3b(119, 90, 43),
    BLACK2 = cc.c3b(99, 50, 21),
    DARK = cc.c3b(165, 137, 96),
}

-- 将表示颜色的参数转换为 cc.c3b/cc.c4b
local function cc_cxb(p)
    local c = clone(p)
    if type(c) == "string" then
        c = COLOR[c]
    end
    if type(c) ~= "table" then
        xx.error(c)
    end
    if c.r and c.g and c.b then
        local x = c.a and 4 or 3
        return c, x
    else
        xx.error("bad color input!", 0)
    end
end

-- color: 统一色值接口，返回 cc.c4b
-- color: string/c3b/c4b, opacity: number
function color(c, opacity)
    local c, x = cc_cxb(c)
    if x == 4 then return c end
    opacity = opacity or 255
    if opacity <= 1 then
    	opacity = math.floor(255 * opacity)
    end
    c.a = opacity
    return c
end