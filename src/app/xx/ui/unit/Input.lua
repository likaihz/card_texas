-- xxui.Button --
module(..., package.seeall)

Input = class("Input", xxui.Img)

--[[ args: {
  panel = 背板尺寸, size = 字体大小, ...  
} ]]
local INPUT_MODE = {
    phone = cc.EDITBOX_INPUT_MODE_PHONENUMBER,
    number = cc.EDITBOX_INPUT_MODE_NUMERIC,
}

-- node = , panel = cc.size(), name = , 
-- img = input box lay
-- mode = INPUT_MODE, color = cc.c3b(),
-- tips = {txt = string, color = cc.c3b()},
-- length = maxlength
function Input:ctor(param)
    local args = clone(param)
    local size = clone(args.panel)
    size.width = size.width + 50
    self:load(args.img or xxres.panel("under"))
    args.scale9 = size
    self:set(args)
    local img = xxres.panel("underlay")
    local w = cc.ui.UIInput.new({
        image = ccui.Scale9Sprite:create(img),
        size = args.panel, x = 0, y = 0,
        listener = args.func
    }):addTo(self)
    w:set {pos = "center"}
    if args.flag then
        w:setInputFlag(args.flag or 1)
    end
    w:setFontName(args.font or "ui/font/dft.ttf")
    w:setFontSize(args.size or 40)
    if args.color then
        w:setFontColor(args.color)
    end
    local tips = args.tips or ""
    if type(tips) == "table" then
        w:setPlaceholderFontName(tips.font or "ui/font/dft.ttf")
        if tips.color then
            w:setPlaceholderFontColor(tips.color)
        end
        tips = tips.txt or ""
    end
    w:setPlaceHolder(xx.translate(tips, "tips") or "")
    if args.mode then
        w:setInputMode(INPUT_MODE[args.mode])
    end
    if args.length then
        self.length = args.length
        w:setMaxLength(args.length)
    end
    w:setReturnType(cc.KEYBOARD_RETURNTYPE_DONE)
    self.core_ = w
end

function Input:getText()
    return self.core_:getText()
end

function Input:setText(str)
    return self.core_:setText(str)
end

function Input:setSwallowTouches(bool)
    self.core_:setSwallowTouches(bool)
end

function Input:checkLen()
    if not self.length then
        return true
    end
    local txt = self:getText()
    return txt:len() >= self.length
end

function Input:setEvent(func)
    self.core_:registerScriptEditBoxHandler(func)
end