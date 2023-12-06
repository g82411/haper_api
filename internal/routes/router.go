package routes

import (
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/routes/paths"
)

func BindingRoutes(app *fiber.App) {
	app.Post("generate", paths.GenerateImage)
	app.Get("auth/callback", paths.AuthCallback)
	app.Post("survey", paths.Survey)
	app.Get("article/:articleId", paths.GetArticle)
}
