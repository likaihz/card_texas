module(..., package.seeall)

-- error: 输出错误信息 
-- mode: nil 打印出错参数的类型, 0 打印错误详情 (p: string)
local err = error

function error(p, mode)
    if not mode then
        p = "input param is a \""..type(p).."\"!!"
    end
    print(_G.debug.traceback("", 2))
    print("------ xx.error ------")
    err(p, 2)
end

-- debug: 输出调试信息
function debug(str)
    print("------ xx.debug ------")
    print("INFO: "..tostring(str or "").. "\n")
    print(debug.traceback("", 3))
end

-- restrict: 将数字 num 限制在区间 [a, b]
function restrict(num, a, b)
    if type(num) ~= "number" then
        error(num)
    end
    if num < a then
        return a, true
    end
    if num > b then
        return b, true
    end
    return num
end

-- range: 在以 num 为中心，半径为 rad 的区间内随机取值
function range(num, r)
    r = math.abs(r)
    local p = 2 * r * math.random()
    return num - r + p
end

-- wave: 在 num*(1-r) 与 num*(1+r) 之间随机取值 
function wave(num, r)
    return range(num, num*r)
end

-- append: return an appended function
-- f1, f2: function or nil
function append(f1, f2)
    f2 = f2 or function() end
    if not f1 then return f2 end
    return function()
        f1(); f2()
    end
end

-- f1, f2: function table, append f2 to f1
function funcAppend(f1, f2)
    if not f1 then return f2 end
    if type(f2) == "function" then
        f2 = {ended = f2}
    end
    for e, fn in pairs(f2) do
        local pre = f1[e]
        f1[e] = function(event)
            if pre then pre(event) end
            fn(event)
        end
    end
    return f1
end

-- expression: an arithmetic expression expressed by table
-- operand: a function act on single object (cant be a table)
-- operation: a function act on two objects
function arithmetic(expression, operand, operation)
    if not expression or #expression == 0 then return end
    local function calcu(x)
        return type(x) == "table" and arithmetic(x, operand, operation) or operand(x)
    end
    local op1, op2 = calcu(expression[1])
    for i = 2, #expression, 2 do
        op2 = calcu(expression[i + 1])
        op1 = operation(op1, op2, expression[i])
    end
    return op1
end

-- dtor: destructor of obj (recursion)
function dtor(data)
    data = clone(data)
    local t = type(data)
    -- data is an obj
    if (t == "table" and data.class) or t == "userdata" then
        -- if an obj has no dtor() return nil
        return data.dtor and data:dtor()
    -- data is a normal table
    elseif t == "table" then
        for k, v in pairs(data) do
            data[k] = dtor(v)
        end
    end
    return data
end

-- inherit: 从多个函数表中继承方法
function inherit(class, ...)
    class = class or {}
    for i, methods in ipairs {...} do
        for key, fn in pairs(methods) do
            class[key] = fn
        end
    end
    return class
end

-- schedule n unschedule: 开启、关闭计时器
local scheduler = require(cc.PACKAGE_NAME .. ".scheduler")

function schedule(func, time)
    if not func then return end
    return scheduler.scheduleGlobal(func, time)
end

function unschedule(sch)
    if not sch then return end
    scheduler.unscheduleGlobal(sch)
end

-- 倒计时过程中每帧执行一次oncount，时间减为0后执行oncomplete并自动取消该计划
-- 每帧会将 实际剩余时间（非整数）time 传入oncount，可根据具体需求将time转化成想要的格式
-- 比如 math.ceil(time)  -- 剩余整数 秒
-- 比如tonumber(string.format("%.2f", time)) -- 保留两位小数
function countdown(time, oncount, oncomplete)
    if time < 0 then
        xx.error("argument 'time' must be positive", 0)
    end
    oncount = oncount or function() end
    local handle
    handle = scheduler.scheduleUpdateGlobal(function(dt)
        if time < 0 then
            oncount(0)
            scheduler.unscheduleGlobal(handle)
            if oncomplete then oncomplete() end
            return
        end
        oncount(time)
        time = time - dt
    end)
    return handle
end

function countsecond(time, oncount, oncomplete)
    local w = xxui.Widget.new()
    local i = math.floor(time)
    local co
    co = coroutine.create(function()
        while i >= 0 do
            if oncount then oncount(i) end
            i = i - 1
            w:delay(1, function()
                coroutine.resume(co)
            end)
            coroutine.yield()
        end
    end)
    coroutine.resume(co)
