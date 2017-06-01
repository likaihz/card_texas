-- mm.result --
module(..., package.seeall)

result = class("result", function()
    return xxui.Widget.new()
end)

function result:ctor(scene)
    self.scene = scene
    self.underlay = xxui.Underlay.new(scene, 0, nil, 0.7)
    self.content = {}
    self:init()
end

-- interface --
function result:load(battle)
    local pos = cc.p(display.width/2, 0)
    local i = 1
    battle:polling(0.03, function(ft)
        local content = self:loadcontent(i, ft)
        content:moveby(0.5, pos)
        i = i + 1
    end)
end

function result:loadcontent(i, ft)
    local content = self.content[i]
    content:setVisible(true)
    local tbl = {name = 0, rank = 0, addscore = 0}
    for k in pairs(tbl) do
        local w = content:getchild(k)
        local v = ft:get(k)
        if k == "rank" then
            v = xx.translate(v, "niuniu")
        elseif k == "addscore" then
            w:set {color = ft:scorecolor(v)}
            if v > 0 then v = "+"..v end
        end
        w:load(v)
    end
    -- load pokers
    return content
end

function result:onclose(func)
    self.btn:setTouchEvent(function()
        if func then func() end
        self:remove()
    end)
end

function result:remove()
    self.underlay:removeSelf()
end

-- create widgets --
function result:init()
    local panel = xxui.create {
        node = self, img = xxres.panel("back"),
    }
    self:setsize(panel)
    self:set {
        node = self.underlay, pos = "center"
    }
    self:createbtn()
    self:creategrid()
    self:createtitle()
end

function result:createbtn()
    local txt = xx.translate("ok", "button")
    self.btn = xxui.Txtbtn.new {
        node = self, mode = "blue", text = txt,
        anch = cc.p(0.5, 0.5), align = cc.p(0.5, 0.08),
        func = function() self:remove() end
    }
end

function result:createtitle()
    local panel = xxui.create {
        node = self, img = xxres.panel("title_red"), name = "title_panel",
        anch = cc.p(0.5, 0.5), align = cc.p(0.5, 0.85)
    }
    xxui.create {
        node = panel, img = xxres.icon("win"), name = "title_icon",
        anch = cc.p(0.5, 0.5), align = cc.p(0.5, 0.75)
    }
end

function result:creategrid()
    for i = 1, 5 do
        self:newgrid(i)
    end
    local grid = self:getchild("grid_1")
    local tbl = {"player", "cards", "rank", "addscore"}
    for i, txt in pairs(tbl) do
        self:newtext(grid, i, txt)
    end
end

function result:newgrid(i)
    local tbl = {[0] = "green", "blue"}
    local img = "bar_"..tbl[i%2]
    local gap = 115
    local y = (1-i) * gap
    local grid = xxui.create {
        node = self, img = xxres.panel(img), name = "grid_"..i,
        anch = cc.p(0.5, 0.5), align = cc.p(0.5, 0.65), pos = cc.p(0, y)
    }
    local content = xxui.Widget.new()
    content:setsize(grid)
    content:set {
        node = grid, name = "content",
        anch = cc.p(0.5, 0.5), align = cc.p(0.5, 0.5),
        pos = cc.p(-display.width/2, 0)
    }
    self:newcontent(content)
    content:setVisible(false)
    self.content[i] = content
end

function result:newtext(node, i, txt)
    txt = xx.translate(txt)
    xxui.create {
        node = node, txt = txt, font = "txt",
        size = 50, color = cc.c3b(190, 190, 224),
        anch = cc.p(0.5, 0), align = cc.p(i/5, 1)
    }
end

function result:newcontent(node)
    xxui.create {
        node = node, txt = "老司机", name = "name",
        anch = cc.p(0.5, 0.5), align = cc.p(1/5, 0.5)
    }
    xxui.create {
        node = node, img = "obj/poker/poker_back.png", name = "poker",
        anch = cc.p(0.5, 0.5), align = cc.p(2/5, 0.5)
    }
    xxui.create {
        node = node, txt = "牛牛", name = "rank",
        anch = cc.p(0.5, 0.5), align = cc.p(3/5, 0.5)
    }
    xxui.create {
        node = node, txt = "+50", name = "addscore",
        anch = cc.p(0.5, 0.5), align = cc.p(4/5, 0.5)
    }
end