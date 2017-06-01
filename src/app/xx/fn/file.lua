-- xx.file --
module(..., package.seeall)

require("lfs")

------------------ 以下所有文件操作函数只接受以"/"分隔的路径----------------

-- 支持执行包内(src开始的）相对全路径文件，区别于只能执行从writablePath开始的全路径的dofile
-- 执行成功，返回文件内table数据
-- 执行失败，打印错误语句
function fexe(path)
    local str = io.readfile(path)
    if str == "" or not str then
        printInfo("The file \""..path.."\" is blank or not exists !!!")
        return {}
    end
    local fun = loadstring(str)
    local success, data = pcall(fun)
    if success then
        return data
    else
        xx.error(data, 0)
    end
end

-- 检查path路径是否存在, 不存在则递归创建此路径
-- 返回以"/"结尾的路径
local function createDir(path)
    if lfs.chdir(path) then
        printInfo("path check OK---> "..path)
        return true
    end
    local temp = string.sub(path, 1, -2)
    local pathinfo = io.pathinfo(temp)
    prePath = pathinfo.dirname
    printInfo("prePath:", prePath)
    if createDir(prePath) then
        lfs.mkdir(path)
        printInfo("path create OK---> "..path)
        return true
    end
end

function fcheck(path)
    local oldpath = lfs.currentdir()
    printInfo("old path---> "..oldpath)
    if string.find(path, "%.") then
        path = io.pathinfo(path).dirname
    end
    local path = string.gsub(device.writablePath, "[\\\\/]+$", "/")..path
    -- local path = string.gsub(device.writablePath, device.directorySeparator, "/") .. path
    -- 必须执行，因io.pathinfo 只能读到"/"分隔的目录
    if createDir(path) then
        if not string.find(path, "/$") then
            path = path.."/"
        end
        lfs.chdir(oldpath)  -- 还原工作路径
        printInfo("****** PATH IS READY ******")
        return path
    else
        xx.error("\""..path.."\" not exists and create it failed !!!", 0)
    end
end


-- 复制文件(目录以src/res开头)
-- string: source, destination, 分别对应源文件和目标文件(即：可重命名，不想重命名就写一样的文件名)
-- number: rowNum 控制是否在目标文件每行添加行数(位数为rowNum)，默认为false不打印(以备不时之需)
function fcopy(source, destination, rowNum)
    local writablePath = string.gsub(device.writablePath, "[\\\\/]+$", "/")
    source = writablePath..source
    if not io.exists(source) then
        xx.error("File \""..source.."\"".." is not EXISTS!", 0)
    end
    fcheck(destination)
    destination = writablePath..destination
    local src = io.open(source)
    local dst = io.open(destination, "w")
    local count = 0
    for line in src:lines() do
        count = count + 1
        if rowNum then
            local fmt = "%"..rowNum.."d  "
            dst:write(string.format(fmt, count), line, "\n")
        else
            dst:write(line, "\n")
        end
    end
    src:close()
    dst:close()
end


-- scan files(keyword: @key) in @path with operation of @func(@fileName, @path)
-- (recursive scan to subfolder if @intoFolder == true)
-- @skip both file and folder
-- return value "stop" of @func to control the mold of fscan()
-- mold 1: stop == nil, scan all the files
-- mold 2: stop == true, find only once n stop
-- mold 3: stop == "folder", scan all folders n find only once in every folder
function fscan(key, path, intoFolder, func, skip, param)
    if not string.find(path, "/$") then
        path = path.."/"
    end
    for fileName in lfs.dir(path) do
        if type(skip) ~= "table" then
            skip = {skip}
        end
        if fileName ~= "." and fileName ~= ".." and not table.contain(skip, fileName) then
            local f = path..fileName
            if string.find(fileName, key) then
                local stop = func(fileName, path, param)
                if stop then return stop end
            end
            local attr = lfs.attributes(f)
            if attr and attr.mode == "directory" and intoFolder then
                f = f.."/"
                local stop = fscan(key, f, intoFolder, func, skip)
                if stop and stop ~= "folder" then return stop end
            end
        end
    end
end