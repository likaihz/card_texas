-- table lib --

local function checkTable(tbl)
    if type(tbl) ~= "table" then
        xx.error(tbl)
    end
    return true
end

-- 输出表tbl内的所有元素，默认为打印到console命令行控制台,支持class“对象”
-- name: tbl表名，需手动指定
-- path: 是输出目录，输出"name".lua文件到该目录
-- writeIndex: (boolean) 控制是否打印数字下标索引[n]，默认为false不打印
local function TABLE_OUTPUT(tbl, indent, write, writeIndex)
    local pre = string.rep("\t", indent)
    local k, v = next(tbl)
    while k do
        local str
        if type(k) == "number" then
            if writeNumIndex then str = pre.."["..k.."] = "
            else str = pre end       
        elseif type(k) == "string" then
            local _, l = string.find(k, "[_%w]*")
            if string.find(k, "^[_%a]") and l == string.len(k) then
                str = pre .. k .. " = "
            else str = pre .. "[\""..k.."\"] = " end
        else
            xx.error("the table index is neither NUMBER nor STRING!", 0)
        end
        if k == "class" then
            local cls = {}
            repeat
                table.insert(cls, v.__cname)
                v = v.super
            until not v
            v = table.concat(cls, "->")
        end

        if type(v) == "table" then
            write(str .. "{")
            TABLE_OUTPUT(v, indent + 1, write, writeNumIndex)
            if next(tbl, k) then write(pre .. "},")
            else write(pre .. "}") break end
        else
            if type(v) == "string" then
                v = string.gsub(v, "[\n]", "")
                v = string.gsub(v, "\"", "\\\"")
                v = "\"" .. v .. "\""
            end
            str = str .. tostring(v)
            if next(tbl, k) then write(str .. ",")
            else write(str) break end
        end
        k, v = next(tbl, k)
    end
end

function table.output(tbl, name, path, writeIndex)
    if not next(tbl) then
        error("The TABLE to output is blank!!! PLS check it!")
        return
    end
    local write = function(str) print(str) end
    if path then
        path = xx.fcheck(path)
        io.output(path .. name .. ".lua")
        write = function(str) io.write(str.."\n") end
        write("local "..name.." = {")
        TABLE_OUTPUT(tbl, 1, write, writeIndex)
        io.write("}\n\nreturn "..name)
        io.close()
        print(name, "输出完毕！此处应有掌声~")
    else
        name = name or ""
        print(name .. " TABLE print:")
        TABLE_OUTPUT(tbl, 1, write, writeIndex)
    end
end

-- kind: type of key want to count
-- kind == nil: return total number of elements 
function table.len(tbl, kind)
    local len, num = 0, {}
    for k in pairs(tbl) do
        len = len + 1
        local tp = type(k)
        if not num[tp] then
            num[tp] = 1
        else
            num[tp] = num[tp] + 1
        end
    end
    if kind then
        return num[kind] or 0
    end
    return len, num
end

-- check whether @tbl contains @item
-- return @true & its corresponding @key
function table.contain(tbl, val)
    checkTable(tbl)
    for k, v in pairs(tbl) do
        if v == val then
            return true, key
        end 
    end
end

function table.merge(tbl, t)
    for k, v in pairs(t) do
        tbl[k] = tbl[k] or v
    end
    return tbl
end

-- include packages to an engine table
-- pkgs: table of package names
function table.include(tbl, pth, pkgs)
    if type(pkgs) == "string" then
        pkgs = {pkgs}
    end
    tbl = tbl or {}
    pth = "app." .. pth .. "."
    for _, pkg in ipairs(pkgs) do
        local dir = pth .. pkg
        local files = require(dir .. ".#catalog")
        for _, name in ipairs(files) do
            local f = require(dir .. "." .. name)
            table.merge(tbl, f)
        end
    end
    return tbl
end

-- 随机返回多维表中的值，kind 控制识别表中的键类型
-- layer为数字（默认值为 1），控制返回第几层的值，tbl本身为第 0 层
-- layer = 2 则返回第2层，layer = -2 则返回倒数第二层
function table.random(tbl, kind, layer)
    if layer == 0 then return tbl end
    local vals, keys = {}, {}
    layer = layer or 1
    for l = 1, layer do
        local tmp = {}
        for k, v in pairs(tbl) do
            if not kind or kind == type(k) then
                table.insert(tmp, k)
            end
        end
        if #tmp == 0 then
            return
        end
        local key = tmp[math.random(#tmp)]
        tbl = tbl[key]
        table.insert(keys, key)
        table.insert(vals, tbl)
    end
    if layer == 1 then
        -- when layer is nil, return a key (not a key table)
        return vals[1], keys[1]
    end
    if layer < 0 then
        layer = layer + #vals + 1
    end
    -- protect the range of layer
    if layer < 1 or layer > #vals then
        xx.error("invalid argument: 'layer'!", 0)
    end
    return vals[layer], keys
end

-- compare two table t1 n t2
-- if same, return true or return false
function table.compare(tbl, t)
    if table.len(tbl) ~= table.len(t) then
        return false
    end
    for k, v1 in pairs(tbl) do
        local v2 = t[k]
        if type(v1) ~= type(v2) then
            return false
        end
        local same
        if type(v1) ~= "table" then
            same = (v1 == v2)
        else
            same = table.compare(v1, v2)
        end
        if not same then
            return false
        end
    end
    return true
end


-- string lib --

-- 解析以 char 分隔的字符串，为起始点为 start
-- 解析得到的字符段保存在返回的 strings 表里
function string.split(str, char, start)
    local pos = start or 0
    local head, tail = 0, 0
    local strings, s = {}
    repeat
        head = pos + 1
        pos = string.find(str, char, head)
        if not pos then break end
        tail = pos - 1
        s = string.sub(str, head, tail)
        table.insert(strings, s)
    until false
    s = string.sub(str, head)
    table.insert(strings, s)
    return strings
end

-- 按类别（数字、英文以及其他字符）分割字符串str, 插入tbl表返回
function string.regularize(str, tbl)
    tbl = tbl or {}
    local i, j = string.find(str, "%w+")
    if i == 1 then
        i, j = string.find(str, "%d+")
        if i ~= 1 then
            i, j = string.find(str, "%a+")  
        end
    elseif i then
        j = i - 1
    end
    local sub = string.sub(str, 1, j)
    table.insert(tbl, sub)
    str = string.sub(str, string.len(sub)+1)
    if string.len(str) ~= 0 then
        string.regularize(str, tbl)
    end
    return tbl
end