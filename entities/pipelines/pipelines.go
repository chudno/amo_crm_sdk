// Package pipelines предоставляет методы для взаимодействия с сущностями "Воронки" в API amoCRM.
package pipelines

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chudno/amo_crm_sdk/client"
)

// Pipeline представляет собой структуру воронки в amoCRM.
type Pipeline struct {
	ID           int      `json:"id,omitempty"`
	Name         string   `json:"name"`
	Sort         int      `json:"sort"`
	IsMain       bool     `json:"is_main"`
	IsUnsortedOn bool     `json:"is_unsorted_on,omitempty"`
	IsActive     bool     `json:"is_active,omitempty"`
	AccountID    int      `json:"account_id,omitempty"`
	Statuses     []Status `json:"statuses,omitempty"`
}

// StatusType представляет тип статуса в воронке.
// API может возвращать это поле как строку ("normal", "win", "lose", "unsorted")
// или как число (0, 1). StatusType обрабатывает оба варианта.
type StatusType string

const (
	// StatusTypeNormal — обычный этап воронки
	StatusTypeNormal StatusType = "normal"
	// StatusTypeUnsorted — неразобранное
	StatusTypeUnsorted StatusType = "unsorted"
	// StatusTypeWin — успешно реализовано
	StatusTypeWin StatusType = "win"
	// StatusTypeLose — закрыто и не реализовано
	StatusTypeLose StatusType = "lose"
)

// UnmarshalJSON обрабатывает как строковые ("normal"), так и числовые (0) значения type.
func (s *StatusType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = StatusType(str)
		return nil
	}

	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		switch num {
		case 0:
			*s = StatusTypeNormal
		case 1:
			*s = StatusTypeUnsorted
		default:
			*s = StatusType(fmt.Sprintf("%d", num))
		}
		return nil
	}

	return fmt.Errorf("cannot unmarshal status type: %s", string(data))
}

// Status представляет собой структуру статуса в воронке amoCRM.
type Status struct {
	ID         int        `json:"id,omitempty"`
	Name       string     `json:"name"`
	Sort       int        `json:"sort"`
	Color      string     `json:"color,omitempty"`
	Type       StatusType `json:"type,omitempty"`
	PipelineID int        `json:"pipeline_id,omitempty"`
	IsEditable bool       `json:"is_editable,omitempty"`
}

// Get получает воронку по её ID.
func Get(ctx context.Context, apiClient *client.Client, pipelineID int) (*Pipeline, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines/%d", apiClient.GetBaseURL(), pipelineID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var pipeline Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&pipeline); err != nil {
		return nil, err
	}

	return &pipeline, nil
}

// Create создает новую воронку в amoCRM.
func Create(ctx context.Context, apiClient *client.Client, pipeline *Pipeline) (*Pipeline, error) {
	requestURL := fmt.Sprintf("%s/api/v4/leads/pipelines", apiClient.GetBaseURL())

	// API ожидает массив объектов со статусами внутри _embedded
	type createRequest struct {
		Name         string `json:"name"`
		Sort         int    `json:"sort"`
		IsMain       bool   `json:"is_main"`
		IsUnsortedOn bool   `json:"is_unsorted_on,omitempty"`
		Embedded     struct {
			Statuses []Status `json:"statuses"`
		} `json:"_embedded"`
	}

	createReq := createRequest{
		Name:         pipeline.Name,
		Sort:         pipeline.Sort,
		IsMain:       pipeline.IsMain,
		IsUnsortedOn: pipeline.IsUnsortedOn,
	}
	createReq.Embedded.Statuses = pipeline.Statuses
	if createReq.Embedded.Statuses == nil {
		createReq.Embedded.Statuses = []Status{}
	}

	pipelineJSON, err := json.Marshal([]createRequest{createReq})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(pipelineJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("неожиданный статус-код: %d, ответ: %s", resp.StatusCode, body)
	}

	var response struct {
		Embedded struct {
			Pipelines []Pipeline `json:"pipelines"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if len(response.Embedded.Pipelines) == 0 {
		return nil, fmt.Errorf("не удалось создать воронку")
	}
	return &response.Embedded.Pipelines[0], nil
}

// Update обновляет существующую воронку в amoCRM.
func Update(ctx context.Context, apiClient *client.Client, pipeline *Pipeline) (*Pipeline, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines/%d", apiClient.GetBaseURL(), pipeline.ID)
	pipelineJSON, err := json.Marshal(pipeline)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(pipelineJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var updatedPipeline Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&updatedPipeline); err != nil {
		return nil, err
	}

	return &updatedPipeline, nil
}

// List получает список воронок.
func List(ctx context.Context, apiClient *client.Client) ([]Pipeline, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines", apiClient.GetBaseURL())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var pipelines struct {
		Embedded struct {
			Items []Pipeline `json:"pipelines"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return nil, err
	}

	return pipelines.Embedded.Items, nil
}

// Delete удаляет воронку по её ID.
func Delete(ctx context.Context, apiClient *client.Client, pipelineID int) error {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines/%d", apiClient.GetBaseURL(), pipelineID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetStatus получает статус воронки по его ID.
func GetStatus(ctx context.Context, apiClient *client.Client, pipelineID int, statusID int) (*Status, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines/%d/statuses/%d", apiClient.GetBaseURL(), pipelineID, statusID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var status Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

// CreateStatus создает новый статус в воронке amoCRM.
func CreateStatus(ctx context.Context, apiClient *client.Client, pipelineID int, status *Status) (*Status, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines/%d/statuses", apiClient.GetBaseURL(), pipelineID)
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(statusJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var newStatus Status
	if err := json.NewDecoder(resp.Body).Decode(&newStatus); err != nil {
		return nil, err
	}

	return &newStatus, nil
}
