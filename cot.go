package main

import (
	"fmt"
)

const (
	SATURDAY = iota
	SUNDAY
	MONDAY
	TUESDAY
	WEDNESDAY
	NUM_DAYS
)

// ddd
const (
	SLOT_8_10 = iota
	SLOT_10_12
	SLOT_14_16
	SLOT_16_18
	NUM_SLOTS
)

type Schedule [NUM_DAYS][NUM_SLOTS]string

type CourseSchedule struct {
	Day      int
	TimeSlot int
}

type Course struct {
	CourseID    int
	CourseName  string
	CourseMajor string
	TeacherID   int
	Credits     int
	Capacity    int
	Schedules   []CourseSchedule
	EnrolledIDs []int
}

type Student struct {
	StudentID    int
	StudentName  string
	StudentMajor string
	Schedule     Schedule
	Courses      []int
	TotalCredits int
}

type Teacher struct {
	TeacherID    int
	TeacherName  string
	TeacherMajor string
	Schedule     Schedule
	Courses      []int
}

func NewSchedule() Schedule {
	var s Schedule
	for d := 0; d < NUM_DAYS; d++ {
		for slot := 0; slot < NUM_SLOTS; slot++ {
			s[d][slot] = ""
		}
	}
	return s
}

func (s Schedule) Display() {
	days := []string{"شنبه", "یک‌شنبه", "دوشنبه", "سه‌شنبه", "چهارشنبه"}
	timeSlots := []string{"8-10", "10-12", "14-16", "16-18"}

	for d := 0; d < NUM_DAYS; d++ {
		fmt.Printf("روز %s:\n", days[d])
		for slot := 0; slot < NUM_SLOTS; slot++ {
			status := s[d][slot]
			if status == "" {
				status = "آزاد"
			}
			fmt.Printf("  بازه %s: %s\n", timeSlots[slot], status)
		}
		fmt.Println()
	}
}

func (s Schedule) IsAvailable(day int, slot int) bool {
	if day < 0 || day >= NUM_DAYS {
		fmt.Println("روز نامعتبر است.")
		return false
	}
	if slot < 0 || slot >= NUM_SLOTS {
		fmt.Println("بازه زمانی نامعتبر است.")
		return false
	}
	return s[day][slot] == ""
}

func (s *Schedule) Reserve(day int, slot int, courseName string) bool {
	if !s.IsAvailable(day, slot) {
		return false
	}
	s[day][slot] = courseName
	return true
}

func (s *Schedule) Unreserve(day int, slot int) bool {
	if s.IsAvailable(day, slot) {
		return false
	}
	s[day][slot] = ""
	return true
}

func AddCourse(courses *[]Course, courseID int, courseName, courseMajor string, teacherID int, credits, capacity int, schedules []CourseSchedule) {
	newCourse := Course{
		CourseID:    courseID,
		CourseName:  courseName,
		CourseMajor: courseMajor,
		TeacherID:   teacherID,
		Credits:     credits,
		Capacity:    capacity,
		Schedules:   schedules,
		EnrolledIDs: []int{},
	}
	*courses = append(*courses, newCourse)
	fmt.Printf("درس '%s' با شناسه %d و ظرفیت %d اضافه شد.\n", courseName, courseID, capacity)
}

func ReserveTeacherTime(teacher *Teacher, course Course) bool {

	for _, sched := range course.Schedules {
		if !teacher.Schedule.IsAvailable(sched.Day, sched.TimeSlot) {
			fmt.Printf("معلم %s در روز %s و بازه %s مشغول است و نمی‌تواند درس '%s' را تدریس کند.\n",
				teacher.TeacherName, getDayName(sched.Day), getSlotName(sched.TimeSlot), course.CourseName)
			return false
		}
	}

	for _, sched := range course.Schedules {
		success := teacher.Schedule.Reserve(sched.Day, sched.TimeSlot, course.CourseName)
		if !success {

			fmt.Printf("خطا در رزرو زمان‌بندی %s روز %s برای درس '%s'.\n",
				getSlotName(sched.TimeSlot), getDayName(sched.Day), course.CourseName)
			return false
		}
	}

	teacher.Courses = append(teacher.Courses, course.CourseID)
	fmt.Printf("درس '%s' با شناسه %d به معلم %s اختصاص یافت.\n", course.CourseName, course.CourseID, teacher.TeacherName)
	return true
}

