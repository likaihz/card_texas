-- models manager: mm --
module(..., package.seeall)

-- mm.obj(): obj 类的创建控制
local OBJ = require "app.models.obj"
poker = require "app.models.poker"

-- id: "tank_A02", cls is "tank"
local function getcls(id)
    local t = string.split(id, "_")
    return t[1]
end

-- data (string): 判断是否为可创建类名
-- data (object): 直接返回 data
-- data (tbl): 将 data 做成对象返回
function obj(data)
    data = clone(data)
    if type(data) == "string" then
        return OBJ[data] and true or false
    end
    if type(data) ~= "table" then
        xx.error(data)
    end
    if not data.cls then
        data.cls = getcls(data.id)
    end
    local cls = data.cls
    if not OBJ[cls] then
        xx.error("invalid input!", 0)
    end
    return OBJ[cls].new(data)
end

-- mm.rank is a function
-- get string (number) by number (string)
rank = xx.map {
    "white", "green", "blue", "purple", "orange", "red"
}

-- pop obj's id
function pop(cls)
    local cfg = GM:getconfig(cls)
    local _, id = table.random(cfg)
    return id
end

-- id: obj's id or cls
-- data: user's data (table) or lv (number)
function create(id, data)
    data = clone(data)
    if obj(id) then  -- id is a cls
        id = pop(id)
    end
    -- config data --
    local cls = getcls(id)
    local cfg = GM:getconfig(cls)
    if not cfg[id] then
        xx.error("no such obj!!", 0)
    end
    data = table.merge(data, cfg[id])
    -- association data --
    -- data.association = getassociation(id)
    return obj(data)
end

function getassociation(id)
    local cfg = GM:getconfig("association")
    local tbl = {}
    for _, asso in pairs(cfg) do
        for _, v in ipairs(asso.comb) do
            if v == id then
                table.insert(tbl, asso)
            end
        end 
    end
    return tbl
end

-- sort obj by rank n lv
function arrange(data)
    table.sort(data, function(a, b)
        local ra = a:get("rank")
        local rb = b:get("rank")
        if ra ~= rb then
            return ra > rb
        end
        return a:get("lv") > b:get("lv")
    end)
    return data
end
