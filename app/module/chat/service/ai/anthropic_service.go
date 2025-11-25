// Package ai содержит интеграцию с Anthropic Claude для генерации AI-ответов.
// Реализует tool calling для бронирования столиков, просмотра и отмены бронирований.
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	reservationPayload "savory-ai-server/app/module/reservation/payload"
	reservationService "savory-ai-server/app/module/reservation/service"
	restaurantService "savory-ai-server/app/module/restaurant/service"
	"savory-ai-server/utils/config"
)

// AnthropicService предоставляет интеграцию с Anthropic Claude API.
// Поддерживает генерацию ответов и tool calling для операций с бронированием.
type AnthropicService struct {
	client             anthropic.Client                      // Клиент Anthropic API
	reservationService reservationService.ReservationService // Сервис бронирования для tool calls
	restaurantService  restaurantService.RestaurantService   // Сервис ресторанов для контекста
}

// ChatMessage представляет сообщение в чате для передачи в API.
type ChatMessage struct {
	Role    string `json:"role"`    // Роль: "user" или "assistant"
	Content string `json:"content"` // Текст сообщения
}

// NewAnthropicService создаёт новый экземпляр сервиса Anthropic.
// Принимает конфигурацию с API ключом и зависимые сервисы для tool calling.
func NewAnthropicService(
	cfg *config.Config,
	reservationSvc reservationService.ReservationService,
	restaurantSvc restaurantService.RestaurantService,
) *AnthropicService {
	client := anthropic.NewClient(
		option.WithAPIKey(cfg.Anthropic.APIKey),
	)

	return &AnthropicService{
		client:             client,
		reservationService: reservationSvc,
		restaurantService:  restaurantSvc,
	}
}

// GenerateResponse генерирует ответ AI на основе истории чата.
// Процесс:
//  1. Получает информацию о ресторане для контекста
//  2. Строит системный промпт с инструкциями для AI
//  3. Конвертирует историю сообщений в формат Anthropic
//  4. Отправляет запрос с определёнными tools
//  5. Обрабатывает tool calls, если они есть
//  6. Возвращает текстовый ответ AI
//
// Использует модель claude-3-5-haiku-20241022.
func (s *AnthropicService) GenerateResponse(ctx context.Context, restaurantID uint, messages []ChatMessage) (string, error) {
	// Get restaurant info for context
	restaurant, err := s.restaurantService.GetByID(restaurantID)
	if err != nil {
		return "", fmt.Errorf("failed to get restaurant: %w", err)
	}

	// Build system prompt
	systemPrompt := buildSystemPrompt(restaurant.Name, restaurant.Description, restaurant.ReservationDuration)

	// Convert messages to Anthropic format
	anthropicMessages := make([]anthropic.MessageParam, len(messages))
	for i, msg := range messages {
		if msg.Role == "user" {
			anthropicMessages[i] = anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content))
		} else {
			anthropicMessages[i] = anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content))
		}
	}

	// Define tools
	tools := s.getTools()

	// Create a message with tools
	response, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5Haiku20241022,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{
				Type: "text",
				Text: systemPrompt,
			},
		},
		Messages: anthropicMessages,
		Tools:    tools,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create message: %w", err)
	}

	// Process response - handle tool calls if any
	return s.processResponse(ctx, restaurantID, response, anthropicMessages)
}

