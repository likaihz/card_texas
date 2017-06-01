-- xxui.LoadingBar --
module(..., package.seeall)

LoadingBar = class("LoadingBar", function()
    return xxui.Widget.new()
end)
--[[
---- SUPPORT SCALE9 ZOOM ----
param: {
    (dir = "loading bar moving direction", pct),
    (bg = prepared outside, offset = cc.p, edge),
    content = progress body: prepared outside with 100% length
}
]]
function LoadingBar:ctor(param)
    local args = clone(param)
    local bar = args.content
    bar:setScale9Enabled(true)
    local bg = args.bg
    if bg then
        self:addChild(bg)
        self:setsize(bg)
    else
        self:setsize(bar)
    end
    self:set(args)
    -- dir: loading bar moving direction
    local dir = args.dir or "right"
    local anchs = {
        left = cc.p(1, 0), right = cc.p(0, 0),
        up = cc.p(0, 0), down = cc.p(0, 1)
    }
    local anch = anchs[dir]
    local anch0 = bar:getAnchorPoint()
    local da = cc.pSub(anch, anch0)
    local w, h = self:getsize(2)
    local align = args.location or cc.p(0, 0)
    local pos = cc.p(w*align.x, h*align.y)
    local offset = args.offset or cc.p(0, 0)
    pos = cc.pAdd(pos, offset) -- init Pos of anch0
    local wid, hei = bar:getsize(2)
    local dp = cc.p(wid*da.x, hei*da.y)
    pos = cc.pAdd(pos, dp) -- new Pos of new anch
    bar:set{
        node = self, anch = anch, pos = pos
    }
    self.edge = args.edge or 20
    if dir == "left" or dir == "right" then
        self.dir = "width"
        self.full = wid - self.edge
    else
        self.dir = "height"
        self.full = hei - self.edge
    end
    self.core = bar
    self:setPercent(args.pct or 0)
end

function LoadingBar:load(num)
    num = num or 0
    num = tonumber(num)
    num = (num > 1) and 1 or num
    num = (num < 0) and 0 or num
    local bar, dir = self.core, self.dir
    local l = self.full*num + self.edge
    local l0 = bar:getVirtualRendererSize()[dir]
    local c = (l >= l0) and 0.2 or 1 - l/l0
    local cap = bar:getCapInsets()
    local p = (dir == "width") and "x" or "y"
    cap[p], cap[dir] = l0 * (0.5 - c/2), l0 * c
    bar:setCapInsets(cap)
    local sz = bar:getContentSize()
    sz[dir] = l
    bar:setContentSize(sz)
end

function LoadingBar:setPercent(pct)
    pct = tonumber(pct)
    if not pct or pct > 100 or pct < 0 then
        xx.error("invalid percent input", 0)
    end
    self:load(pct/100)
end