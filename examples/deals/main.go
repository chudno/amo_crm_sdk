package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/deals"
)

func main() {
	// Инициализация клиента API
	apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

	// Пример создания новой сделки
	newDeal := &deals.Deal{
		Name:              "Новая сделка",
		Value:             10000,
		ResponsibleUserID: 12345,
		PipelineID:        67890,
		StatusID:          54321,
	}

	createdDeal, err := deals.CreateDeal(apiClient, newDeal)
	if err != nil {
		log.Fatalf("Ошибка при создании сделки: %v", err)
	}
	fmt.Printf("Создана сделка: ID=%d, Имя=%s, Сумма=%d\n", createdDeal.ID, createdDeal.Name, createdDeal.Value)

	// Пример получения сделки по ID
	dealID := createdDeal.ID // или любой существующий ID сделки
	deal, err := deals.GetDeal(apiClient, dealID)
	if err != nil {
		log.Fatalf("Ошибка при получении сделки: %v", err)
	}
	fmt.Printf("Получена сделка: ID=%d, Имя=%s, Сумма=%d\n", deal.ID, deal.Name, deal.Value)

	// Пример обновления сделки (изменение статуса)
	dealToUpdate := &deals.Deal{
		ID:       dealID,
		StatusID: 98765, // ID нового статуса
	}

	updatedDeal, err := deals.UpdateDeal(apiClient, dealToUpdate)
	if err != nil {
		log.Fatalf("Ошибка при обновлении сделки: %v", err)
	}
	fmt.Printf("Сделка обновлена: ID=%d, Новый статус=%d\n", updatedDeal.ID, updatedDeal.StatusID)

	// Пример получения списка сделок с фильтрацией по ответственному
	filter := map[string]string{
		"responsible_user_id": "12345",
		"status_id":           "54321",
	}

	dealsList, err := deals.GetDeals(apiClient, 1, 50, filter) // страница 1, лимит 50
	if err != nil {
		log.Fatalf("Ошибка при получении списка сделок: %v", err)
	}
	fmt.Printf("Получено %d сделок\n", len(dealsList))
	for i, d := range dealsList {
		if i < 5 { // Выводим только первые 5 сделок
			fmt.Printf("Сделка %d: ID=%d, Имя=%s, Сумма=%d\n", i+1, d.ID, d.Name, d.Value)
		}
	}
}