// processResponse обрабатывает ответ от Anthropic API.
// Если ответ содержит tool calls, выполняет их и продолжает диалог.
// Рекурсивно обрабатывает результаты tool calls до получения текстового ответа.
func (s *AnthropicService) processResponse(ctx context.Context, restaurantID uint, response *anthropic.Message, messages []anthropic.MessageParam) (string, error) {
	// Check if there are tool uses in the response
	var toolUses []struct {
		ID    string
		Name  string
		Input json.RawMessage
	}
	var textContent string

	for _, block := range response.Content {
		switch b := block.AsAny().(type) {
		case anthropic.ToolUseBlock:
			inputBytes, _ := json.Marshal(b.Input)
			toolUses = append(toolUses, struct {
				ID    string
				Name  string
				Input json.RawMessage
			}{
				ID:    b.ID,
				Name:  b.Name,
				Input: inputBytes,
			})
		case anthropic.TextBlock:
			textContent = b.Text
		}
	}

	// If no tool uses, return text content
	if len(toolUses) == 0 {
		return textContent, nil
	}

	// Process tool calls
	var toolResults []anthropic.ContentBlockParamUnion

	for _, toolUse := range toolUses {
		result, err := s.executeToolCall(restaurantID, toolUse.Name, toolUse.Input)
		if err != nil {
			toolResults = append(toolResults, anthropic.NewToolResultBlock(
				toolUse.ID,
				fmt.Sprintf("Error: %s", err.Error()),
				true,
			))
		} else {
			toolResults = append(toolResults, anthropic.NewToolResultBlock(
				toolUse.ID,
				result,
				false,
			))
		}
	}

	// Build assistant message content blocks
	assistantContentBlocks := make([]anthropic.ContentBlockParamUnion, 0)
	for _, block := range response.Content {
		switch b := block.AsAny().(type) {
		case anthropic.ToolUseBlock:
			assistantContentBlocks = append(assistantContentBlocks,
				anthropic.ContentBlockParamUnion{
					OfToolUse: &anthropic.ToolUseBlockParam{
						Type:  "tool_use",
						ID:    b.ID,
						Name:  b.Name,
						Input: b.Input,
					},
				})
		case anthropic.TextBlock:
			assistantContentBlocks = append(assistantContentBlocks,
				anthropic.NewTextBlock(b.Text))
		}
	}

	// Add assistant message with tool uses and user message with tool results
	newMessages := append(messages,
		anthropic.MessageParam{
			Role:    anthropic.MessageParamRoleAssistant,
			Content: assistantContentBlocks,
		},
	)

	// Add tool results as user message
	newMessages = append(newMessages, anthropic.MessageParam{
		Role:    anthropic.MessageParamRoleUser,
		Content: toolResults,
	})

	// Get restaurant info for system prompt
	restaurant, _ := s.restaurantService.GetByID(restaurantID)
	systemPrompt := buildSystemPrompt(restaurant.Name, restaurant.Description, restaurant.ReservationDuration)

	// Continue conversation with tool results
	continueResponse, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5Haiku20241022,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{
				Type: "text",
				Text: systemPrompt,
			},
		},
		Messages: newMessages,
		Tools:    s.getTools(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to continue conversation: %w", err)
	}

	// Extract text from continue response
	for _, block := range continueResponse.Content {
		if tb, ok := block.AsAny().(anthropic.TextBlock); ok {
			return tb.Text, nil
		}
	}

	return "", nil
}

// executeToolCall выполняет вызов инструмента по его имени.
// Поддерживаемые инструменты:
//   - get_available_slots: получить доступные слоты для бронирования
//   - create_reservation: создать бронирование
//   - get_my_reservations: получить бронирования по телефону
//   - cancel_reservation: отменить бронирование
//   - get_restaurant_info: получить информацию о ресторане
func (s *AnthropicService) executeToolCall(restaurantID uint, toolName string, input json.RawMessage) (string, error) {
	switch toolName {
	case "get_available_slots":
		return s.handleGetAvailableSlots(restaurantID, input)
	case "create_reservation":
		return s.handleCreateReservation(restaurantID, input)
	case "get_my_reservations":
		return s.handleGetMyReservations(input)
	case "cancel_reservation":
		return s.handleCancelReservation(input)
	case "get_restaurant_info":
		return s.handleGetRestaurantInfo(restaurantID)
	default:
		return "", fmt.Errorf("unknown tool: %s", toolName)
	}
}

// GetAvailableSlotsInput параметры для инструмента get_available_slots.
type GetAvailableSlotsInput struct {
	Date       string `json:"date"`        // Дата в формате YYYY-MM-DD
	GuestCount int    `json:"guest_count"` // Количество гостей (по умолчанию 2)
}

