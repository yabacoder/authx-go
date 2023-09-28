package controllers

import (
	// "authx/initializers"
	// "authx/model"
	// "fmt"
	"authx/model"
	"fmt"

	// "os/user"

	middleware "authx/middleware"

	"github.com/gofiber/fiber/v2"
)

type Result struct {
	ID        uint64 `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Photo     string `json:"photo"`
}

func convertSearchToJSON(user model.User) Result {
	// search { Term, model.User}
	return Result{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Photo:     user.Photo,
	}
}

func Search(c *fiber.Ctx) error {
	c.Accepts("applicatin/json")

	// retrieve the search term
	// search from the term
	// save the search term with the userId
	// var search model.Search
	// err := c.BodyParser(&search)
	term := c.Params("term")
	// if err != nil {
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"message": "error parsing JSON " + err.Error(),
	// 	})
	// }

	// fmt.Println(term)
	userResp, err := model.SearchTerm(term)
	users := []Result{}

	// resp := convertToJSON(userResp)
	for _, user := range userResp {
		fmt.Println(user)
		// search.User = user
		conv := convertSearchToJSON(user)
		users = append(users, conv)
	}

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "user not found " + err.Error(),
		})
	}
	loggedInUser := middleware.GetLoggedInUser(c)

	err = model.SaveUserSearch(term, loggedInUser)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "error logging search " + err.Error(),
		})
	}

	// rst :=  Result { term, rs}

	// resp := convertSearchToJSON(search)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"query":  term,
		"data":   users,
	})
}
