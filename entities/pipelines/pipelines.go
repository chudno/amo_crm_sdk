// Пакет pipelines предоставляет методы для взаимодействия с сущностями "Воронки" в API amoCRM.
package pipelines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chudno/amo_crm_sdk/client"
	"net/http"
)

// Pipeline представляет собой структуру воронки в amoCRM.
type Pipeline struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Sort     int      `json:"sort"`
	IsMain   bool     `json:"is_main"`
	IsActive bool     `json:"is_active"`
	Statuses []Status `json:"statuses,omitempty"`
}

// Status представляет собой структуру статуса в воронке amoCRM.
type Status struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Sort       int    `json:"sort"`
	Color      string `json:"color"`
	Type       int    `json:"type"`
	PipelineID int    `json:"pipeline_id"`
	IsEditable bool   `json:"is_editable"`
}

// GetPipeline получает воронку по её ID.
func GetPipeline(apiClient *client.Client, pipelineID int) (*Pipeline, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines/%d", apiClient.GetBaseURL(), pipelineID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pipeline Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&pipeline); err != nil {
		return nil, err
	}

	return &pipeline, nil
}

// CreatePipeline создает новую воронку в amoCRM.
func CreatePipeline(apiClient *client.Client, pipeline *Pipeline) (*Pipeline, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines", apiClient.GetBaseURL())
	pipelineJSON, err := json.Marshal(pipeline)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pipelineJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var newPipeline Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&newPipeline); err != nil {
		return nil, err
	}

	return &newPipeline, nil
}

// UpdatePipeline обновляет существующую воронку в amoCRM.
func UpdatePipeline(apiClient *client.Client, pipeline *Pipeline) (*Pipeline, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines/%d", apiClient.GetBaseURL(), pipeline.ID)
	pipelineJSON, err := json.Marshal(pipeline)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(pipelineJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updatedPipeline Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&updatedPipeline); err != nil {
		return nil, err
	}

	return &updatedPipeline, nil
}

// ListPipelines получает список воронок.
func ListPipelines(apiClient *client.Client) ([]Pipeline, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines", apiClient.GetBaseURL())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pipelines struct {
		Embedded struct {
			Items []Pipeline `json:"items"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return nil, err
	}

	return pipelines.Embedded.Items, nil
}

// DeletePipeline удаляет воронку по её ID.
func DeletePipeline(apiClient *client.Client, pipelineID int) error {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines/%d", apiClient.GetBaseURL(), pipelineID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetStatus получает статус воронки по его ID.
func GetStatus(apiClient *client.Client, pipelineID int, statusID int) (*Status, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines/%d/statuses/%d", apiClient.GetBaseURL(), pipelineID, statusID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

// CreateStatus создает новый статус в воронке amoCRM.
func CreateStatus(apiClient *client.Client, pipelineID int, status *Status) (*Status, error) {
	url := fmt.Sprintf("%s/api/v4/leads/pipelines/%d/statuses", apiClient.GetBaseURL(), pipelineID)
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(statusJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var newStatus Status
	if err := json.NewDecoder(resp.Body).Decode(&newStatus); err != nil {
		return nil, err
	}

	return &newStatus, nil
}
