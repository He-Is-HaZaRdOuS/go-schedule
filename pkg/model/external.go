package model

type External struct {
	Section                  int        `csv:"Section"`
	Course_Code              string     `csv:"Course_Code"`
	Course_Name              string     `csv:"Course_Name"`
	Number_of_Students       int        `csv:"Number_of_Students"`
	Course_Environment       string     `csv:"Course_Environment"`
	TplusU                   string     `csv:"T+U"`
	AKTS                     int        `csv:"AKTS"`
	Class                    int        `csv:"Class"`
	Department               string     `csv:"Department"`
	Lecturer                 string     `csv:"Lecturer"`
	StartingTimeSTR          string     `csv:"Starting_Time"`
	DaySTR                   string     `csv:"Day"`
	CourseRef                *Course    `csv:"_"`
	Duration                 int        `csv:"-"`
	CourseID                 CourseID   `csv:"-"`
	ConflictingCourses       []CourseID `csv:"-"`
	Placed                   bool       `csv:"-"`
	Classroom                *Classroom `csv:"-"`
	NeedsRoom                bool       `csv:"-"`
	NeededSlots              int        `csv:"-"`
	Reserved                 bool       `csv:"-"`
	ReservedStartingTimeSlot int        `csv:"-"`
	ReservedDay              int        `csv:"-"`
	BusyDays                 []int      `csv:"-"`
	Compulsory               bool       `csv:"-"`
	ConflictProbability      float64    `csv:"_"`
	DisplayName              string     `csv:"_"`
	ServiceCourse            bool       `csv:"_"`
	HasBeenSplit             bool       `csv:"_"`
	IsFirstHalf              bool       `csv:"_"`
	HasLab                   bool       `csv:"_"`
	PlacedDay                int        `csv:"_"`
	AreEqual                 bool       `csv:"_"`
	IsBiggerHalf             bool       `csv:"_"`
	OtherHalfID              CourseID   `csv:"_"`
}
