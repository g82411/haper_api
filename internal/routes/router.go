package routes

import (
	"github.com/gofiber/fiber/v2"
	"hyper_api/internal/middleware"
	"hyper_api/internal/routes/paths"
)

func BindingRoutes(app *fiber.App) {
	app.Get("userInfo", middleware.AuthMiddleware, paths.GetUserInfo)
	app.Post("generate", middleware.AuthMiddleware, paths.GenerateImage)
	app.Post("survey", middleware.AuthMiddleware, paths.Survey)
	app.Get("user/article", middleware.AuthMiddleware, paths.QueryArticleByUser)

	app.Get("article/:articleId", paths.GetArticle)
	app.Patch("article/:articleId", middleware.AuthMiddleware, paths.UpdateArticle)
	app.Get("articles", paths.TakeImages)
	app.Get("templates", paths.GetTemplates)
}
