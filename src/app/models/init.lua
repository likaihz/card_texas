mm = require("app.models.mm")

mm = table.include(mm, "widgets", {
	"unit", "window",
})
