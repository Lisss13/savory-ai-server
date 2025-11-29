// Package controller содержит HTTP обработчики для модуля поддержки.
package controller

// Controller агрегирует все контроллеры модуля поддержки
type Controller struct {
	Support SupportController
}

// NewControllers создаёт новый экземпляр агрегатора контроллеров
func NewControllers(support SupportController) *Controller {
	return &Controller{Support: support}
}
