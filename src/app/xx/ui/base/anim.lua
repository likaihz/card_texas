module(..., package.seeall)

-- delay: 延时执行函数，func 为必传参数
function delay(w, time, func)
    func = func or function() end
    if time <= 0 then return func() end
    local seq = cc.Sequence:create(
        cc.DelayTime:create(time),
        cc.CallFunc:create(func)
    )
    return w:runAction(seq)
end

function schedule(w, time, num, func, oncomplete)
    func = func or function() end
    local seq = cc.Sequence:create(
    	cc.CallFunc:create(func),
        cc.DelayTime:create(time)
    )
    local act
    if num > 0 then
        act = cc.Repeat:create(seq, num)
    else
        act = cc.RepeatForever:create(seq)
    end
    if oncomplete then
        act = cc.Sequence:create(act, cc.CallFunc:create(oncomplete))
    end
    return w:runAction(act)
end

-- implementation of fade --
local function callback(w, delete, func)
	return function()
        if func then func() end
        if delete then w:removeSelf() end
    end
end

-- fadeout: 渐变淡出
function fadeout(w, time, delete, func)
    local children = w:getChildren()
    for k, v in pairs(children or {}) do
        fadeout(v, time)
    end
    transition.fadeOut(w, {
    	time = time, onComplete = callback(w, delete, func)
    })
end

-- fadeout: 渐变淡入
function fadein(w, time, delete, func)
    local children = w:getChildren()
    for k, v in pairs(children or {}) do
        fadein(v, time)
    end
    transition.fadeIn(w, {
    	time = time, onComplete = callback(w, delete, func)
    })
end

-- moveto: 移动到给定位置
function moveto(w, time, pos, func)
    local seq = cc.Sequence:create(
        cc.MoveTo:create(time, pos),
        cc.CallFunc:create(function()
            if func then func() end
        end))
    return w:runAction(seq)
end

-- moveby: 移动给定位移
function moveby(w, time, pos, func)
    local seq = cc.Sequence:create(
        cc.MoveBy:create(time, pos),
        cc.CallFunc:create(function()
            if func then func() end
        end))
    return w:runAction(seq)
end

-- vibrate: 控件振动
-- mode = "x": horizontal, "y": vertical, "#": cross, "@": rotation
-- time 单次振动时间，range 振动幅度，num 振动次数，func 振动结束后回调
local function duplicate(tbl, num)
    local t = {}
    for i = 1, num do
        for i, v in ipairs(tbl) do
            table.insert(t, v)
        end
    end
    return t
end

function vibrate(w, mode, time, range, num, func)
    if num == 0 then
        if func then func() end; return
    end
    local t = (mode == "#") and time/8 or time/4
    local function m(mod, p)
        if mod == "@" then
            return cc.RotateBy:create(math.abs(p) * t, p * range)
        end
        local pos = cc.p(0, 0)
        pos[mod] = p * range
        return cc.MoveBy:create(math.abs(p) * t, pos)
    end
    local tbl = {}
    local function period(mod)
        table.insert(tbl, m(mod, 1))
        table.insert(tbl, m(mod, -2))
        table.insert(tbl, m(mod, 1))
    end
    if mode == "#" then
        period("x"); period("y")
    else
        period(mode)
    end
    local anch = w:getAnchorPoint()
    local x, y = w:getPosition()
    local rot = w:getRotation()
    if mode == "@" then
        w:setAnchorPoint(cc.p(0.5, 0.5))
        local wid, hei = xxui.getsize(w, 2)
        w:setPosition(x + (0.5-anch.x)*wid, y + (0.5-anch.y)*hei)
    end
    local f = function()
        w:setAnchorPoint(anch)
        w:setPosition(x, y)
        w:setRotation(rot)
    end
    local call = (not func) and f or function()
        f(); func()
    end
    local seq
    if num > 0 then
        tbl = duplicate(tbl, num)
        table.insert(tbl, cc.CallFunc:create(call))
        seq = cc.Sequence:create(unpack(tbl))
    else
        seq = cc.Sequence:create(unpack(tbl))
        seq = cc.RepeatForever:create(seq)
    end
    return w:runAction(seq)
end