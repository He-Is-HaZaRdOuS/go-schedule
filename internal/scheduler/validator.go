package scheduler

import (
	"fmt"

	"github.com/rhyrak/go-schedule/pkg/model"
)

// Validate checks schedule for conflicts and unassigned courses.
// Returns false and a message for invalid schedules.
func Validate(courses []*model.Course, labs []*model.Laboratory, schedule *model.Schedule, rooms []*model.Classroom, congestedDepartments map[string]int, CongestionLimit int) ([]*model.Course, bool, bool, string, int) {
	var message string
	var valid bool = true
	var allAssigned bool
	var hasCourseCollision bool
	var hasClassroomCollision bool

	// Find and store unassigned courses
	unassignedCount := 0
	var unassignedCourses []*model.Course
	var unassignedLabs []*model.Laboratory
	for _, c := range courses {
		if !c.Placed {
			unassignedCount++
			unassignedCourses = append(unassignedCourses, c)
		}
	}

	for _, l := range labs {
		if !l.Placed {
			unassignedCount++
			unassignedLabs = append(unassignedLabs, l)
		}
	}

	// Display error message if any course remains unassigned
	if unassignedCount > 0 {
		message = fmt.Sprintf("- There are %d unassigned courses:\n", unassignedCount)
		for _, un := range unassignedCourses {
			var ce string = "Compulsory"
			if !un.Compulsory {
				ce = "Elective"
			}
			message += fmt.Sprintf("IN-CLASS    %s %s %s %d %s\n", ce, un.Course_Code, un.Department, un.Number_of_Students, un.Lecturer)
		}
		for _, un := range unassignedLabs {
			var ce string = "Compulsory"
			if !un.Compulsory {
				ce = "Elective"
			}
			message += fmt.Sprintf("LABORATORY    %s %s %s %d %s\n", ce, un.Course_Code, un.Department, un.Number_of_Students, un.Lecturer)
		}
		message += "\n"
	}

	allAssigned = unassignedCount == 0

	// Check for course conflict
	ok, msg := checkCourseCollision(schedule)
	hasCourseCollision = !ok
	message += msg

	// Check for classroom conflict
	ok, msg = checkClassroomCollision(schedule)
	hasClassroomCollision = !ok
	message += msg

	var sufficientRooms bool = true
	message = "\n" + message

	// Display messages accordingly
	if hasClassroomCollision {
		message = "[FAIL]: Classroom collision check.\n" + message
		valid = false
	} else {
		message = "[  OK]: Classroom collision check.\n" + message
	}
	if hasCourseCollision {
		message = "[FAIL]: Course collision check.\n" + message
		valid = false
	} else {
		message = "[  OK]: Course collision check.\n" + message
	}
	if !allAssigned {
		message = "[FAIL]: Course has classroom check.\n" + message
		valid = false
		sufficientRooms = false
	} else {
		message = "[  OK]: Course has classroom check.\n" + message
	}

	return unassignedCourses, valid, sufficientRooms, message, unassignedCount
}

func checkCourseCollision(schedule *model.Schedule) (bool, string) {
	valid := true
	message := ""
	for _, day := range schedule.Days {
		for _, slot := range day.Slots {
			for _, c1 := range slot.CourseRefs {
				for _, c2 := range slot.Courses {
					if contains(c1.ConflictingCourses, c2) && !c1.ServiceCourse {
						valid = false
						message += "Conflicting courses placed at the same time\n"
					}
				}
			}
		}
	}
	return valid, message
}

func checkClassroomCollision(schedule *model.Schedule) (bool, string) {
	valid := true
	message := ""
	for _, day := range schedule.Days {
		for _, slot := range day.Slots {
			var usedRooms map[string]bool = make(map[string]bool)
			for _, c := range slot.CourseRefs {
				if c.Classroom == nil {
					continue
				}
				_, usedBefore := usedRooms[c.Classroom.ID]
				if usedBefore {
					schedule.Cost++
					valid = false
					message += "- Classroom " + c.Classroom.ID + " assigned multiple times\n"
				} else {
					usedRooms[c.Classroom.ID] = true
				}
			}
		}
	}
	return valid, message
}

func contains(s []model.CourseID, e model.CourseID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