end

-- res = string: image
-- res = {path, name, frame, (time)}: animation
function cache(res, func)
    if type(res) == "string" then
        display.addImageAsync(res, func or function() end)
        return
    end
    local path, name = res.path, res.name
    local time = res.time or 1/17
    display.addSpriteFrames(path..".plist", path..".png", function()
        local frames = display.newFrames(name.."%04d.png", 1, res.frame)
        local ani = display.newAnimation(frames, time)
        display.setAnimationCache(name, ani)
        if func then func() end
    end)
end

-- tbl: a table hold number key n string value
-- return a function to get key (value) by value (key)
-- i: string (value) or number (key)
function map(tbl)
    return function(i)
        local tp = type(i) 
        if not i or tp == "number" then
            i = i or math.random(#tbl)
            return tbl[i]
        elseif type(i) == "string" then
            for j, k in pairs(tbl) do
                if i == k then
                    return j
                end
            end
        else
            xx.error(i)
        end
    end
end

-- translate: interface of dictionary
local Dictionary = require "app.xx.fn.Dictionary"

dictionary = Dictionary.new("CHS")
-- p1: word (string/number) or table of words
-- p2: key (class of the words)
function translate(p1, p2)
    if p1 == nil or p1 == "" then
        return ""
    end
    local t = type(p1)
    if t == "string" or t == "number" then
        return dictionary:search(p1, p2)
    elseif t == "table" then
        return dictionary:translate(p1, p2)
    else
        return xx.error(p1)
    end
end

-- read/write file
function io(pth, tbl)
    local state = require("framework.cc.utils.GameState")
    -- GameState.init(function, path, secretKey)
    state.init(function(param)
        if param.errorCode then
            print("io error:", param.errorCode)
            return
        end
        local name = param.name
        local val = param.values
        if name == "save" then
            local str = json.encode(val)
            -- str = crypto.encryptXXTEA(str, "XuanXiang")
            return {data = str} 
        elseif name == "load" then
            local str = val.data
            -- str = crypto.decryptXXTEA(str, "XuanXiang")
            return json.decode(str)
        end
    end, pth, "XX")
    if not tbl then
        return state.load()
    end
    state.save(tbl)
end


---- test ----
-- 文件网络下载
local function receive(conn)
    conn:settimeout(0)
    local s, status, partial = conn:receive(2^10)
    if status == "timeout" then
        coroutine.yield(conn)
    end
    return s or partial, status
end

-- host: "www.w3.org", file: "/TR/REC-html32.html"
local function download(host, file)
    require "socket"
    local conn = assert(socket.connect(host, 80))
    local count = 0
    conn:send("GET "..file.." HTTP/1.0\r\n\r\n")
    local out = {}
    while true do
        local s, status, partial = receive(conn)
        -- print(s or partial)
        table.insert(out, s or partial)
        count = count + #(s or partial)
        if status == "closed" then break end
    end
    conn:close()
    local name = string.reverse(file)
    name = string.sub(file, 1-string.find(name, "/"))
    table.output(out, name, "src/download test/")
    print(file, count)
end

-- download("www.w3.org", "/TR/REC-html32.html")
local threads = {}

local function get(host, file)
    local co = coroutine.create(function()
        download(host, file)
    end)
    table.insert(threads, co)
end

local function dispatch()
    local i = 1
    local connections = {}
    while true do
        if not threads[i] then
            if not threads[1] then break end
            i = 1
            connections = {}
        end
        local status, res = coroutine.resume(threads[i])
        if not res then
            table.remove(threads, i)
        else
            i = i + 1
            connections[#connections+1] = res
            if #connections == #threads then
                socket.select(connections)
            end
        end
    end
end

local function downloadTest()
    local host = "www.w3.org"
    local files = {
        "/TR/html401/html40.txt",
        "/TR/2002/REC-xhtml1-20020801/xhtml1.pdf",
        "/TR/REC-html32.html",
        "/TR/2000/REC-DOM-Level-2-Core-20001113/DOM2-Core.txt",
    }
    for i, file in ipairs(files) do
        get(host, file)
    end
    dispatch()
end

-- downloadTest()

-- weak table --
-- local a = {}
-- local meta = {__mode = "k"}

-- setmetatable(a, meta)

-- local key

-- for i = 1, 5 do
--     key = {}
--     a[key] = i
-- end

-- key = 1
-- a[key] = "hello"
-- collectgarbage()

-- for k, v in pairs(a) do print(v) end




