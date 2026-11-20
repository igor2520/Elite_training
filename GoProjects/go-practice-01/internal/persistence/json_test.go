package persistence

import (
	"os"
	"path/filepath"
	"testing"

	"example.com/go-practice-01/internal/tasks"
)

func TestMissingFileIsNewStore(t *testing.T) {
	state, err := Load(filepath.Join(t.TempDir(), "absent.json"))
	if err != nil || state.NextID != 1 || len(state.Tasks) != 0 {
		t.Fatalf("ожидалось пустое состояние, получено %+v, ошибка: %v", state, err)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	initialState := tasks.State{
		NextID: 3,
		Tasks: []tasks.Task{
			{ID: 1, Title: "Изучить Go", Status: "done", Priority: 2},
			{ID: 2, Title: "Написать API", Status: "new", Priority: 3},
		},
	}

	if err := Save(path, initialState); err != nil {
		t.Fatalf("ошибка сохранения: %v", err)
	}

	loadedState, err := Load(path)
	if err != nil {
		t.Fatalf("ошибка загрузки: %v", err)
	}

	if loadedState.NextID != initialState.NextID {
		t.Errorf("NextID не совпадает: ожидалось %d, получено %d", initialState.NextID, loadedState.NextID)
	}
	if len(loadedState.Tasks) != len(initialState.Tasks) {
		t.Fatalf("количество задач не совпадает: ожидалось %d, получено %d", len(initialState.Tasks), len(loadedState.Tasks))
	}
	for i := range initialState.Tasks {
		if loadedState.Tasks[i] != initialState.Tasks[i] {
			t.Errorf("Задача %d не совпадает: ожидалось %+v, получено %+v", i, initialState.Tasks[i], loadedState.Tasks[i])
		}
	}
}

func TestLoadScenarios(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "пустой файл",
			content:     "",
			expectError: true,
		},
		{
			name:        "битый JSON",
			content:     "{oops",
			expectError: true,
		},
		{
			name:        "два JSON объекта подряд",
			content:     `{"next_id": 1, "tasks": []}{"next_id": 2, "tasks": []}`,
			expectError: true,
		},
		{
			name:        "неизвестное поле в JSON",
			content:     `{"next_id": 1, "tasks": [], "unknown_field": 123}`,
			expectError: true,
		},
		{
			name:        "некорректный next_id",
			content:     `{"next_id": 0, "tasks": []}`,
			expectError: true,
		},
		{
			name:        "дубликаты ID",
			content:     `{"next_id": 3, "tasks": [{"id": 1, "title": "A", "status": "new", "priority": 1}, {"id": 1, "title": "B", "status": "new", "priority": 1}]}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "test.json")

			// Записываем тестовые данные напрямую в файл, минуя функцию Save
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("не удалось создать тестовый файл: %v", err)
			}

			_, err := Load(path)
			if tt.expectError && err == nil {
				t.Errorf("ожидалась ошибка, но ее нет")
			}
			if !tt.expectError && err != nil {
				t.Errorf("не ожидалась ошибка, но получена: %v", err)
			}
		})
	}
}