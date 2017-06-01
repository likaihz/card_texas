module(..., package.seeall)

-- func = {began =, moved =, ended =, canceled =}
local function check(func)
	func = func or {}
	if type(func) == "function" then
        func = {ended = func}
    end
    if type(func) ~= "table" then
    	xx.error(func)
    end
    return func
end

local register = {}

local function registerEvent(w, func)
    if type(func) ~= "table" then
        xx.error(func)
    end
    if w.type_ == "txt" then
        register.txt(w, func)
    else
        register.btn(w, func)
    end
end

-- interface --
function isCantouch(w)
	if w.canTouch == nil then
		return w:isTouchEnabled()
	else
		return w.canTouch
	end
end

function isInview(w, event)
	local rect = w:getBoundingBox()
    local pos = w:convertToWorldSpace(cc.p(0, 0))
    rect.x, rect.y = pos.x, pos.y
    return cc.rectContainsPoint(rect, event)
end

function isClick(w, dp, dt)
	local range, time = 60, 0.2
	if cc.pGetLength(dp) < range and dt < time then
		return true
	end
end

function setTouchEvent(w, func)
    w.func = check(func)
    registerEvent(w, w.func)
end

function addTouchEvent(w, func)
    if not func then return end
    func = xx.funcAppend(w.func, check(func))
    setTouchEvent(w, func)
end

-- set w's event which LIMITED by widget(e.g. scrollview)
function setTouchEventBy(w, widget, func)
	func = check(func)
	w.func = func
	local function f(evt)
		return widget:isCantouch() and widget:isInview(evt)
	end
	local start, time
	w:setTouchEvent {
		began = function(event)
			start = event
			time = os.clock()
			if func.began and f(event) then
				func.began(event)
			end
		end,
		moved = function(event)
			if func.moved and f(event) then
				func.moved(event)
			end
		end,
		ended = function(event)
			local dp = cc.pSub(event, start)
			local dt = os.clock() - time
			if func.ended and f(event) and widget:isClick(dp, dt) then
				func.ended(event)
			end
		end
	}
	w:setSwallowTouches(false)
end

function addTouchEventBy(w, widget, func)
	if not func then return end
	func = xx.funcAppend(w.func, check(func))
    setTouchEventBy(w, widget, func)
end

-- implementation of register --
-- register event on Txt
local TouchType = {
    began = "EVENT_TOUCH_BEGAN", moved = "EVENT_TOUCH_MOVED",
    ended = "EVENT_TOUCH_ENDED", canceled = "EVENT_TOUCH_CANCELLED"
}
function register.txt(w, func)
    local listener = cc.EventListenerTouchOneByOne:create()
    for e, t in pairs(TouchType) do
        local callback = function(touch, event)
            if not w.canTouch then return end
            event = touch:getLocation()
            if func[e] and isInview(w, event) then
                event.name = e
                event.target = w
                func[e](event)
            end
            return true
        end
        listener:registerScriptHandler(callback, cc.Handler[t])
    end
    local dispacther = cc.Director:getInstance():getEventDispatcher()
    dispacther:removeEventListenersForTarget(w)
    dispacther:addEventListenerWithSceneGraphPriority(listener, w)
end

-- register event on Btn n Img
local TouchPos = {
    began = "getTouchBeganPosition", moved = "getTouchMovePosition",
    ended = "getTouchEndPosition", canceled = "getTouchMovePosition"
}
function register.btn(w, func)
    w:setTouchEnabled(true)
    w:addTouchEventListener(function(sender, evt)
        for e, f in pairs(TouchPos) do
            if evt == ccui.TouchEventType[e] then
                if func[e] then
                    local event = ccui.Widget[f](sender)
                    event.name = e
                    event.target = sender
                    func[e](event)
                end
                return true
            end
        end
    end)
end