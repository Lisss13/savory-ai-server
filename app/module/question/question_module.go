package question

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/question/controller"
	question_repo "savory-ai-server/app/module/question/repository"
	"savory-ai-server/app/module/question/service"
)

type QuestionRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewQuestionRouter(fiber *fiber.App, controller *controller.Controller) *QuestionRouter {
	return &QuestionRouter{
		App:        fiber,
		Controller: controller,
	}
}

var QuestionModule = fx.Options(
	fx.Provide(question_repo.NewQuestionRepository),
	fx.Provide(service.NewQuestionService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewQuestionRouter),
)

func (r *QuestionRouter) RegisterQuestionRoutes(auth fiber.Handler) {
	questionController := r.Controller.Question
	r.App.Route("/questions", func(router fiber.Router) {
		router.Get("/", auth, questionController.GetAll)
		router.Post("/", auth, questionController.Create)
		router.Delete("/:id", auth, questionController.Delete)
	})
}
