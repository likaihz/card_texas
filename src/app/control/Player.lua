-- class Player --
local Player = class("Player")

function Player:ctor(data)
	self.data = data
end

-- interface --
-- key: "name", "gold", "jade"...
function Player:get(key)
	local data = self.data
	if data[key] then
		return data[key]
	end
end

-- key: "roomcard", "gold"...
function Player:set(key, val)
	local data = self.data
	data[key] = val
	local w = GM:getscene().topbar
	if w then
		w:load(key)
	end
end

-- key: "gold", "jade"...
function Player:add(key, val)
	if not val then return end
	local data = self.data
	local v = data[key]
	self:set(key, v+val)
end

return Player