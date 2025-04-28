package main

import (
	"fmt"
	"log"
	"time"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/tasks"
)

func main() {
	// Инициализация клиента API
	apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

	// Пример создания новой задачи
	// Время выполнения задачи - через 3 дня
	completeTill := time.Now().AddDate(0, 0, 3).Unix()

	newTask := &tasks.Task{
		Text:              "Позвонить клиенту",
		TaskTypeID:        1, // 1 - звонок
		ResponsibleUserID: 12345,
		EntityID:          67890,   // ID сущности (например, ID сделки)
		EntityType:        "leads", // Тип сущности (leads - сделка)
		CompleteTill:      completeTill,
	}

	createdTask, err := tasks.CreateTask(apiClient, newTask)
	if err != nil {
		log.Fatalf("Ошибка при создании задачи: %v", err)
	}
	fmt.Printf("Создана задача: ID=%d, Текст=%s\n", createdTask.ID, createdTask.Text)

	// Пример получения задачи по ID
	taskID := createdTask.ID // или любой существующий ID задачи
	task, err := tasks.GetTask(apiClient, taskID)
	if err != nil {
		log.Fatalf("Ошибка при получении задачи: %v", err)
	}
	fmt.Printf("Получена задача: ID=%d, Текст=%s\n", task.ID, task.Text)

	// Пример обновления задачи
	taskToUpdate := &tasks.Task{
		ID:           taskID,
		Text:         "Позвонить клиенту и обсудить детали сделки",
		CompleteTill: time.Now().AddDate(0, 0, 5).Unix(), // Перенос на 5 дней вперед
	}

	updatedTask, err := tasks.UpdateTask(apiClient, taskToUpdate)
	if err != nil {
		log.Fatalf("Ошибка при обновлении задачи: %v", err)
	}
	fmt.Printf("Задача обновлена: ID=%d, Новый текст=%s\n", updatedTask.ID, updatedTask.Text)

	// Пример получения списка задач с фильтрацией
	filter := map[string]interface{}{
		"responsible_user_id": 12345,
		"task_type_id":        1, // 1 - звонок
	}

	tasksList, err := tasks.ListTasks(apiClient, 50, 1, filter) // лимит 50, страница 1
	if err != nil {
		log.Fatalf("Ошибка при получении списка задач: %v", err)
	}
	fmt.Printf("Получено %d задач\n", len(tasksList))
	for i, t := range tasksList {
		if i < 5 { // Выводим только первые 5 задач
			fmt.Printf("Задача %d: ID=%d, Текст=%s\n", i+1, t.ID, t.Text)
		}
	}

	// Пример выполнения задачи
	completeResult, err := tasks.CompleteTask(apiClient, taskID, "Звонок выполнен, клиент проинформирован")
	if err != nil {
		log.Fatalf("Ошибка при выполнении задачи: %v", err)
	}
	fmt.Printf("Задача выполнена: ID=%d, Статус=%v, Результат=%s\n", completeResult.ID, completeResult.IsCompleted, completeResult.Result)
}
