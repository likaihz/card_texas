module(..., package.seeall)

local crypto = require(cc.PACKAGE_NAME..".crypto")

-- rerquest: net request
function request(url, data, callback)
    local conn = network.createHTTPRequest(function(event)
    	local r = event.request
    	local err, message
	    if event.name ~= "completed" then
	    	err = r:getErrorCode()
	    	if err == 0 then
	    	-- retry to connect...
	    		return
	    	end
	    	message = r:getErrorMessage()
	        print("request error: ", err, message)
	        return callback {
	        	status = "fail", err = err, message = message
	        }
	    end
	    err = r:getResponseStatusCode()
	    if err ~= 200 then
	        print("no response: ", err)
	        return callback {
	        	status = "fail", err = err, message = "no response!"
	        }
	    end
	    -- request succeed
	    local res = r:getResponseData()
	    res = crypto.decryptXXTEA(res, "玄襄科技")
    	callback(json.decode(res))
    end, url, "POST")
    data = json.encode(data)
    data = crypto.encryptXXTEA(data, "玄襄科技")
    conn:setPOSTData(data)
    conn:start()
end


-- class websocket --
websocket = class("websocket")

function websocket:ctor(url)
	self.ws_ = cc.WebSocket:create(url)
end

function websocket:close()
	self.ws_:close()
end

function websocket:send(msg)
	msg = json.encode(msg)
	self.ws_:sendString(msg)
end

-- register callback
function websocket:onopen(func)
	self.ws_:registerScriptHandler(func, cc.WEBSOCKET_OPEN)
end

function websocket:onmessage(func)
	local function receive(msg)
		msg = json.decode(msg)
		if table.len(msg) == 0 then
			return
		end
		func(msg)
	end
	self.ws_:registerScriptHandler(receive, cc.WEBSOCKET_MESSAGE)
end

function websocket:onclose(func)
	self.ws_:registerScriptHandler(func, cc.WEBSOCKET_CLOSE)
end

function websocket:onerror(func)
	self.ws_:registerScriptHandler(func, cc.WEBSOCKET_ERROR)
end