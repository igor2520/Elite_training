package tasks

import (
	"errors"
	"fmt"
)

type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority int    `json:"priority"`
}

type State struct {
	NextID int    `json:"next_id"`
	Tasks  []Task `json:"tasks"`
}

type MemoryStore struct {
	nextID int
	tasks  []Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		nextID: 1,
		tasks:  []Task{},
	}
}

// Restore восстанавливает состояние из JSON, проверяя данные на корректность
func (s *MemoryStore) Restore(state State) error {
	if state.NextID < 1 {
		return errors.New("next_id не может быть меньше 1")
	}

	seenIDs := make(map[int]bool)
	maxID := 0

	for _, t := range state.Tasks {
		if t.ID <= 0 {
			return fmt.Errorf("некорректный ID задачи: %d", t.ID)
		}
		if seenIDs[t.ID] {
			return fmt.Errorf("дубликат ID задачи: %d", t.ID)
		}
		if t.Title == "" {
			return fmt.Errorf("пустое название задачи (ID: %d)", t.ID)
		}
		if t.Priority < 1 || t.Priority > 5 {
			return fmt.Errorf("некорректный приоритет (ID: %d): %d", t.ID, t.Priority)
		}
		if t.Status != "new" && t.Status != "done" {
			return fmt.Errorf("некорректный статус (ID: %d): %s", t.ID, t.Status)
		}

		seenIDs[t.ID] = true
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	if state.NextID <= maxID {
		return fmt.Errorf("next_id (%d) должен быть больше максимального ID (%d)", state.NextID, maxID)
	}

	s.nextID = state.NextID
	s.tasks = make([]Task, len(state.Tasks))
	copy(s.tasks, state.Tasks)
	return nil
}

func (s *MemoryStore) Add(title string, priority int) (Task, error) {
	if title == "" {
		return Task{}, errors.New("название не может быть пустым")
	}
	if priority < 1 || priority > 5 {
		return Task{}, errors.New("приоритет должен быть от 1 до 5")
	}
	task := Task{
		ID:       s.nextID,
		Title:    title,
		Status:   "new",
		Priority: priority,
	}
	s.nextID++
	s.tasks = append(s.tasks, task)
	return task, nil
}

func (s *MemoryStore) Complete(id int) error {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].Status = "done"
			return nil
		}
	}
	return errors.New("задача не найдена")
}

func (s *MemoryStore) Delete(id int) error {
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("задача не найдена")
}

func (s *MemoryStore) List() []Task {
	return s.tasks
}

func (s *MemoryStore) GetState() State {
	return State{
		NextID: s.nextID,
		Tasks:  s.tasks,
	}
}