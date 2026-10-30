package main

import "fmt"
import "os"

func main() {
	const minutesInHour = 60

	testCases := []int{0, 45, 60, 150, 360, 1440}

	file, err := os.Create("Readme.md")
	if err != nil {
		fmt.Println("Не удалось создать файл:", err)
		return
	}
	defer file.Close()

	fmt.Fprintln(file, "# Фактические результаты проверок работы программы")
	fmt.Fprintln(file, "| Входные минуты | Вычисленные полные часы | Оставшиеся минуты |")
	fmt.Fprintln(file, "| :--- | :--- | :--- |")
	
	for _, totalMinutes := range testCases {
		if totalMinutes < 0 {
			fmt.Fprintf(file, "| %d минут | Ошибка: не может быть < 0 | Ошибка |\n", totalMinutes)
		continue
		}
		hours := totalMinutes / minutesInHour
		remainder := totalMinutes % minutesInHour

		fmt.Fprintf(file, "| %d минут | %d ч. | %d мин. |\n", totalMinutes, hours, remainder)

	}

	fmt.Println("Результаты успешно вычислены и записаны в README.md!")
}
