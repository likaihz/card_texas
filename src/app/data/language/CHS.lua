local tbl = {}

tbl.char = {
    [","] = "，", [":"] = "：", [";"] = "；", ["!"] = "！",
    ["["] = "『", ["]"] = "』", ["."] = "。", ["?"] = "？",
}

tbl.default = {
    roundnum = "局数", turning  = "坐庄", roomnum  = "房号",
    rotated = "轮庄", fixed = "固定庄", challenged = "抢庄",
    raise = "倍", player = "玩家", cards = "手牌", rank = "牌型",
    addscore = "积分",
}

tbl.niuniu = {
    [0] = "没牛", "牛一", "牛二", "牛三", "牛四",
    "牛五", "牛六", "牛七", "牛八", "牛九", "牛牛",
    [15] = "顺子", [16] = "同花", [17] = "葫芦",
    [18] = "小牛", [19] = "花牛", [20] = "炸弹",
    [25] = "同花顺",
}

tbl.tips = {
    raise = "请选择下分倍数：", roomnum = "请输入房间号",
    room = "亲，该房间", fail = "创建失败", unexist = "不存在",
    full = "已满",
}

tbl.number = {
    [0] = "零", "一", "二", "三", "四", "五",
    "六", "七", "八", "九", "十", 
}

tbl.button = {
    ready = "准备", ok = "确定", logout = "登出", exit = "退出",
    cancel = "取消",
}

tbl.title = {
    sharing = "分享", setting = "设置", friends = "朋友圈",
    wechat = "微信群", sound = "音效", music = "音乐",
    roomnum = "房号",
}

tbl.verb = {
    clean = "重输", delete = "删除",
}

return tbl