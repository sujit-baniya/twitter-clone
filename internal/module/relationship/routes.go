package relationship

import (
	"github.com/sujit-baniya/twitter-clone/internal/common/cache"
	"github.com/sujit-baniya/twitter-clone/internal/common/database"
	"github.com/sujit-baniya/twitter-clone/internal/common/middleware"
	"github.com/sujit-baniya/twitter-clone/internal/module/relationship/action"
	"github.com/sujit-baniya/twitter-clone/internal/module/relationship/service"
	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, db database.Database, cache cache.Cache) {
	authMiddleware := middleware.NewAuthMiddleware()
	r.Post("/follow/:userID", authMiddleware.Execute(), buildFollowUserHandler(db))
	r.Delete("/unfollow/:userID", authMiddleware.Execute(), buildUnfollowUserHandler(db))
	r.Get("/followers/:userID", buildListFollowersHandler(db))
	r.Get("/followings/:userID", buildListFollowingsHandler(db))
}

func buildFollowUserHandler(db database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		service := service.NewFollowUserService(db)
		action := action.NewFollowUserAction(service)

		return action.Execute(c)
	}
}

func buildUnfollowUserHandler(db database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		service := service.NewUnfollowUserService(db)
		action := action.NewUnfollowUserAction(service)

		return action.Execute(c)
	}
}

func buildListFollowersHandler(db database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		service := service.NewListFollowersService(db)
		action := action.NewListFollowersAction(service)

		return action.Execute(c)
	}
}

func buildListFollowingsHandler(db database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		service := service.NewListFollowingsService(db)
		action := action.NewListFollowingsAction(service)

		return action.Execute(c)
	}
}
