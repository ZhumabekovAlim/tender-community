package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, makeResponseJSON)

	dynamicMiddleware := alice.New()

	mux := pat.New()

	// USERS
	mux.Post("/users/signup", dynamicMiddleware.ThenFunc(app.userHandler.SignUp))                    // sign up
	mux.Post("/users/login", dynamicMiddleware.ThenFunc(app.userHandler.LogIn))                      // login
	mux.Get("/users", standardMiddleware.ThenFunc(app.userHandler.GetAllUsers))                      // get all users
	mux.Get("/users/details/:id", standardMiddleware.ThenFunc(app.userHandler.GetUserByID))          // get one user info http://localhost:4000/clients/details/1
	mux.Get("/users/balance/:id", standardMiddleware.ThenFunc(app.userHandler.GetBalance))           // get user balance by id
	mux.Put("/users/balance/update/:id", standardMiddleware.ThenFunc(app.userHandler.UpdateBalance)) // update user balance

	// PERMISSIONS
	mux.Post("/permissions", dynamicMiddleware.ThenFunc(app.permissionHandler.AddPermission))                        // add a new permission
	mux.Get("/permissions/user/:user_id", standardMiddleware.ThenFunc(app.permissionHandler.GetPermissionsByUserID)) // get all permissions by user ID
	mux.Put("/permissions/:id", standardMiddleware.ThenFunc(app.permissionHandler.UpdatePermission))                 // update a permission by id
	mux.Del("/permissions/:id", standardMiddleware.ThenFunc(app.permissionHandler.DeletePermission))                 // delete a permission by id

	// COMPANY
	mux.Post("/companies", dynamicMiddleware.ThenFunc(app.companyHandler.CreateCompany))      // Create a new company
	mux.Get("/companies", standardMiddleware.ThenFunc(app.companyHandler.GetAllCompanies))    // Get all companies
	mux.Get("/companies/:id", standardMiddleware.ThenFunc(app.companyHandler.GetCompanyByID)) // Get company by ID
	mux.Put("/companies/:id", standardMiddleware.ThenFunc(app.companyHandler.UpdateCompany))  // Update company by ID
	mux.Del("/companies/:id", standardMiddleware.ThenFunc(app.companyHandler.DeleteCompany))  // Delete company by ID

	// TRANSACTION ROUTES
	mux.Post("/transactions", dynamicMiddleware.ThenFunc(app.transactionHandler.CreateTransaction))      // Create a new transaction
	mux.Get("/transactions", standardMiddleware.ThenFunc(app.transactionHandler.GetAllTransactions))     // Get all transactions
	mux.Get("/transactions/:id", standardMiddleware.ThenFunc(app.transactionHandler.GetTransactionByID)) // Get transaction by ID
	mux.Put("/transactions/:id", standardMiddleware.ThenFunc(app.transactionHandler.UpdateTransaction))  // Update transaction by ID
	mux.Del("/transactions/:id", standardMiddleware.ThenFunc(app.transactionHandler.DeleteTransaction))  // Delete transaction by ID

	// PERSONAL EXPENSES ROUTES
	mux.Post("/expenses", dynamicMiddleware.ThenFunc(app.expenseHandler.CreatePersonalExpense))      // Create a new expense
	mux.Get("/expenses", standardMiddleware.ThenFunc(app.expenseHandler.GetAllPersonalExpenses))     // Get all expenses
	mux.Get("/expenses/:id", standardMiddleware.ThenFunc(app.expenseHandler.GetPersonalExpenseByID)) // Get expense by ID
	mux.Put("/expenses/:id", standardMiddleware.ThenFunc(app.expenseHandler.UpdatePersonalExpense))  // Update expense by ID
	mux.Del("/expenses/:id", standardMiddleware.ThenFunc(app.expenseHandler.DeletePersonalExpense))  // Delete expense by ID

	return standardMiddleware.Then(mux)
}
