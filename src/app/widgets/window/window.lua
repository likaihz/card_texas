module(..., package.seeall)

-- 所有一、二级弹窗统一继承于 window
window = class("window", function()
	return xxui.Widget.new()
end)

-- stage: 所添加到的父节点
-- mode: 大弹窗 1，小弹窗 2
-- delete: boolean, 点击关闭时是否删除自身
function window:ctor(stage, mode, delete)
	self.mode_ = mode
	self.delete_ = delete
	self.func_ = {}
	self.btn_ = {}
	-- underlay --
	local m = delete and 2 or -2
	self.underlay = xxui.Underlay.new(stage, m, nil, 0.4)
	self:createUI(mode)
end

-- interface --
-- 获取子控件
-- name: "title", "board", "close", "return"
function window:getWidget(name)
	return self:getleaf("win_"..name)
end

-- 获取弹窗底部控件
-- name: btn's txt in English
function window:getBtn(name)
	if not name then
		return self.btn_
	end
	return self.btn_[name]
end

-- 设置弹窗底部按钮
-- name: btn's name or table of btns' names
function window:setBtn(name)
	for _, btn in ipairs(self.btn_) do
		btn:removeSelf()
	end
	self.btn_ = {}
	if not name then return {} end
	if type(name) == "string" then
		name = {name}
	end
	local num = table.len(name)
	local posy = {30, 20}
	local y = posy[self.mode_]
	for i, v in ipairs(name) do
		local x = self:align("x", num, i, 0, 0)
		local mode = self:getBtnMode(v)
		local txt = xx.translate(v, "button")
		local btn = xxui.Txtbtn.new {
			node = self, mode = mode, text = txt,
			anch = cc.p(0.5, 0), pos = cc.p(x, y)
		}
		self.btn_[v] = btn
	end
	return self.btn_
end

-- 设置弹窗标题文字
function window:setTitle(txt)
	txt = xx.translate(txt, "title")
	local w = self:getleaf("win_title")
	w:load(txt)
end

-- 控制弹窗打开和关闭
-- flag: boolean
function window:switch(flag)
	local key = flag and "opened" or "closed"
	local func = self.func_[key]
	if func then func() end
	self.underlay:switch(flag)
end

-- 删除弹窗控件
function window:remove()
	self.underlay:removeSelf()
end

-- 设置弹窗事件
-- key: "opened", "closed"
function window:setEvent(key, func)
	self.func_[key] = func
	if key == "closed" then
		self.underlay:setTouchEvent(func)
	end
end

-- 添加弹窗事件
function window:addEvent(key, func)
	if not func then return end
	local old = self.func_[key]
	local new = not old and func or function()
		old(); new()
	end
	self:setEvent(key, new)
end

-- 设置弹窗出现时所隐藏的节点（组）
-- 该节点（组）在弹窗关闭时恢复显示
-- node: widget or table of widgets
function window:setHidden(node)
	self.underlay:setHidden(node)
end

-- 添加弹窗出现时所隐藏的节点（组）
function window:addHidden(node)
	self.underlay:addHidden(node)
end

-- 控制全屏触摸事件
function window:disableScreenTouch()
	self.underlay:disable()
end

-- implementation --
function window:getBtnMode(txt)
	local color = {
		ok = "blue", register = "green", cancel = "orange",
		logout = "orange", exit = "orange",
	}
	return color[txt] or "yellow"
end

function window:close()
	if self.delete_ then
		self:remove()
	else
		self:switch(false)
	end
end


-- create widgets --
function window:createUI(mode)
	local img = "window"
	if mode == 2 then
		img = img.."_small"
	end
	local panel = xxui.create {
		node = self, img = xxres.panel(img), name = "win_panel"
	}
	self:setsize(panel)
	self:set {node = self.underlay, pos = "center"}
	self:setTouchEnabled(true)
	local btn = xxui.create {
		node = self, btn = xxres.button("close"), name = "win_close",
		anch = cc.p(1, 1), align = cc.p(1, 1), pos = cc.p(-30, -30),
		func = function() self:close() end
	}
	self:createTitle()
end

function window:createTitle()
	xxui.create {
		node = self, txt = "弹窗", name = "win_title",
		font = "txt", size = 40, anch = cc.p(0.5, 1),
		align = cc.p(0.5, 1), pos = cc.p(0, -35)
	}
	xxui.create {
		node = self, img = xxres.panel("window_light"),
		anch = cc.p(0.5, 1), align = cc.p(0.5, 1), pos = cc.p(0, -4)
	}
end
