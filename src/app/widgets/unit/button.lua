-- mm.button --
module(..., package.seeall)

button = class("button", xxui.Btn)

function button:ctor(param)
    local args = clone(param)
    args.btn = xxres.button(args.btn)
    button.super.ctor(self, args)
    self:setTouchEvent(args.func)
end

-- interface --
function button:setTouchEvent(func)
    xxui.setTouchEvent(self, func)
    xxui.addTouchEvent(self, function()
        -- GM:playSound("btn/sound_9.mp3")
    end)
end