// handleGetAvailableSlots обрабатывает запрос на получение доступных слотов.
// Возвращает отформатированный список слотов или сообщение об отсутствии мест.
func (s *AnthropicService) handleGetAvailableSlots(restaurantID uint, input json.RawMessage) (string, error) {
	var params GetAvailableSlotsInput
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	// Set default guest count
	if params.GuestCount == 0 {
		params.GuestCount = 2
	}

	// Use today's date if not specified
	if params.Date == "" {
		params.Date = time.Now().Format("2006-01-02")
	}

	slots, err := s.reservationService.GetAvailableSlots(restaurantID, params.Date, params.GuestCount)
	if err != nil {
		return "", err
	}

	if len(slots.Slots) == 0 {
		return fmt.Sprintf("No available slots on %s for %d guests. The restaurant might be closed on this day or all tables are booked.", params.Date, params.GuestCount), nil
	}

	// Group slots by time for cleaner output
	result := fmt.Sprintf("Available slots on %s for %d guests:\n", params.Date, params.GuestCount)
	for _, slot := range slots.Slots {
		result += fmt.Sprintf("- %s to %s (Table: %s, seats %d guests)\n",
			slot.StartTime, slot.EndTime, slot.TableName, slot.Capacity)
	}

	return result, nil
}

// CreateReservationInput параметры для инструмента create_reservation.
type CreateReservationInput struct {
	Date          string `json:"date"`           // Дата бронирования YYYY-MM-DD
	StartTime     string `json:"start_time"`     // Время начала HH:MM
	GuestCount    int    `json:"guest_count"`    // Количество гостей
	CustomerName  string `json:"customer_name"`  // Имя клиента (обязательно)
	CustomerPhone string `json:"customer_phone"` // Телефон клиента (обязательно)
	CustomerEmail string `json:"customer_email"` // Email клиента (опционально)
	TableID       uint   `json:"table_id"`       // ID столика (автовыбор если не указан)
	Notes         string `json:"notes"`          // Примечания к бронированию
}

// handleCreateReservation обрабатывает запрос на создание бронирования.
// Валидирует входные данные, автоматически подбирает столик если не указан,
// создаёт бронирование и возвращает подтверждение.
func (s *AnthropicService) handleCreateReservation(restaurantID uint, input json.RawMessage) (string, error) {
	var params CreateReservationInput
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	// Validate required fields
	if params.CustomerName == "" {
		return "Please provide the customer's name for the reservation.", nil
	}
	if params.CustomerPhone == "" {
		return "Please provide a phone number for the reservation.", nil
	}
	if params.Date == "" {
		return "Please specify the date for the reservation (format: YYYY-MM-DD).", nil
	}
	if params.StartTime == "" {
		return "Please specify the time for the reservation (format: HH:MM).", nil
	}
	if params.GuestCount == 0 {
		params.GuestCount = 2
	}

	// If no table specified, find one
	if params.TableID == 0 {
		slots, err := s.reservationService.GetAvailableSlots(restaurantID, params.Date, params.GuestCount)
		if err != nil {
			return "", err
		}

		// Find a slot matching the requested time
		for _, slot := range slots.Slots {
			if slot.StartTime == params.StartTime {
				params.TableID = slot.TableID
				break
			}
		}

		if params.TableID == 0 {
			return fmt.Sprintf("Sorry, no table available at %s on %s for %d guests. Please try a different time.", params.StartTime, params.Date, params.GuestCount), nil
		}
	}

	// Create reservation
	req := &reservationPayload.CreateReservationReq{
		RestaurantID:    restaurantID,
		TableID:         params.TableID,
		CustomerName:    params.CustomerName,
		CustomerPhone:   params.CustomerPhone,
		CustomerEmail:   params.CustomerEmail,
		GuestCount:      params.GuestCount,
		ReservationDate: params.Date,
		StartTime:       params.StartTime,
		Notes:           params.Notes,
	}

	reservation, err := s.reservationService.Create(req)
	if err != nil {
		return fmt.Sprintf("Could not create reservation: %s", err.Error()), nil
	}

	return fmt.Sprintf("Reservation confirmed!\n"+
		"- Date: %s\n"+
		"- Time: %s - %s\n"+
		"- Table: %s\n"+
		"- Guests: %d\n"+
		"- Name: %s\n"+
		"- Confirmation ID: %d\n\n"+
		"We look forward to seeing you!",
		reservation.ReservationDate,
		reservation.StartTime,
		reservation.EndTime,
		reservation.TableName,
		reservation.GuestCount,
		reservation.CustomerName,
		reservation.ID,
	), nil
}

