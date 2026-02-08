// Package tasks предоставляет методы для взаимодействия с сущностями "Задачи" в API amoCRM.
package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chudno/amo_crm_sdk/client"
	"net/http"
	"net/url"
	"time"
)

// Task представляет собой структуру задачи в amoCRM.
type Task struct {
	ID                int             `json:"id,omitempty"`
	CreatedBy         int             `json:"created_by,omitempty"`
	UpdatedBy         int             `json:"updated_by,omitempty"`
	CreatedAt         int64           `json:"created_at,omitempty"`
	UpdatedAt         int64           `json:"updated_at,omitempty"`
	ResponsibleUserID int             `json:"responsible_user_id,omitempty"`
	GroupID           int             `json:"group_id,omitempty"`
	EntityID          int             `json:"entity_id,omitempty"`
	EntityType        string          `json:"entity_type,omitempty"`
	Duration          int             `json:"duration,omitempty"`
	IsCompleted       bool            `json:"is_completed,omitempty"`
	TaskTypeID        int             `json:"task_type_id,omitempty"`
	Text              string          `json:"text,omitempty"`
	Result            json.RawMessage `json:"result,omitempty"`
	CompleteTill      int64           `json:"complete_till,omitempty"`
	AccountID         int             `json:"account_id,omitempty"`
}

// EntityType определяет тип сущности, к которой привязана задача
const (
	EntityTypeLead     = "leads"
	EntityTypeContact  = "contacts"
	EntityTypeCompany  = "companies"
	EntityTypeCustomer = "customers"
)

// Get получает задачу по её ID.
func Get(ctx context.Context, apiClient *client.Client, taskID int) (*Task, error) {
	url := fmt.Sprintf("%s/api/v4/tasks/%d", apiClient.GetBaseURL(), taskID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

// Create создает новую задачу в amoCRM.
func Create(ctx context.Context, apiClient *client.Client, task *Task) (*Task, error) {
	url := fmt.Sprintf("%s/api/v4/tasks", apiClient.GetBaseURL())

	taskData, err := json.Marshal([]*Task{task})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

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

// Update обновляет существующую задачу в amoCRM.
func Update(ctx context.Context, apiClient *client.Client, task *Task) (*Task, error) {
	if task.ID == 0 {
		return nil, fmt.Errorf("ID задачи не указан")
	}

	url := fmt.Sprintf("%s/api/v4/tasks/%d", apiClient.GetBaseURL(), task.ID)

	taskData, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(taskData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var updatedTask Task
	if err := json.NewDecoder(resp.Body).Decode(&updatedTask); err != nil {
		return nil, err
	}

	return &updatedTask, nil
}

// Complete отмечает задачу как выполненную.
func Complete(ctx context.Context, apiClient *client.Client, taskID int, result string) (*Task, error) {
	resultJSON, _ := json.Marshal(map[string]string{"text": result})
	task := &Task{
		ID:          taskID,
		IsCompleted: true,
		Result:      resultJSON,
	}

	return Update(ctx, apiClient, task)
}

// List получает список задач с возможностью фильтрации и пагинации.
func List(ctx context.Context, apiClient *client.Client, limit int, page int, filter map[string]any) ([]*Task, error) {
	baseURL := fmt.Sprintf("%s/api/v4/tasks", apiClient.GetBaseURL())

	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("page", fmt.Sprintf("%d", page))

	if len(filter) > 0 {
		filterData, err := json.Marshal(filter)
		if err != nil {
			return nil, err
		}
		params.Add("filter", string(filterData))
	}

	url := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

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

// CreateForEntity создает новую задачу, привязанную к сущности (лид, контакт, компания).
func CreateForEntity(ctx context.Context, apiClient *client.Client, entityType string, entityID int, taskTypeID int, text string, completeTill time.Time, responsibleUserID int) (*Task, error) {
	task := &Task{
		EntityType:        entityType,
		EntityID:          entityID,
		TaskTypeID:        taskTypeID,
		Text:              text,
		CompleteTill:      completeTill.Unix(),
		ResponsibleUserID: responsibleUserID,
	}

	return Create(ctx, apiClient, task)
}
