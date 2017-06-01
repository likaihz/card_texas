local Fighter = require("app.models.battle.Fighter")

local Battle = class("Battle")

function Battle:ctor(scene)
	self.scene = scene
	self:init()
	self.roundcnt = 0
end

-- interface --
-- k: "scene", "fighter", "roundcnt"
function Battle:get(k)
	return self[k]
end

function Battle:receive(msg)
	local func = {
		seat = 0, ready = 0, start = 0,
		card = 0, cards = 0, raise = 0, result = 0
	}
	local opt = msg.opt
	if not func[opt] then return end
	msg.opt = nil
	local k = "load"..opt
	self[k](self, msg)
end

-- opt: "enter", "ready", "leave", "raise"
function Battle:send(opt, raise)
	local msg = {opt = opt, raise = raise}
	if raise then
		msg.idx = self.me
	end
	GM:send(msg)
end

-- load: call by server --
function Battle:loadseat(msg)
	local data = msg.data
	if not self.fighter then
		self.me = self:myseat(data)
		self.fighter = self:createfighter()
	end
	for i, ft in pairs(self.fighter) do
		local e = data[tostring(i)]
		ft:load(e)
	end
end

function Battle:loadready(msg)
	local idx = tonumber(msg.idx)
	if self.me == idx then
		self.scene:removechild("ready")
	end
	local ft = self.fighter[idx]
	ft:ready()
end

function Battle:loadstart(msg)
	self.roundcnt = msg.round
	local scene = self.scene
	scene:load("roundnum")
	print("game start!")  -- animation
	xxui.delay(scene, 0.3, function()
		self:clean()
		self:setdealer(msg.dealer)
		xxui.delay(scene, 0.3, function()
			self:deal()
		end)
	end)
end

function Battle:loadcards(msg)
	local me = self:getme()
	self.pokernum = me:loadcards(msg.data)
end

function Battle:loadraise(msg)
	local ft = self.fighter[msg.idx]
	ft:showraise(msg.raise)
end

function Battle:loadcard(msg)
	local me = self:getme()
	me:loadcard(msg.data)
	self:addpoker()
end

function Battle:loadresult(msg)
	for s, data in pairs(msg.data) do
		local i = tonumber(s)
		local ft = self.fighter[i]
		ft:cache(data)
	end
end

-- procedure --
function Battle:deal()
	self:polling(0.25, function(ft)
		ft:obtain(self.pokernum)
	end, function()
		self:chooseraise()
	end)
end

function Battle:addpoker()
	self:polling(0.25, function(ft)
		ft:obtainone()
	end, function()
		self:checkrank()
		self.scene:countdown(6, function()
			self.rank:setVisible(false)
			self:showpoker()
		end)
	end)
end

function Battle:showpoker()
	local me = self:getme()
	me:coverpoker()
	self:polling(0.25, function(ft)
		ft:showpoker(true)
		ft:showrank()
	end, function()
		self:showscore()
		self:showresult()
	end)
end

function Battle:showscore()
	for i, ft in pairs(self:fighters()) do
		ft:showscore()
	end
end

function Battle:showresult()
	local result = mm.result.new(self.scene)
	result:load(self)
	result:onclose(function()
		self:send("ready")
		if self:ended() then
			GM:popscene()
		end
	end)
end

-- implementation --
function Battle:ended()
	local roundnum = self.scene:get("roundnum")
	return self.roundcnt >= roundnum
end

function Battle:clean()
	for _, ft in pairs(self:fighters()) do
		ft:seticon()
	end
end

function Battle:polling(time, func, oncomplete)
	local co
	co = coroutine.create(function()
		for i, ft in pairs(self:fighters()) do
			if func then func(ft) end
			ft:delay(time, function()
				coroutine.resume(co)
			end)
			coroutine.yield()
		end
		if oncomplete then oncomplete() end
	end)
	coroutine.resume(co)
end

function Battle:fighters()
	local tbl = {}
	for i, ft in pairs(self.fighter) do
		if ft:present() then
			tbl[i] = ft
		end
	end
	return tbl
end

function Battle:myseat(data)
	local name = GM.player:get("name")
	for s, v in pairs(data) do
		if v.name == name then
			return tonumber(s)
		end
	end
end

function Battle:arange()
	local me = self:getme()
	-- me:loadpokers()
	me:showpoker(true)
end

function Battle:getme()
	return self.fighter[self.me]
end

function Battle:setdealer(idx)
	if self.dealer == idx then
		return
	end
	local i = self.dealer
	if i then
		local old = self.fighter[i]
		old:showdealer(false)
	end
	local new = self.fighter[idx]
	new:showdealer(true)
	self.dealer = idx
end

function Battle:chooseraise()
	local me = self:getme()
	if me:is("dealer") then
		return
	end
	self.raise:load(true)
end

function Battle:checkrank()
	local btn = self.rank
	btn:setVisible(true)
	btn:setTouchEnabled(true)
	local me = self:getme()
	local rank = me:get("rank")
	local txt = xx.translate(rank, "niuniu")
	btn:loadtxt(txt)
	if rank <= 0 then
		btn:loadimg("gray", true)
	elseif rank < 10 then
		btn:loadimg("blue", true)
	else
		btn:loadimg("orange", true)
	end
end

function Battle:fighternum()
	return 5
end

-- create widgets --
function Battle:init()
	self.raise = mm.raise.new(self.scene)
	self.raise:setevent(function(i)
		self:send("raise", i)
	end)
	self.raise:setcomplete(function()
		self:send("raise", 1)
	end)
	self.rank = self:createrank()
end

function Battle:createfighter()
	local num = self:fighternum()
	local tbl = {}
	for i = 0, num-1 do
		tbl[i] = Fighter.new(self, i)
	end
	return tbl
end

function Battle:createrank()
	local btn
	btn = xxui.Txtbtn.new {
		node = self.scene, mode = "blue", txt = "",
		anch = cc.p(0.5, 0.5), align = cc.p(0.9, 0.1),
		func = function()
			self:arange()
			btn:setTouchEnabled(false)
		end
	}
	btn:setVisible(false)
	return btn
end

return Battle