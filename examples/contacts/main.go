package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/contacts"
	"github.com/chudno/amo_crm_sdk/utils/custom_fields"
)

func main() {
	// Инициализация клиента API
	apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

	// Пример создания нового контакта
	newContact := &contacts.Contact{
		Name:              "Иван Иванов",
		FirstName:         "Иван",
		LastName:          "Иванов",
		ResponsibleUserID: 12345,
		CustomFieldsValues: []custom_fields.CustomFieldValue{
			{
				FieldID: 67890, // ID пользовательского поля "Телефон"
				Values: []custom_fields.FieldValue{
					{
						Value: "+7 (999) 123-45-67",
					},
				},
			},
			{
				FieldID: 54321, // ID пользовательского поля "Email"
				Values: []custom_fields.FieldValue{
					{
						Value: "ivan@example.com",
					},
				},
			},
		},
	}

	createdContact, err := contacts.CreateContact(apiClient, newContact)
	if err != nil {
		log.Fatalf("Ошибка при создании контакта: %v", err)
	}
	fmt.Printf("Создан контакт: ID=%d, Имя=%s\n", createdContact.ID, createdContact.Name)

	// Пример получения контакта по ID
	contactID := createdContact.ID // или любой существующий ID контакта
	contact, err := contacts.GetContact(apiClient, contactID)
	if err != nil {
		log.Fatalf("Ошибка при получении контакта: %v", err)
	}
	fmt.Printf("Получен контакт: ID=%d, Имя=%s\n", contact.ID, contact.Name)

	// Пример получения списка контактов
	contactsList, err := contacts.GetContacts(apiClient, 1, 50) // страница 1, лимит 50
	if err != nil {
		log.Fatalf("Ошибка при получении списка контактов: %v", err)
	}
	fmt.Printf("Получено %d контактов\n", len(contactsList))
	for i, c := range contactsList {
		if i < 5 { // Выводим только первые 5 контактов
			fmt.Printf("Контакт %d: ID=%d, Имя=%s\n", i+1, c.ID, c.Name)
		}
	}
}
