package modules

import (
	student "app/modules/student/migrate"
	advisor "app/modules/advisor/migrate"
	headOfSubject "app/modules/headOfSubject/migrate"
	council "app/modules/council/migrate"
	facultyOffice "app/modules/facultyOffice/migrate"
	thesis "app/modules/thesis/migrate"
)

func MigrateModule() bool {
	student.MigrateTable();
	advisor.MigrateTable();
	headOfSubject.MigrateTable();
	council.MigrateTable();
	facultyOffice.MigrateTable();
	thesis.MigrateTable();
	return true
}