// GetMyReservationsInput параметры для инструмента get_my_reservations.
type GetMyReservationsInput struct {
	Phone string `json:"phone"` // Телефон клиента для поиска бронирований
}

// handleGetMyReservations обрабатывает запрос на получение бронирований клиента.
// Ищет бронирования по номеру телефона и возвращает отформатированный список.
func (s *AnthropicService) handleGetMyReservations(input json.RawMessage) (string, error) {
	var params GetMyReservationsInput
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	if params.Phone == "" {
		return "Please provide your phone number to look up your reservations.", nil
	}

	reservations, err := s.reservationService.GetByPhone(params.Phone)
	if err != nil {
		return "", err
	}

	if len(reservations.Reservations) == 0 {
		return fmt.Sprintf("No reservations found for phone number %s.", params.Phone), nil
	}

	result := fmt.Sprintf("Your reservations (phone: %s):\n\n", params.Phone)
	for _, res := range reservations.Reservations {
		result += fmt.Sprintf("ID: %d\n", res.ID)
		result += fmt.Sprintf("- Restaurant: %s\n", res.RestaurantName)
		result += fmt.Sprintf("- Date: %s\n", res.ReservationDate)
		result += fmt.Sprintf("- Time: %s - %s\n", res.StartTime, res.EndTime)
		result += fmt.Sprintf("- Table: %s\n", res.TableName)
		result += fmt.Sprintf("- Guests: %d\n", res.GuestCount)
		result += fmt.Sprintf("- Status: %s\n\n", res.Status)
	}

	return result, nil
}

// CancelReservationInput параметры для инструмента cancel_reservation.
type CancelReservationInput struct {
	ReservationID uint   `json:"reservation_id"` // ID бронирования для отмены
	Phone         string `json:"phone"`          // Телефон для верификации владельца
}

// handleCancelReservation обрабатывает запрос на отмену бронирования.
// Верифицирует владельца по номеру телефона и отменяет бронирование.
func (s *AnthropicService) handleCancelReservation(input json.RawMessage) (string, error) {
	var params CancelReservationInput
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	if params.ReservationID == 0 {
		return "Please provide the reservation ID to cancel.", nil
	}
	if params.Phone == "" {
		return "Please provide your phone number to verify the reservation.", nil
	}

	_, err := s.reservationService.CancelByPhone(params.ReservationID, params.Phone)
	if err != nil {
		return fmt.Sprintf("Could not cancel reservation: %s", err.Error()), nil
	}

	return fmt.Sprintf("Reservation #%d has been cancelled successfully.", params.ReservationID), nil
}

// handleGetRestaurantInfo обрабатывает запрос на получение информации о ресторане.
// Возвращает название, описание, адрес, телефон и рабочие часы.
func (s *AnthropicService) handleGetRestaurantInfo(restaurantID uint) (string, error) {
	restaurant, err := s.restaurantService.GetByID(restaurantID)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Restaurant: %s\n", restaurant.Name)
	if restaurant.Description != "" {
		result += fmt.Sprintf("Description: %s\n", restaurant.Description)
	}
	result += fmt.Sprintf("Address: %s\n", restaurant.Address)
	result += fmt.Sprintf("Phone: %s\n", restaurant.Phone)

	if len(restaurant.WorkingHours) > 0 {
		result += "\nWorking Hours:\n"
		days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
		for _, wh := range restaurant.WorkingHours {
			result += fmt.Sprintf("- %s: %s - %s\n", days[wh.DayOfWeek], wh.OpenTime, wh.CloseTime)
		}
	}

	return result, nil
}

