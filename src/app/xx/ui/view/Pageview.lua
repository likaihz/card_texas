-- xxui.Pageview --
module(..., package.seeall)

Pageview = class("Pageview", xxui.Scrollview)
-- param = { 
-- 		node, size, (dir = "x" or "y"), (bounce = boolean or number),
--		(range = number < 1 -- when to turn page)
--		(zorder = "avoid to cover other widgets")
--}
function Pageview:ctor(param)
	self.super.ctor(self, param)
	local a = (self.dir == "y") and 1 or 0
	self.content = xxui.Widget.new {
		node = self.clip, align = cc.p(0, a)
	}
    self.initpos = self:getLocation()
	local args = clone(param)
	self.bounce = self:param(args.bounce, 1/2)
	self.turn_range = args.range or 0.4
	self.pages = {}
	self.length = 0
	self.cell = table.maxn(self.pages) + 1 -- special for Pageview
	self.curPageIdx = 1
	self:setTouchEnabled(false)
end

-- interface --
function Pageview:newPage(i)
	if not i then
		i = 1
		while self.pages[i] do
			i = i + 1
		end
	elseif i < 1 or math.ceil(i) ~= i then
		xx.error(i.." must be a positive integer!!", 0)
	elseif self.pages[i] then
		xx.error("page"..i.."already exists!!", 0)
	end
	self.pages[i] = 0
	local page = xxui.Widget.new()
	page:setsize(self)
	page:set {
		node = self.content, pos = self:calcpos(i), tag = i
	}
	self:update()
	if i == 1 then
		self:setTouchEnabled(true)
	end
	return page
end

function Pageview:getCurPageIdx()
	return self.curPageIdx
end

function Pageview:setPage(i)
	if i < 1 or i > table.maxn(self.pages) then
		xx.error("page index "..i.." overflow!!", 0)
	elseif math.ceil(i) ~= i then
		xx.error(i.." must be an integer!!", 0)
	end
	local pos = self:calcpos(i)
	pos[self.dir] = -pos[self.dir]
	self.content:set {pos = pos}
	self.curPageIdx = i
end

function Pageview:remove(i)
	self.content:getchild(i):removeSelf()
	self.pages[i] = nil
	if i > table.maxn(self.pages) then
		self:update()
	end
end

function Pageview:rearrange()
	local tbl = {}
	for index in pairs(self.pages) do
		tbl[#tbl+1] = index
	end
	table.sort(tbl)
	self.pages = {}
	for i, index in ipairs(tbl) do
		if i ~= index then
			local page = self.content:getchild(index)
			page:set { pos = self:calcpos(i) }
		end
		self.pages[i] = 0
	end
	self:update()
	self:setPage(1)
end

-- function Pageview:setTouchEvent(func)
-- 	local f
-- 	if type(func) == "table" then
-- 		f = func.ended
-- 		if f then 
-- 			func.ended = function()
-- 				f(self.curPageIdx)
-- 			end
-- 		end
-- 		f = func
-- 	end
-- 	if type(func) == "function" then
-- 		f = {ended = function() func(self.curPageIdx) end}
-- 	end
-- 	self.func = f or {}
-- end

function Pageview:regularize(func)
	func = func or {}
	if type(func) == "function" then
		func = {ended = func}
	end
	local f = func.ended
	if f then
		func.ended = function()
			f(self.curPageIdx)
		end
	end
	return func
end

-- func: nil, function or table of function
function Pageview:setTouchEvent(func)
	self.func = self:regularize(func)
end

function Pageview:addTouchEvent(func)
	func = self:regularize(func)
	self.func = xx.funcAppend(self.func, func)
end


-- override functions --
function Pageview:setTouchEnabled(bool)
	self.canTouch = bool
	if bool then
		local i = self:getLocation() / self:getlen()
		self.curPageIdx = math.round(i)
		if self.dir == "x" then
			self.curPageIdx = -self.curPageIdx+1
		end
	end
end

-- implementation --
function Pageview:calcpos(i)
	local pos = cc.p(0, 0)
	local w, h = self:getsize(2)
	if self.dir == "x" then
		pos.x = (i - 1) * w
	else
		pos.y = - i * h
	end
	return pos
end

function Pageview:update()
	self.length = self:getlen() * table.maxn(self.pages)
	self.cell = table.maxn(self.pages) + 1
end