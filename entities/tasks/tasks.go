// Пакет tasks предоставляет методы для взаимодействия с сущностями "Задачи" в API amoCRM.
package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chudno/amo_crm_sdk/client"
	"net/http"
	"net/url"
	"time"
)

// Task представляет собой структуру задачи в amoCRM.
type Task struct {
	ID                int    `json:"id,omitempty"`
	CreatedBy         int    `json:"created_by,omitempty"`
	UpdatedBy         int    `json:"updated_by,omitempty"`
	CreatedAt         int64  `json:"created_at,omitempty"`
	UpdatedAt         int64  `json:"updated_at,omitempty"`
	ResponsibleUserID int    `json:"responsible_user_id,omitempty"`
	GroupID           int    `json:"group_id,omitempty"`
	EntityID          int    `json:"entity_id,omitempty"`
	EntityType        string `json:"entity_type,omitempty"`
	Duration          int    `json:"duration,omitempty"`
	IsCompleted       bool   `json:"is_completed,omitempty"`
	TaskTypeID        int    `json:"task_type_id,omitempty"`
	Text              string `json:"text,omitempty"`
	Result            string `json:"result,omitempty"`
	CompleteTill      int64  `json:"complete_till,omitempty"`
	AccountID         int    `json:"account_id,omitempty"`
}

// EntityType определяет тип сущности, к которой привязана задача
const (
	EntityTypeLead     = "leads"
	EntityTypeContact  = "contacts"
	EntityTypeCompany  = "companies"
	EntityTypeCustomer = "customers"
)

// GetTask получает задачу по её ID.
func GetTask(apiClient *client.Client, taskID int) (*Task, error) {
	url := fmt.Sprintf("%s/api/v4/tasks/%d", apiClient.GetBaseURL(), taskID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

// CreateTask создает новую задачу в amoCRM.
func CreateTask(apiClient *client.Client, task *Task) (*Task, error) {
	url := fmt.Sprintf("%s/api/v4/tasks", apiClient.GetBaseURL())

	taskData, err := json.Marshal([]*Task{task})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Embedded struct {
			Tasks []*Task `json:"tasks"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Embedded.Tasks) == 0 {
		return nil, fmt.Errorf("не удалось создать задачу")
	}

	return response.Embedded.Tasks[0], nil
}

// UpdateTask обновляет существующую задачу в amoCRM.
func UpdateTask(apiClient *client.Client, task *Task) (*Task, error) {
	if task.ID == 0 {
		return nil, fmt.Errorf("ID задачи не указан")
	}

	url := fmt.Sprintf("%s/api/v4/tasks/%d", apiClient.GetBaseURL(), task.ID)

	taskData, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updatedTask Task
	if err := json.NewDecoder(resp.Body).Decode(&updatedTask); err != nil {
		return nil, err
	}

	return &updatedTask, nil
}

// CompleteTask отмечает задачу как выполненную.
func CompleteTask(apiClient *client.Client, taskID int, result string) (*Task, error) {
	task := &Task{
		ID:          taskID,
		IsCompleted: true,
		Result:      result,
	}

	return UpdateTask(apiClient, task)
}

// ListTasks получает список задач с возможностью фильтрации и пагинации.
func ListTasks(apiClient *client.Client, limit int, page int, filter map[string]interface{}) ([]*Task, error) {
	baseURL := fmt.Sprintf("%s/api/v4/tasks", apiClient.GetBaseURL())

	// Добавляем параметры запроса
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("page", fmt.Sprintf("%d", page))

	// Если указаны фильтры, добавляем их в запрос
	if len(filter) > 0 {
		filterData, err := json.Marshal(filter)
		if err != nil {
			return nil, err
		}
		params.Add("filter", string(filterData))
	}

	url := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Embedded struct {
			Tasks []*Task `json:"tasks"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Tasks, nil
}

// DeleteTask удаляет задачу по её ID.
func DeleteTask(apiClient *client.Client, taskID int) error {
	url := fmt.Sprintf("%s/api/v4/tasks/%d", apiClient.GetBaseURL(), taskID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// CreateTaskForEntity создает новую задачу, привязанную к сущности (лид, контакт, компания).
func CreateTaskForEntity(apiClient *client.Client, entityType string, entityID int, taskTypeID int, text string, completeTill time.Time, responsibleUserID int) (*Task, error) {
	task := &Task{
		EntityType:        entityType,
		EntityID:          entityID,
		TaskTypeID:        taskTypeID,
		Text:              text,
		CompleteTill:      completeTill.Unix(),
		ResponsibleUserID: responsibleUserID,
	}

	return CreateTask(apiClient, task)
}
