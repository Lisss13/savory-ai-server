// Package ai_reservation предоставляет AI-сервис для бронирования столиков.
// Интегрируется с Anthropic Claude для обработки запросов на естественном языке.
//
// Поддерживаемые операции через tool calling:
//   - get_available_slots: проверка доступности столиков
//   - create_reservation: создание бронирования
//   - get_my_reservations: просмотр бронирований по телефону
//   - cancel_reservation: отмена бронирования
//   - get_restaurant_info: информация о ресторане
package ai_reservation

import (
	"go.uber.org/fx"
	"savory-ai-server/app/module/ai_reservation/service"
)

// AIReservationModule определяет FX-модуль для AI-бронирования.
// Регистрирует сервис AIReservationService с зависимостями от:
//   - config.Config (API ключ Anthropic)
//   - ReservationService (операции с бронированием)
//   - RestaurantService (информация о ресторане)
var AIReservationModule = fx.Options(
	fx.Provide(service.NewAIReservationService),
)