func EnrollStudentInCourse(student *Student, course *Course) bool {

	if student.StudentMajor != course.CourseMajor {
		fmt.Printf("رشته تحصیلی دانشجو %s با رشته تحصیلی درس '%s' مطابقت ندارد.\n",
			student.StudentName, course.CourseName)
		return false
	}

	if len(course.EnrolledIDs) >= course.Capacity {
		fmt.Printf("ظرفیت درس '%s' پر است و دانشجو %s نمی‌تواند ثبت‌نام کند.\n",
			course.CourseName, student.StudentName)
		return false
	}

	if student.TotalCredits+course.Credits > 20 {
		fmt.Printf("دانشجو %s نمی‌تواند بیش از 20 واحد ثبت‌نام کند.\n", student.StudentName)
		return false
	}

	for _, sched := range course.Schedules {
		if !student.Schedule.IsAvailable(sched.Day, sched.TimeSlot) {
			fmt.Printf("دانشجو %s در روز %s و بازه %s مشغول است و نمی‌تواند در درس '%s' ثبت‌نام کند.\n",
				student.StudentName, getDayName(sched.Day), getSlotName(sched.TimeSlot), course.CourseName)
			return false
		}
	}

	for _, sched := range course.Schedules {
		student.Schedule.Reserve(sched.Day, sched.TimeSlot, course.CourseName)
	}

	course.EnrolledIDs = append(course.EnrolledIDs, student.StudentID)

	student.Courses = append(student.Courses, course.CourseID)

	student.TotalCredits += course.Credits

	fmt.Printf("دانشجو %s با موفقیت در درس '%s' ثبت‌نام کرد.\n", student.StudentName, course.CourseName)
	return true
}

func getDayName(day int) string {
	days := []string{"شنبه", "یک‌شنبه", "دوشنبه", "سه‌شنبه", "چهارشنبه"}
	return days[day]
}

func getSlotName(slot int) string {
	slots := []string{"8-10", "10-12", "14-16", "16-18"}
	return slots[slot]
}

func main() {
	courses := []Course{}
	students := []Student{}
	teachers := []Teacher{}

	AddCourse(&courses, 1, "ریاضی 1", "ریاضی", 101, 3, 10, []CourseSchedule{{SATURDAY, SLOT_8_10}})
	AddCourse(&courses, 2, "فیزیک 1", "فیزیک", 102, 4, 10, []CourseSchedule{{SUNDAY, SLOT_10_12}})
	AddCourse(&courses, 3, "برنامه‌نویسی", "کامپیوتر", 103, 4, 10, []CourseSchedule{{MONDAY, SLOT_14_16}})

	student1 := Student{StudentID: 1, StudentName: "علی", StudentMajor: "کامپیوتر", Schedule: NewSchedule()}
	student2 := Student{StudentID: 2, StudentName: "مریم", StudentMajor: "ریاضی", Schedule: NewSchedule()}
	students = append(students, student1, student2)

	teacher1 := Teacher{TeacherID: 101, TeacherName: "آقای رضایی", TeacherMajor: "ریاضی", Schedule: NewSchedule()}
	teacher2 := Teacher{TeacherID: 102, TeacherName: "خانم محمدی", TeacherMajor: "کامپیوتر", Schedule: NewSchedule()}
	teachers = append(teachers, teacher1, teacher2)

	ReserveTeacherTime(&teacher1, courses[0])
	ReserveTeacherTime(&teacher2, courses[2])

	EnrollStudentInCourse(&students[0], &courses[2])
	EnrollStudentInCourse(&students[1], &courses[0])

	fmt.Println("\nزمان‌بندی دانشجویان:")
	for _, student := range students {
		fmt.Printf("دانشجو: %s\n", student.StudentName)
		student.Schedule.Display()
	}
	fmt.Println("\nزمان‌بندی معلمان:")
	for _, teacher := range teachers {
		fmt.Printf("معلم: %s\n", teacher.TeacherName)
		teacher.Schedule.Display()
	}
}
