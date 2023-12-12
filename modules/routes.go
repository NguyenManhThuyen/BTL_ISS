package modules

import (
	authenRoute "app/modules/authen/routes"
	studentRoute "app/modules/student/routes"
	advisorRoute "app/modules/advisor/routes"
	headOfSubjectRoute "app/modules/headOfSubject/routes"
	councilRoute "app/modules/council/routes"
	facultyOfficeRoute "app/modules/facultyOffice/routes"
	thesisRoute "app/modules/thesis/routes"
	usersRoute "app/modules/users/routes"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	authenRoute.InitAuthenRoutes(app)
	studentRoute.InitStudentRoutes(app)
	advisorRoute.InitAdvisorRoutes(app)
	headOfSubjectRoute.InitHeadOfSubjectRoutes(app)
	councilRoute.InitCouncilRoutes(app)
	facultyOfficeRoute.InitFacultyOfficeRoutes(app)
	thesisRoute.InitThesisRoutes(app)
	usersRoute.InitUsersRoutes(app)
}
