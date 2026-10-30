package main

import "fmt"

func main() {
	// Название курса храним в константе
	const courseName = "Бэкенд-разработка на языке Go"

	var studentName string = "Алексей"
	var groupName string = "Учебная-01"

	completedJobs := 1
	interestedInBackend := true

	// Выводим значения с использованием fmt.Println и fmt.Printf
	fmt.Printf("Участник: %s\n", studentName)
	fmt.Println("Группа:", groupName)
	fmt.Println("Курс:", courseName)
	fmt.Printf("Завершено работ: %d\n", completedJobs)
	fmt.Printf("Интерес к бэкенду: %t\n", interestedInBackend)
}
