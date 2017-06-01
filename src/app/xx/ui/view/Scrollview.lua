-- xxui.Scrollview --
module(..., package.seeall)

Scrollview = class("Scrollview", function()
    return xxui.Widget.new()
end)
-- param: {
--     node, (size), (dir = "x" or "y"), (bounce = bool or num), 
--     (content, location, offset),(edge = {head, tail}),
--     length = length of content, cell, (func),
--	   (zorder = "avoid to cover other widgets")
-- }
function Scrollview:ctor(param)
	local args = clone(param)
	self.clip = cc.ClippingRectangleNode:create():addTo(self)
	self:setTouchEvent(args.func)
	args.func = nil
    self.dir = args.dir or "x"
    self.cell = args.cell
    self.speed = args.speed or 1000
    self.touchLayer = self:initTouch()   
    -- set size/pos/content/touchLayerPos
    self:setsize(args.size or xxui.screen())
    self:set(args)
    local edge = args.edge or {}
    self.head = self:param(edge[1], 0)
    self.tail = self:param(edge[2], 1)
    self.turn_range = args.range or 0.4
	self.bounce = self:param(args.bounce, 1/6)
	self.slide = 0 --slide speed (can be > 0 or < 0)
end

function Scrollview:onExit()
	xx.unschedule(self.scheduler)
end

-- interface --
function Scrollview:set(args)
	if args.size then
		self:setsize(args.size)
	end
	local wid, hei = self:getsize(2)
	self.clip:setClippingRegion(cc.rect(0, 0, wid, hei))
	xxui.setsize(self.clip, cc.size(wid, hei))
	xxui.set(self, args)
	-- update touchLayer pos
	local p = self:convertToWorldSpace(cc.p(0, 0))
	self.touchLayer:set {pos = cc.p(-p.x, -p.y)}
	self.viewRect = cc.rect(p.x, p.y, wid, hei)
end

function Scrollview:reset()
	self.length = self:getLength() 
	self:setLocation()
end

-- args = {content, length, anch, align, pos}
function Scrollview:setContent(args)
	self.length = args.length or self.length
	local w = args.content
	if w then
		if self.content then
			self.content:removeFromParent()
		end
		w:addTo(self.clip)
		self.content = w
	end
	self:regularTouchEvent(self.content)
	self.content:set(args)
	self.initpos = self:getLocation()
end

function Scrollview:getLocation()
	local x, y = self.content:getPosition()
    return (self.dir == "x") and x or y
end

