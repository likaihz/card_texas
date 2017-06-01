-- xxui.ProgressTimer --
module(..., package.seeall)

ProgressTimer = class("ProgressTimer", function()
    return xxui.Widget.new()
end)
--[[
param: {
    (kind = "bar" or "ring", dir)
    (bg = "string" or "Widget", pct)
    content = progress body: "string" or Sprite, 
}
]]
function ProgressTimer:ctor(param)
    local args = clone(param)
    local sp = args.content
    if type(sp) == "string" then
        sp = cc.Sprite:create(sp)
    end
    local progress = cc.ProgressTimer:create(sp)
    local kind = args.kind or "bar"
    local TYPE = {
        bar = cc.PROGRESS_TIMER_TYPE_BAR,
        ring = cc.PROGRESS_TIMER_TYPE_RADIAL
    }
    progress:setType(TYPE[kind] or cc.PROGRESS_TIMER_TYPE_RADIAL)
    if kind == "bar" then
        -- dir: progress bar moving direction
        local direction = {
            left = cc.p(1, 0.5), right = cc.p(0, 0.5),
            up = cc.p(0.5, 0), down = cc.p(0.5, 1)
        }
        local dir = args.dir or "right"
        progress:setMidpoint(direction[dir])
        if dir == "right" or dir == "left" then
            progress:setBarChangeRate(cc.p(1, 0))
        else
            progress:setBarChangeRate(cc.p(0, 1))
        end
    else -- circle progress, defalut:clockwise,
        if args.dir == "anticlockwise" then  
            progress:setReverseDirection(true) -- anticlockwise
        end
    end
    local bg = args.bg or args.img
    if bg then
        if type(bg) == "string" then
            bg = xxui.create { img = bg }
        end
        self:addChild(bg)
        xxui.setsize(self, bg)
    else
        xxui.setsize(self, sp)
    end
    self:setAnchorPoint(0, 0)
    xxui.set(self, args)
    xxui.set(progress, {
        node = self, anch = cc.p(0.5, 0.5),
        align = args.location or cc.p(0.5, 0.5), pos = args.offset
    })
    progress:setPercentage(args.pct or 100)
    self.core = progress
end

function ProgressTimer:progressTo(time, percent, callback)
    local act = cc.ProgressTo:create(time, percent)
    if callback then
        act = cc.Sequence:create(act, cc.CallFunc:create(callback))
    end
    self.core:runAction(act)
end

function ProgressTimer:progressFromTo(time, from, to, callback)
    local act = cc.ProgressFromTo:create(time, from, to)
    if callback then
        act = cc.Sequence:create(act, cc.CallFunc:create(callback))
    end
    self.core:runAction(act)
end

function ProgressTimer:setPercent(pct)
    self.core:setPercentage(pct)
end

function ProgressTimer:setpercent(num)
    self.core:setPercentage(num)
end

function ProgressTimer:load(num)
    num = num or 0
    num = tonumber(num)
    num = (num > 1) and 1 or num
    num = (num < 0) and 0 or num
    self.core:setPercentage(num * 100)
end