// getTools возвращает определения инструментов для Anthropic API.
// Инструменты:
//   - get_available_slots: проверка доступности столиков
//   - create_reservation: создание бронирования
//   - get_my_reservations: просмотр бронирований по телефону
//   - cancel_reservation: отмена бронирования
//   - get_restaurant_info: информация о ресторане
func (s *AnthropicService) getTools() []anthropic.ToolUnionParam {
	return []anthropic.ToolUnionParam{
		{
			OfTool: &anthropic.ToolParam{
				Name:        "get_available_slots",
				Description: anthropic.String("Get available time slots for table reservations on a specific date. Use this when the customer asks about availability or wants to know what times are free."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type: "object",
					Properties: map[string]interface{}{
						"date": map[string]interface{}{
							"type":        "string",
							"description": "The date to check availability for in YYYY-MM-DD format (e.g., 2024-12-25)",
						},
						"guest_count": map[string]interface{}{
							"type":        "integer",
							"description": "Number of guests for the reservation (default: 2)",
						},
					},
					Required: []string{"date"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "create_reservation",
				Description: anthropic.String("Create a new table reservation. Use this when the customer wants to book a table. You need the customer's name, phone number, date, and time. Ask for any missing information before calling this tool."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type: "object",
					Properties: map[string]interface{}{
						"date": map[string]interface{}{
							"type":        "string",
							"description": "Reservation date in YYYY-MM-DD format",
						},
						"start_time": map[string]interface{}{
							"type":        "string",
							"description": "Reservation time in HH:MM format (24-hour)",
						},
						"guest_count": map[string]interface{}{
							"type":        "integer",
							"description": "Number of guests",
						},
						"customer_name": map[string]interface{}{
							"type":        "string",
							"description": "Customer's full name",
						},
						"customer_phone": map[string]interface{}{
							"type":        "string",
							"description": "Customer's phone number",
						},
						"customer_email": map[string]interface{}{
							"type":        "string",
							"description": "Customer's email (optional)",
						},
						"table_id": map[string]interface{}{
							"type":        "integer",
							"description": "Specific table ID (optional, will auto-select if not provided)",
						},
						"notes": map[string]interface{}{
							"type":        "string",
							"description": "Special requests or notes",
						},
					},
					Required: []string{"date", "start_time", "customer_name", "customer_phone"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "get_my_reservations",
				Description: anthropic.String("Get customer's existing reservations by their phone number. Use this when the customer asks about their reservations or wants to check their bookings."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type: "object",
					Properties: map[string]interface{}{
						"phone": map[string]interface{}{
							"type":        "string",
							"description": "Customer's phone number to look up reservations",
						},
					},
					Required: []string{"phone"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "cancel_reservation",
				Description: anthropic.String("Cancel an existing reservation. The customer must provide their reservation ID and phone number for verification."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type: "object",
					Properties: map[string]interface{}{
						"reservation_id": map[string]interface{}{
							"type":        "integer",
							"description": "The reservation ID to cancel",
						},
						"phone": map[string]interface{}{
							"type":        "string",
							"description": "Customer's phone number for verification",
						},
					},
					Required: []string{"reservation_id", "phone"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "get_restaurant_info",
				Description: anthropic.String("Get information about the restaurant including name, address, phone, and working hours."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
		},
	}
}

// buildSystemPrompt строит системный промпт для AI с контекстом ресторана.
// Включает инструкции по работе с бронированиями, правила общения,
// информацию о ресторане и текущую дату.
func buildSystemPrompt(restaurantName, description string, reservationDuration int) string {
	return fmt.Sprintf(`You are a helpful AI assistant for %s restaurant. Your role is to help customers:
1. Check available time slots for table reservations
2. Make table reservations
3. View their existing reservations
4. Cancel their reservations
5. Answer questions about the restaurant

Restaurant info:
- Name: %s
%s
- Reservation duration: %d minutes

Guidelines:
- Be friendly and professional
- Customers are anonymous - they don't need to be registered in the system
- To identify customers and their reservations, use their phone number
- When a customer wants to book a table, gather all necessary information: date, time, number of guests, name, and phone number
- Use the get_available_slots tool to check availability before suggesting times
- Use the create_reservation tool only when you have all required information
- If the requested time is not available, suggest alternative times
- Confirm reservation details before finalizing
- When a customer asks about their reservations, ask for their phone number and use the get_my_reservations tool
- When a customer wants to cancel, get their phone number and reservation ID, then use the cancel_reservation tool
- Respond in the same language the customer uses (Russian or English)
- Keep responses concise but informative

Today's date: %s`,
		restaurantName,
		restaurantName,
		func() string {
			if description != "" {
				return fmt.Sprintf("- Description: %s", description)
			}
			return ""
		}(),
		reservationDuration,
		time.Now().Format("2006-01-02"),
	)
}
