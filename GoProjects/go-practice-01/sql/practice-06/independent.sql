SELECT id, title, priority 
FROM tasks 
WHERE project_id = 1 AND status = 'new' 
ORDER BY priority DESC, id ASC;

SELECT p.name AS project_name, COUNT(t.id) AS task_count
FROM projects p
LEFT JOIN tasks t ON p.id = t.project_id
GROUP BY p.id, p.name
ORDER BY p.id;

-- Добавляем проект
INSERT INTO projects (name) VALUES ('Аналитика');
-- Проверяем (должно быть 3 проекта, у Аналитики 0 задач)
SELECT p.name, COUNT(t.id) 
FROM projects p
LEFT JOIN tasks t ON p.id = t.project_id
GROUP BY p.id, p.name;

-- Обновляем задачу с ID = 1
UPDATE tasks 
SET status = 'done' 
WHERE id = 1;

BEGIN;
-- Вставляем тестовую задачу
INSERT INTO tasks (project_id, title, priority) VALUES (1, 'Тестовая задача', 1);
-- Проверяем (должно стать 5)
SELECT COUNT(*) FROM tasks;
-- Откатываем изменения
ROLLBACK;
-- Проверяем снова (должно снова стать 4)
SELECT COUNT(*) FROM tasks;