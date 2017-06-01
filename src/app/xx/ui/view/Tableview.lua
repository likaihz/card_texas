-- xxui.Tableview --
module(..., package.seeall)

local DIRECTION = {
    x = cc.SCROLLVIEW_DIRECTION_HORIZONTAL,
    y = cc.SCROLLVIEW_DIRECTION_VERTICAL
}
local TABLECELL = {
    size = cc.TABLECELL_SIZE_FOR_INDEX,
    number = cc.NUMBER_OF_CELLS_IN_TABLEVIEW,
    create = cc.TABLECELL_SIZE_AT_INDEX,
    event = cc.TABLECELL_TOUCHED
}
-- class Tableview --
-- use args.setcell (param) or setCell() (interface) to create cell
-- use args.click (param) or setClick() (interface) to set cell touched event
Tableview = class("Tableview", function()
    return xxui.Widget.new()
end)
-- args: {
--     data = {} or nil, num = 3.5, dir = "y", size = cc.size(),
--     setcell = func, click = func, fade = bool, bounce = bool,
--     anch, align, pos ...
-- }
function Tableview:ctor(args)
    self.data = args.data
    -- initial param
    local d = args.dir or "y"
    local k = d == "x" and "width" or "height"
    local sz = clone(args.size)
    sz[k] = sz[k] / args.num
    self.cellsize_ = sz
    self.dir_, self.len_ = d, sz[k]
    self.viewnum = args.num
    -- create tableview
    self.core_ = self:init(args)
    self.initpos = self:getpos()
    self:setCell(args.setcell)
    self:setClick(args.click)
    -- set fade --
    self.fade = args.fade
    if self.fade then
        self:setfade()
    end
end

function Tableview:onExit()
    self.sch_ = xx.unschedule(self.sch_)
end

-- interface --
function Tableview:load(data)
    if data then
        if type(data) ~= "table" then
            xx.error(data)
        end
        self.data = data
    end
    self.core_:reloadData()
    if self.fade then
        self:setopac()
    end
end

function Tableview:getdata(index)
    local data = self.data
    if not index then
        return data
    end
    return data and data[index]
end

function Tableview:getnum()
    local data = self.data
    return data and #data or 0
end

function Tableview:getpos()
    local w = self.core_
    local pos = w:getContentOffset()
    return pos[self.dir_]
end

function Tableview:getlen()
    local w, h = self:getsize(2)
    return (self.dir_ == "x") and w or h
end

function Tableview:setCell(func)
    self.setcell_ = func
    self:load()
end

function Tableview:setClick(func)
    self.click_ = func
end

-- implementation (callbacks) --
local tableCell = {}

function tableCell.size(w, tbl, idx)
    local sz = w.cellsize_
    return sz.height, sz.width
end

function tableCell.number(w, tbl)
    return w:getnum()
end

function tableCell.create(w, tbl, idx)
    local cell = tbl:dequeueCell()
    if not cell then
        cell = cc.TableViewCell:new()
        xxui.setsize(cell, w.cellsize_)
    end
    local data = w:getdata(idx+1)
    if not data then
        cell:setVisible(false)
        return cell
    end
    cell:setVisible(true)
    local f = w.setcell_
    if f then f(cell, idx+1) end
    return cell
end

function tableCell.event(w, tbl, cell)
    local idx = cell:getIdx()
    local f = w.click_
    if f then f(cell, idx+1) end
end

-- create widgets --
function Tableview:init(args)
    self:setsize(args.size)
    self:set(args)
    local w = cc.TableView:create(args.size):addTo(self)
    local bounce = args.bounce
    w:setBounceable(bounce == nil or bounce)
    w:setDirection(DIRECTION[self.dir_])
    w:setDelegate()
    w:setVerticalFillOrder(cc.TABLEVIEW_FILL_TOPDOWN)
    for k, v in pairs(TABLECELL) do
        w:registerScriptHandler(handler(self, tableCell[k]), v)
    end
    return w
end

-- implementation of cell fade in n out --
function Tableview:setfade()
    self.head_ = 0
    self.tail_ = self.viewnum
    local layer = xxui.Widget.new {node = self}
    xxui.setsize(layer, self)
    layer:setTouchEvent {
        moved = function()
            if self.sch_ then return end
            self.sch_ = xx.schedule(function()
                self:sliding()
            end, 0.03)
        end
    }
    layer:setSwallowTouches(false)
end

function Tableview:sliding()
    local pos = self:getpos()
    local stop = pos == self.pre_
    stop = stop and not self:bouncing()
    self.pre_ = pos
    if stop then
        self.sch_ = xx.unschedule(self.sch_)
        return
    end
    self:setopac()
end

-- check if the list is in bouncing region
function Tableview:bouncing()
    if self.head_ < 0 then
        return true
    end
    local n = self:getnum()
    if n >= self.viewnum then
        return self.tail_ > n
    else
        return self.head_ > 0
    end
end

-- return: cellIdx, opacity
function Tableview:getopac(num, tail)
    local i = math.floor(num)
    local f = num-i
    local opac = tail and f or 1-f
    return i, opac
end

-- head, tail: head (tail) of view region
function Tableview:setopac()
    local n = self:getnum()
    local p = self:getpos() / self.len_
    local v = self.viewnum
    self.head_ = (self.dir_ == "y") and n+p-v or -p
    self.tail_ = self.head_ + v
    local opacity = {}
    local head, opac = self:getopac(self.head_)
    opacity[head] = opac
    local tail, opac = self:getopac(self.tail_, true)
    opacity[tail] = opac
    -- set opacity --
    for i = head, tail do
        opac = opacity[i] or 1
        local cell = self.core_:cellAtIndex(i)
        if cell then
            xxui.setopacity(cell, opac, true)
        end
    end
end

