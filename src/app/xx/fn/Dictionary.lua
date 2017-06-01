
local PATH = "app.data.language."
local Dictionary = class("Dictionary")

-- language = "CHS": 简体中文
function Dictionary:ctor(language)
    if language then
        self:load(language)
    end
end

function Dictionary:dtor()
    return self.name
end

function Dictionary:load(language)
    self.name = language
    self.lib = require(PATH..language)
end

local function findInTable(word, tbl)
    if string.find(word, "%w+") then
        return tbl[word]
    end
    for k, v in pairs(tbl) do
        if v == word then
            return k
        end
    end
end

function Dictionary:search(word, label)
    if not word then return "" end
    if type(word) == "string" then
        word = string.gsub(word, " ", "")
    end
    local lib = self.lib
    -- search in label "char" first
    local w
    if lib.char then
        w = lib.char[word]
        if w then return w end
    end
    if not label or not lib[label] then
        label = "default"
    end
    local w = findInTable(word, lib[label])
    if w then return w end
    for k, tbl in pairs(lib) do
        w = findInTable(word, tbl)
        if w then return w end
    end
    print(word.." is not in the dictionary!")
    return
end

-- input: a strings table, output a translated string
function Dictionary:translate(tbl, label)
    if not tbl then return "" end
    if type(tbl) ~= "table" then
        xx.error(tbl)
    end
    local str, word = {}
    for i = 1, table.maxn(tbl) do
        word = tbl[i] and self:search(tbl[i], label) or ""
        table.insert(str, word)
    end
    return table.concat(str)
end

-- test --
function Dictionary:test()
    local data = require(PATH.."test")
    local CN, EN
    for k, tbl in pairs(data) do
        print("in cls: ", k)
        for key, word in pairs(tbl) do
            CN = self:search(key) or "not found"
            EN = self:search(word) or "not found"
            print(key..": ", EN)
        end
    end
end
   
return Dictionary