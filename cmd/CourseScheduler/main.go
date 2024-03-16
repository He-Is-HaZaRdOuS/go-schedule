package main

import (
	"fmt"
	"math/rand"

	"github.com/rhyrak/CourseScheduler/internal/csvio"
	"github.com/rhyrak/CourseScheduler/internal/scheduler"
	"github.com/rhyrak/CourseScheduler/pkg/model"
)

// Program parameters
const (
	ClassroomsFile   = "./res/private/classrooms.csv"
	CoursesFile      = "./res/private/courses2.csv"
	ExportFile       = "schedule.csv"
	NumberOfDays     = 5
	TimeSlotDuration = 60
	TimeSlotCount    = 9
)

func main() {
	classrooms := csvio.LoadClassrooms(ClassroomsFile, ';')
	ignoredCourses := []string{"ENGR450", "IE101", "CENG404"}
	courses := csvio.LoadCourses(CoursesFile, ';', ignoredCourses)

	var schedule *model.Schedule
	for limit := 0; limit < 100; limit++ {
		for _, c := range classrooms {
			c.CreateSchedule(NumberOfDays, TimeSlotCount)
		}
		for _, c := range courses {
			c.Placed = false
		}
		rand.Shuffle(len(courses), func(i, j int) {
			courses[i], courses[j] = courses[j], courses[i]
		})
		schedule = model.NewSchedule(NumberOfDays, TimeSlotDuration, TimeSlotCount)
		scheduler.FillCourses(courses, schedule, classrooms)
		if valid, _ := scheduler.Validate(courses, schedule, classrooms); valid {
			break
		}
	}

	csvio.ExportSchedule(schedule, ExportFile)
	valid, msg := scheduler.Validate(courses, schedule, classrooms)
	if !valid {
		fmt.Println("Invalid schedule:")
	} else {
		fmt.Println("Passed all tests")
	}
	fmt.Println(msg)
	schedule.CalculateCost()
	fmt.Printf("Cost: %d\n", schedule.Cost)
}