-- p: integer(self.cell isn't nil) or number
function Scrollview:setLocation(p)
	p = p or self.initpos
	if self.cell then
		local len = self.length / (self.cell - 1)
		local sign = self.dir == "x" and -1 or 1
		p = self.initpos + sign * (p-1) * len
	end
	if self.dir == "x" then
		self.content:setPositionX(p)
	else
		self.content:setPositionY(p)
	end
end

function Scrollview:getLength()
	local w, h = self.content:getsize(2)
    return (self.dir == "x") and w or h
end

-- set scrollview's touch event
function Scrollview:setTouchEvent(func)
	if type(func) == "function" then
		func = {moved = func}
	end
	self.func = func or {}
end

-- set layer's child clicked event
function Scrollview:setChildEvent(w, func)
	w:setTouchEventBy(self, func)
end

-- func: callback function
function Scrollview:contentMoveBy(dis, func, speed)
	local pos = self:getLocation()
	local des, r = self:outrange(pos + dis)
	return self:move(des - pos, 2, function()
		if self.bounce and r ~= 0 then
			return self:move(r, 2, func)
		end
		if func then func() end
	end, speed)
end

function Scrollview:contentMoveTo(des, func)
	local pos = self:getLocation()
	return self:contentMoveBy(des-pos, func)
end

-- implementation --
function Scrollview:getlen()
	local w, h = self:getsize(2)
	return (self.dir == "x") and w or h
end

function Scrollview:param(p, dft)
	if p == nil then p = dft end
	local t = type(p)
	if t ~= "boolean" and t ~= "number" then
		xx.error(p)
	end
	if t == "boolean" then
		if not p then return false end
		p = dft
	end
	if p < 0 then
		xx.error("positive param p needed!!", 0)
	end
	return (p < 10) and p * self:getlen() or p
end

function Scrollview:initTouch()
	local cover = xxui.Widget.new {node = self, zorder = 911}
	cover:setsize(xxui.screen())
	cover:setTouchEnabled(true)
	cover:setSwallowTouches(false)
	self.canTouch = true
	local dir = self.dir
	local flag = false
	local r, cp, began, t
	local event
    cover:addTouchEventListener(function(sender, evt)
    	if not self.canTouch and self.slide == 0 then return end
    	if evt == 0 then
    		event = sender:getTouchBeganPosition()
    		event.name = "began"
    		if not flag and self:isInview(event) then
    			flag = true
    			r, cp = 0, self:getLocation()
    			began, t = event, os.clock()
    		end
        elseif evt == 1 then  -- "moved"
        	event = sender:getTouchMovePosition()
        	event.name = "moved"
        	if flag then
	        	if self.slide ~= 0 then return end
	        	local dp = event[dir] - began[dir]
	            r = self:drag(cp + dp)
	        end
        elseif evt == 2 or evt == 3 then -- "ended" or "canceled"
            if flag then
            	flag = false
            	if self.bounce and r ~= 0 then
            		return self:move(r)  -- bounce back
            	end
            	event = sender:getTouchEndPosition()
            	event.name = "ended"
	            local dp = cc.pSub(event, began)
	            local dt = os.clock() - t
	            -- click to stop sliding
	            if self.slide ~= 0 and self:isClick(dp, dt) then
	            	xx.unschedule(self.scheduler)
	            	self.content:stopAllActions()
	            	self.slide = 0
	            	self:setTouchEnabled(true)
	            end
            	if self.cell then -- fixPos view
            		self:cellUpdate(dp[dir], dt)
            	else
            		self:sliding(dp[dir], dt)
            	end
	        end
        end
        local f = self.func[event.name]
        if f then f(event) end
    end)
    return cover
end

function Scrollview:coordinate(p)
	return self.dir == "x" and p or self:getlen()-p
end

function Scrollview:outrange(pos)
	local p = self:coordinate(pos)
	local head, tail = p, p+self.length
	local len = self.tail-self.head
	local r, bounce = 0, self.bounce or 0
	local f = self.dir == "x" and 1 or -1
	if head > self.head then
		if head < self.head + bounce then
			r = f * (self.head - head)
		else
			r = f * (-bounce)
			p = self.head + bounce
		end
	elseif self.length < len then
		if head > self.head - bounce then
			r = f * (self.head - head)
		else
			r = f * bounce
			p = self.head - bounce
		end
	elseif tail < self.tail then
		if tail > self.tail - bounce then
			r = f * (self.tail - tail)
		else
			r = f * bounce
			p = self.tail - bounce - self.length
		end
	end
	return self:coordinate(p), r
end

function Scrollview:drag(des)
	local des, r = self:outrange(des)
	if self.dir == "x" then
		self.content:setPositionX(des)
	else
		self.content:setPositionY(des)
	end
	return r
end

-- set content at fix position, (e.g. Pageview)
function Scrollview:cellUpdate(dp, dt)
	local TOUCH_SPEED = 1000
	local len = self.length / (self.cell - 1)
	local p = self:getLocation() - self.initpos
	local s = p - math.floor(p / len) * len
	if s == 0 then return end
	if math.abs(dp/dt) > TOUCH_SPEED then
		s = dp > 0 and len-s or -s
	else
		local r = self.turn_range
		local exc = (dp > 0) and r or 1-r
		s = (s < len * exc) and -s or len-s
	end
	self:move(s)
end

-- content can stop at arbitrary position
function Scrollview:sliding(dp, dt)
	if math.abs(dp/dt) < 1000 then return end
	local SLIDE_SPEED = 4000
	if self.slide ~= 0 then
		self.slide = self.slide + dp/dt
		self.content:stopAllActions()
	else
		self.slide = dp/dt
	end
	local dis = self.slide/SLIDE_SPEED * self:getlen()
	self:contentMoveBy(dis, function() 
		self.slide = 0
	end, math.abs(self.slide))
end

function Scrollview:move(dis, slow, func, speed)
	if dis == 0 then return end
	self:setTouchEnabled(false)
	local speed = speed or self.speed
	local time = math.abs(dis) / speed
	local s, moving = cc.p(0, 0), self.func.moved
	s[self.dir] = dis
	if not self.scheduler then 
		self.scheduler = xx.schedule(moving, 0.01)
	end
	local act = cc.MoveBy:create(time, s)
	act = cc.EaseOut:create(act, slow or 1)
    local seq = cc.Sequence:create(act,
        cc.CallFunc:create(function()
        	self.scheduler = xx.unschedule(self.scheduler)
            self:setTouchEnabled(true)
            local ended = self.func.ended
            if ended then ended() end
            if func then func() end
        end)
    )
	return self.content:runAction(seq)
end

function Scrollview:setTouchEnabled(flag)
	self.canTouch = flag
end

function Scrollview:regularTouchEvent(w)
	if w.func then
		local isTouch = w:isTouchEnabled()
		w:setTouchEventBy(self, w.func)
		w:setTouchEnabled(isTouch)
	end
	local children = w:getChildren()
	for k, widget in pairs(children) do
		self:regularTouchEvent(widget)
	end
end