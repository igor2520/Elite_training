package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"example.com/go-practice-01/internal/persistence"
	"example.com/go-practice-01/internal/tasks"
)

func main() {
	path := os.Getenv("TASKS_FILE")
	if path == "" {
		path = "tasks.json"
	}

	state, err := persistence.Load(path)
	if err != nil {
		fmt.Printf("Ошибка загрузки: %v\n", err)
		os.Exit(1)
	}

	store := tasks.NewMemoryStore()
	if err := store.Restore(state); err != nil {
		fmt.Printf("Ошибка инициализации: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Команды: add <название> <приоритет>, list, complete <id>, delete <id>, exit")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		args := strings.Fields(line)
		cmd := args[0]

		switch cmd {
		case "add":
			if len(args) < 3 {
				fmt.Println("Использование: add <название> <приоритет>")
				continue
			}
			priority, err := strconv.Atoi(args[len(args)-1])
			if err != nil {
				fmt.Println("Приоритет должен быть числом")
				continue
			}
			title := strings.Join(args[1:len(args)-1], " ")
			_, err = store.Add(title, priority)
			if err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			} else {
				fmt.Println("Задача добавлена")
			}

		case "list":
			list := store.List()
			if len(list) == 0 {
				fmt.Println("Список пуст")
			} else {
				for _, t := range list {
					fmt.Printf("[%d] %s | Статус: %s | Приоритет: %d\n", t.ID, t.Title, t.Status, t.Priority)
				}
			}

		case "complete":
			if len(args) < 2 {
				fmt.Println("Использование: complete <id>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("ID должен быть числом")
				continue
			}
			if err := store.Complete(id); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			} else {
				fmt.Println("Задача выполнена")
			}

		case "delete":
			if len(args) < 2 {
				fmt.Println("Использование: delete <id>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("ID должен быть числом")
				continue
			}
			if err := store.Delete(id); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			} else {
				fmt.Println("Задача удалена")
			}

		case "exit":
			return

		default:
			fmt.Println("Неизвестная команда")
			continue
		}

		// Сохраняем состояние после каждой успешной команды (кроме list и exit)
		if cmd != "list" && cmd != "exit" {
			if err := persistence.Save(path, store.GetState()); err != nil {
				fmt.Printf("ОШИБКА СОХРАНЕНИЯ: %v\n", err)
				os.Exit(1)
			}
		}
	}
}	