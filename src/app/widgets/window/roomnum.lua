module(..., package.seeall)

roomnum = class("roomnum", mm.window)

local margin = 20

function roomnum:ctor(stage, func)
	roomnum.super.ctor(self, stage, 1, true)
	self.func = func
	self.number = ""
	self.current = 0
	self.numposition = {}
	self:init()
	self:setevent()
	self:setTitle("roomnum")
end

-- interface --
function roomnum:getnumber()
	local n = self.number
	local l = string.len(n)
	if l ~= 6 then
		return
	end
	return tonumber(n)
end

function roomnum:action(i)
	if i < 1 or i > 12 then
		return
	end
	if i == 10 then
		self:cleannum()
	elseif i == 12 then
		self:decnum()
	else
		self:addnum(i)
	end
end

function roomnum:addnum(i)
	if self.current >= 6 then
		return
	end
	self.current = self.current + 1
	local idx = self.current
	i = i == 11 and 0 or i
	self:setinput(idx, i)
	self.number = self.number .. i
	if self.current >= 6 and self.func then
		local num = self:getnumber()
		self.func(num)
	end
end

function roomnum:decnum()
	if self.current <= 0 then
		return
	end
	local idx = self.current
	self:setinput(idx)
	self.number = string.sub(self.number, 1, -2)
	self.current = self.current - 1
end

function roomnum:cleannum()
	for idx = 1, self.current do
		self:setinput(idx)
	end
	self.number = ""
	self.current = 0
end

-- implementation --
function roomnum:setevent()
	local board = self:getchild("board")
	local offset = cc.p(margin, margin)
	self.initpos = board:convertToWorldSpace(offset)
	local i, ok
	board:setTouchEvent {
		began = function(evt)
			i = self:numidx(evt)
			if not i then return end
			ok = true
			self:highlight(i, ok)
		end,
		moved = function(evt)
			local j = self:numidx(evt)
			ok = i == j
			self:highlight(i, ok)
	    end,
	    ended = function(evt)
	    	self:highlight(i, false)
	    	if ok then
	    		self:action(i)
	    	end
	    end,
	    canceled = function()
	    	self:highlight(i, false)
		end
	}
end

function roomnum:setinput(idx, num)
	local w = self:getleaf("input_"..idx)
	w:load(num or "")
end

function roomnum:highlight(i, bool)
	if not i then return end
	self.hl:setVisible(bool)
	local text = self:getleaf("num_"..i)
	local c = cc.c3b(70, 150, 210)
	if bool then
		c = cc.c3b(160, 200, 244)
	end
	text:set {color = c}
	if not i then return end
	local pos = self:numpos(i)
	self.hl:set {pos = pos}
end

function roomnum:numpos(i)
	if self.numposition[i] then
		return self.numposition[i]
	end
	local board = self:getchild("board")
	local col = i%3 == 0 and 3 or i%3
	local row = math.ceil(i/3)
	local x = board:align("x", 3, col, 0, margin)
	local y = board:align("y", 4, row, 0, margin)
	local pos = cc.p(x, y)
	self.numposition[i] = pos
	return pos
end

function roomnum:numidx(pos)
	local p = self.initpos
	local x = pos.x - p.x
	local y = pos.y - p.y
	local wid, hei = self:gridsize()
	local col = math.ceil(x/wid)
	local row = 5 - math.ceil(y/hei)
	local i = (row-1)*3+col
	if i >= 1 and i <= 12 then
		return i
	end
end

-- create widgets --
function roomnum:init()
	local board = xxui.create {
		node = self, img = xxres.panel("num"), name = "board",
		anch = cc.p(0.5, 0.5), align = cc.p(0.5, 0.39)
	}
	xxui.create {
		node = self, txt = xx.translate("roomnum", "tips"),
		size = 40, color = cc.c3b(70, 150, 210),
		anch = cc.p(0.5, 1), align = cc.p(0.5, 0.85)
	}
	for i = 1, 6 do
		self:newinput(board, i)
	end
	for i = 1, 12 do
		self:newnum(board, i)
	end
	self.hl = xxui.create {
		node = board, img = xxres.grid("num_hl"),
		anch = cc.p(0.5, 0.5) 
	}
	self.hl:setVisible(false)
end

function roomnum:newinput(node, i)
	local d = 140
	local x = (i-3.5) * d
	local bar = xxui.create {
		node = node, img = xxres.grid("num"),
		anch = cc.p(0.5, 0), align = cc.p(0.5, 1),
		pos = cc.p(x, -20)
	}
	xxui.create {
		node = bar, txt = "", size = 60, name = "input_"..i,
		anch = cc.p(0.5, 0), align = cc.p(0.5, 0.5)
	}
end

function roomnum:newnum(node, i)
	local pos = self:numpos(i)
	local tbl = {
		[10] = "clean", [11] = 0, [12] = "delete"
	}
	local txt = tbl[i] or i
	if type(txt) == "string" then
		txt = xx.translate(txt, "verb")
	end
	xxui.create {
		node = node, txt = txt, name = "num_"..i,
		size = 60, color = cc.c3b(70, 150, 210),
		anch = cc.p(0.5, 0.5), pos = pos
	}
end

function roomnum:gridsize()
	local board = self:getchild("board")
	local wid, hei = board:getsize(2)
	local w = (wid - 2*margin) / 3
	local h = (hei - 2*margin) / 4
	return w, h
end