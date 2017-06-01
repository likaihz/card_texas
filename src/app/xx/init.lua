-- xx n xxui --
require "app.xx.fn.extn"

xx, xxui = {}, {}

table.include(xx, "xx", "fn")

table.include(xxui, "xx.ui", {
    "base", "unit", "view"
})

-- resource path --
local function respth()
    local tbl = {
        grid = 0, button = 0, icon = 0, panel = 0
    }
    for k in pairs(tbl) do
        tbl[k] = xxui.imgpath("ui/" .. k .. "/")
    end
    tbl.scene = xxui.imgpath("scene/")
    tbl.ani = function(pth)
        local p = "res/animation/".. pth
        return p .. ".plist", p .. ".png"
    end
    return tbl
end

xxres = respth()
