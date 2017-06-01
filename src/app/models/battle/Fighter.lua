-- class Fighter --
local Fighter = class("Fighter", function()
	return xxui.Widget.new()
end)

function Fighter:ctor(battle, idx)
	self.battle = battle
	self.scene = battle:get("scene")
	self.idx = idx
	self.data = {}
	self:init()
end

-- interface --
function Fighter:get(k)
	if self[k] then
		return self[k]
	end
	return self.data[k]
end

-- k: "me", "dealer"
function Fighter:is(k)
	local i = self.battle:get(k)
	return self.idx == i
end

function Fighter:load(data)
	local have = data and true or false
	self:setVisible(have)
	if self:present() and not have then
		self:leave()
	end
	self:cache(data)
	if not have then
		return
	end
	self:showname()
	self:showscore()
end

function Fighter:cache(data)
	data = data or {}
	local new = data.score or 0
	local old = self:get("score") or 0
	self.addscore = new - old
	self.data = data
	self:loadcards(data.cards)
end

function Fighter:present()
	local name = self:get("name")
	return name and true or false
end

function Fighter:leave()
	self:cleanpoker()
	-- tips --
end

-- load pokers' data in fighter
function Fighter:loadcard(data)
	table.insert(self.data.cards, data)
end

function Fighter:loadcards(data)
	local tbl = {}
	for s, v in pairs(data or {}) do
		local i = tonumber(s) + 1
		tbl[i] = v
	end
	self.data.cards = tbl
	return table.len(tbl)
end

function Fighter:newpoker(i)
	local poker = mm.poker.new(self.scene)
	table.insert(self.poker, poker)
	self:loadpoker(i)
	return poker
end

-- load data in poker
function Fighter:loadpoker(i)
	local data = self:get("cards")
	i = i or #data
	local poker = self.poker[i]
	poker:load(data[i])
end

function Fighter:loadpokers()
	for i, poker in pairs(self.poker) do
		self:loadpoker(i)
	end
end

function Fighter:showpoker(bool)
	self:loadpokers()
	for i, poker in pairs(self.poker) do
		poker:show(bool)
	end
end

function Fighter:coverpoker()
	for i, poker in pairs(self.poker) do
		poker:cover(i)
	end
end

function Fighter:cleanpoker()
	for i, poker in pairs(self.poker or {}) do
		poker:removeSelf()
	end
	self.poker = {}
	self.data.cards = {}
end

function Fighter:ready()
	self:cleanpoker()
	self:showraise()
	local img = "obj/rank/ready.png"
	local pos = self:cardpos()
	self:showicon(img, pos)
end

function Fighter:obtain(num)
	local co
	co = coroutine.create(function()
		for i = 1, num do
			local poker = self:newpoker(i)
			poker:move(self, 1, function()
				if i == num then
					self:slide()
				end
			end)
			self:delay(0.15, function()
				coroutine.resume(co)
			end)
			coroutine.yield()
		end
	end)
	coroutine.resume(co)
end

function Fighter:obtainone()
	local poker = self:newpoker(5)
	poker:move(self, 5, function()
		if self:is("me") then
			poker:show(true, "flip")
		end
	end)
end

function Fighter:slide()
	local mode = self:is("me") and "flip"
	for i, poker in pairs(self.poker) do
		poker:slide(i, mode)
	end
end

function Fighter:showname()
	local name = self:get("name")
	local w = self:getleaf("name")
	w:load(name)
end

function Fighter:showscore()
	local num = self:get("score") or 0
	local color = self:scorecolor(num)
	local w = self:getleaf("score")
	w:load(num)
	w:set {color = color}
end

function Fighter:showdealer(bool)
	local w = self:getleaf("player")
	local img = bool and "dealer" or "player"
	w:load(xxres.panel(img))
	w = self:getleaf("raise")
	w:setVisible(bool)
	if bool then
		w:load(xxres.icon("dealer"))
	end
end

function Fighter:showraise(num)
	local icon = self:getchild("raise")
	local bool = num and true or false
	icon:setVisible(bool)
	if num then
		icon:load(xxres.icon("x"..num))
	end
end
function Fighter:showrank(num)
	local num = self:get("rank")
	if num == 0 then
		for i, poker in pairs(self.poker) do
			poker:setlight(0.5)
		end
	end
	local img = "obj/rank/rank_"..num..".png"
	local pos = self:cardpos(true)
	pos.y = pos.y-40
	self:showicon(img, pos)
end

-- w: img or nil
function Fighter:seticon(w)
	if not self.icon then
		self.icon = w
		return
	end
	self.icon:removeSelf()
	self.icon = w
end

function Fighter:scorecolor(num)
	local color = cc.c3b(220, 100, 190)
	if num > 0 then
		color = cc.c3b(220, 174, 80)
	end
	return color
end

-- implementation --
function Fighter:showicon(img, pos, time)
	local w = xxui.create {
		node = self.scene, img = img, scale = 3.5,
		anch = cc.p(0.5, 0.5), pos = pos
	}
	time = time or 0.15
	local a1 = cc.FadeIn:create(time)
	local a2 = cc.ScaleTo:create(time, 1)
	w:setOpacity(0)
	w:runAction(a1)
	w:runAction(a2)
	self:seticon(w)
end

-- create widgets --
function Fighter:init()
	local panel = xxui.create {
		node = self, img = xxres.panel("player"), name = "player"
	}
	self:setsize(panel)
	self:set {
		node = self.scene, anch = cc.p(0.5, 0.5),
		align = self:getalign()
	}
	self:createinfo()
end

local ALIGN = {
	cc.p(0.1, 0.1), cc.p(0.9, 0.6), cc.p(0.65, 0.9),
	cc.p(0.35, 0.9), cc.p(0.1, 0.6)
}
function Fighter:getalign()
	local me = self.battle:get("me")
	local i = self.idx - me
	i = i >= 0 and i or i+5
	return ALIGN[i+1]
end

function Fighter:cardpos(cover)
	local x, y
	if self:is("me") then
		x = display.cx
		y = cover and 220 or 100
		return cc.p(x, y)
	end
	x, y = self:getPosition()
	return cc.p(x, y-160)
end

function Fighter:createinfo()
	xxui.create {
		node = self, txt = "玩家名字", name = "name",
		color = cc.c3b(127, 176, 234), size = 28,
		anch = cc.p(0, 0.5), align = cc.p(0.5, 0.63)
	}
	xxui.create {
		node = self, txt = "+6666", name = "score", size = 28, 
		anch = cc.p(0, 0.5), align = cc.p(0.5, 0.37)
	}
	local raise = xxui.create {
		node = self, img = xxres.icon("x1"), name = "raise",
		anch = cc.p(1, 1), align = cc.p(1, 1)
	}
	raise:setVisible(false)
end

return Fighter