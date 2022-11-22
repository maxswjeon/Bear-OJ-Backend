package routes

import (
	"github.com/gin-gonic/gin"

	admin_contests "github.com/maxswjeon/contest-backend/routes/admin/contests"
	admin_contests_problems "github.com/maxswjeon/contest-backend/routes/admin/contests/problems"
	admin_internalproblems "github.com/maxswjeon/contest-backend/routes/admin/internalproblems"
	admin_problems "github.com/maxswjeon/contest-backend/routes/admin/problems"
	admin_problems_judge "github.com/maxswjeon/contest-backend/routes/admin/problems/judge"
	admin_report_reset "github.com/maxswjeon/contest-backend/routes/admin/report/reset"
	admin_submits "github.com/maxswjeon/contest-backend/routes/admin/submits"

	// admin_submits_file "github.com/maxswjeon/contest-backend/routes/admin/submits/file"
	admin_session "github.com/maxswjeon/contest-backend/routes/admin/session"
	admin_users "github.com/maxswjeon/contest-backend/routes/admin/users"

	// admin_users_file "github.com/maxswjeon/contest-backend/routes/admin/users/file"
	contests "github.com/maxswjeon/contest-backend/routes/contests"
	problems "github.com/maxswjeon/contest-backend/routes/problems"
	problems_submits "github.com/maxswjeon/contest-backend/routes/problems/submits"
	report "github.com/maxswjeon/contest-backend/routes/report"
	report_init "github.com/maxswjeon/contest-backend/routes/report/init"
	scoreboard "github.com/maxswjeon/contest-backend/routes/scoreboard"
	session "github.com/maxswjeon/contest-backend/routes/session"

	// submits "github.com/maxswjeon/contest-backend/routes/submits"

	healthcheck "github.com/maxswjeon/contest-backend/routes/healthcheck"
)

func Apply(engine *gin.Engine) {
	engine.GET("/admin/contests", admin_contests.GET)
	engine.POST("/admin/contests", admin_contests.POST)
	engine.PATCH("/admin/contests/:id", admin_contests.PATCH)

	engine.PUT("/admin/contests/:id/problems/:pid", admin_contests_problems.PUT_)
	engine.DELETE("/admin/contests/:id/problems/:pid", admin_contests_problems.DELETE)

	engine.GET("/admin/internalproblems", admin_internalproblems.GET)

	engine.GET("/admin/problems", admin_problems.GET)
	engine.POST("/admin/problems", admin_problems.POST)
	engine.PATCH("/admin/problems/:id", admin_problems.PATCH_)
	engine.DELETE("/admin/problems/:id", admin_problems.DELETE_)

	engine.GET("/admin/problems/judge", admin_problems_judge.GET)

	engine.POST("/admin/report/reset", admin_report_reset.POST)

	engine.GET("/admin/session", admin_session.GET)
	engine.POST("/admin/session", admin_session.POST)

	engine.GET("/admin/submits", admin_submits.GET)

	// engine.GET("/admin/submits/file", admin_submits_file.GET)

	engine.GET("/admin/users", admin_users.GET)
	engine.POST("/admin/users", admin_users.POST)
	engine.PATCH("/admin/users/:id", admin_users.PATCH_)
	engine.DELETE("/admin/users/:id", admin_users.DELETE_)
	// engine.POST("/admin/users/file", admin_users_file.POST)

	engine.GET("/contests", contests.GET)
	engine.GET("/contests/:id", contests.GET_)

	engine.GET("/problems", problems.GET)
	engine.GET("/problems/:id", problems.GET_)
	engine.GET("/problems/:id/submits", problems_submits.GET)
	engine.POST("/problems/:id/submits", problems_submits.POST)

	engine.POST("/report", report.POST)
	engine.POST("/report/init", report_init.POST)

	engine.GET("/scoreboard", scoreboard.GET)

	engine.GET("/session", session.GET)
	engine.POST("/session", session.POST)
	engine.DELETE("/session", session.DELETE)

	// engine.GET("/submits", submits.GET)

	engine.GET("/healthcheck", healthcheck.GET)
}
