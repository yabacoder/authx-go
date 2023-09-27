package controllers

import (
	// "authx/initializers"
	// "authx/model"
	// "fmt"
	"authx/model"
	// "os/user"

	"github.com/gofiber/fiber/v2"

)

type Result struct {
	Term string	`json:"term"`
	User model.User	`json:"user"`
}

func convertSearchToJSON(search model.Search) Result {
	// search { Term, model.User}
	return  Result {
		Term: search.Term,
		User: search.User,
	}
}

func Search (c *fiber.Ctx) error {
	c.Accepts("applicatin/json")

	// retrieve the search term
	// search from the term
	// save the search term with the userId
	var search model.Search
	// err := c.BodyParser(&search)
	term := c.Params("search")
	// if err != nil {
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"message": "error parsing JSON " + err.Error(),
	// 	})
	// }

	_, err := model.SearchTerm(term)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "user not found " + err.Error(),
		})
	}
	err = model.SaveUserSearch(term, "1")

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "error logging search " + err.Error(),
		})
	}

	// rst :=  Result { term, rs}
	
	resp := convertSearchToJSON(search)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": resp, 
	})
}