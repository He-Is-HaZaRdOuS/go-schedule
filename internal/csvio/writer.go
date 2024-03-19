package csvio

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/rhyrak/go-schedule/pkg/model"
)

// ExportSchedule formats the schedule data into ScheduleCSVRow structs and
// writes it to the CSV file specified by the given path.
func ExportSchedule(schedule *model.Schedule, path string) {
	nice := formatAndFilterSchedule(schedule)
	out, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	err = gocsv.MarshalFile(&nice, out)
	if err != nil {
		panic(err)
	}
}

// PrintSchedule prints weekly schedule grouped by department name.
func PrintSchedule(schedule *model.Schedule) {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	deps := make(map[string]bool, 10)
	nice := formatAndFilterSchedule(schedule)
	slices.SortFunc(nice, func(c1 *model.ScheduleCSVRow, c2 *model.ScheduleCSVRow) int {
		if dep := strings.Compare(c1.Department, c2.Department); dep != 0 {
			return dep
		}
		if grade := c1.Class - c2.Class; grade != 0 {
			return grade
		}
		if day := c1.Day - c2.Day; day != 0 {
			return day
		}
		if time := c1.Time - c2.Time; time != 0 {
			return time
		}
		return strings.Compare(c1.CourseCode, c2.CourseCode)
	})
	for _, c := range nice {
		if _, seen := deps[c.Department]; !seen {
			deps[c.Department] = true
			fmt.Printf("\n%s %s %s\n", strings.Repeat("-", (32-len(c.Department))/2), c.Department, strings.Repeat("-", int(0.5+(32-float32(len(c.Department)))/2.0)))
		}
		fmt.Printf("%-12s %-0.2d:%-0.2d   %-11s %d\n", days[c.Day], 8+(c.Time+30)/60, (c.Time+30)%60, c.CourseCode, c.Class)
	}
	fmt.Printf("Printed rows: %d\n", len(nice))
}

func formatAndFilterSchedule(schedule *model.Schedule) []*model.ScheduleCSVRow {
	var formatted []*model.ScheduleCSVRow
	var seen map[model.CourseID]bool = make(map[model.CourseID]bool)
	for dayIndex, day := range schedule.Days {
		for slotOffset, s := range day.Slots {
			for _, c := range s.CourseRefs {
				if _, ok := seen[c.CourseID]; !ok {
					seen[c.CourseID] = true
				} else {
					continue
				}
				classroom := c.Course_Environment
				if c.NeedsRoom {
					classroom = c.Classroom.ID
				}
				formatted = append(formatted, &model.ScheduleCSVRow{
					CourseCode: c.Course_Code,
					Day:        dayIndex,
					Duration:   c.Duration,
					Time:       slotOffset * schedule.TimeSlotDuration,
					Classrooms: classroom,
					Class:      c.Class,
					Department: c.DepartmentCode,
					CourseName: c.Course_Name,
					Lecturer:   c.Lecturer,
				})
			}
		}
	}
	return formatted
}
