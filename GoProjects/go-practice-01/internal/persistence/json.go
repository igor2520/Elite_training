package persistence

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"example.com/go-practice-01/internal/tasks"
)

// Load читает состояние из файла. Если файла нет — возвращает пустое состояние.
func Load(path string) (tasks.State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return tasks.State{NextID: 1}, nil
		}
		return tasks.State{}, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	var state tasks.State
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&state); err != nil {
		return tasks.State{}, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	// Проверка, что в файле ровно один JSON-документ
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return tasks.State{}, errors.New("в файле больше одного JSON-объекта")
	}

	// Проверяем корректность данных
	store := tasks.NewMemoryStore()
	if err := store.Restore(state); err != nil {
		return tasks.State{}, fmt.Errorf("некорректные данные в файле: %w", err)
	}

	return state, nil
}

// Save сохраняет состояние в файл (через временный файл)
func Save(path string, state tasks.State) error {
	// Валидация перед сохранением
	store := tasks.NewMemoryStore()
	if err := store.Restore(state); err != nil {
		return fmt.Errorf("нельзя сохранить некорректное состояние: %w", err)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}
	data = append(data, '\n') // Добавляем перевод строки в конце

	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, ".tasks-*.tmp")
	if err != nil {
		return fmt.Errorf("ошибка создания временного файла: %w", err)
	}
	
	// В случае ошибки удаляем временный файл
	defer func() {
		if err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}
	}()

	if _, err = tmpFile.Write(data); err != nil {
		return fmt.Errorf("ошибка записи: %w", err)
	}
	if err = tmpFile.Sync(); err != nil {
		return fmt.Errorf("ошибка синхронизации: %w", err)
	}
	if err = tmpFile.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия файла: %w", err)
	}

	if err = os.Rename(tmpFile.Name(), path); err != nil {
		return fmt.Errorf("ошибка замены файла: %w", err)
	}

	return nil
}