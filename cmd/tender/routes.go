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
	mux.Del("/users/:id", standardMiddleware.ThenFunc(app.userHandler.DeleteUserByID))               // delete user by id
	mux.Put("/users/:id", standardMiddleware.ThenFunc(app.userHandler.UpdateUser))                   // update user by id
	mux.Get("/users/balance/:id", standardMiddleware.ThenFunc(app.userHandler.GetBalance))           // get user balance by id
	mux.Put("/users/balance/update/:id", standardMiddleware.ThenFunc(app.userHandler.UpdateBalance)) // update user balance
	mux.Put("/users/password/:id", standardMiddleware.ThenFunc(app.userHandler.ChangePassword))      // update user balance

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

	// TRANSACTION
	mux.Post("/transactions", dynamicMiddleware.ThenFunc(app.transactionHandler.CreateTransaction))              // Create a new transaction
	mux.Get("/transactions", standardMiddleware.ThenFunc(app.transactionHandler.GetAllTransactions))             // Get all transactions
	mux.Get("/transactions/:id", standardMiddleware.ThenFunc(app.transactionHandler.GetTransactionByID))         // Get transaction by ID
	mux.Get("/transactions/user/:id", standardMiddleware.ThenFunc(app.transactionHandler.GetTransactionsByUser)) // Get transaction by user ID
	mux.Put("/transactions/:id", standardMiddleware.ThenFunc(app.transactionHandler.UpdateTransaction))          // Update transaction by ID
	mux.Del("/transactions/:id", standardMiddleware.ThenFunc(app.transactionHandler.DeleteTransaction))          // Delete transaction by ID

	// EXTRA TRANSACTIONS
	mux.Post("/extra_transactions", dynamicMiddleware.ThenFunc(app.extraTransactionHandler.CreateExtraTransaction))              // Create a new extra transaction
	mux.Get("/extra_transactions", standardMiddleware.ThenFunc(app.extraTransactionHandler.GetAllExtraTransactions))             // Get all extra transactions
	mux.Get("/extra_transactions/:id", standardMiddleware.ThenFunc(app.extraTransactionHandler.GetExtraTransactionByID))         // Get extra transaction by ID
	mux.Get("/extra_transactions/user/:id", standardMiddleware.ThenFunc(app.extraTransactionHandler.GetExtraTransactionsByUser)) // Get extra transactions by user ID
	mux.Put("/extra_transactions/:id", standardMiddleware.ThenFunc(app.extraTransactionHandler.UpdateExtraTransaction))          // Update extra transaction by ID
	mux.Del("/extra_transactions/:id", standardMiddleware.ThenFunc(app.extraTransactionHandler.DeleteExtraTransaction))          // Delete extra transaction by ID

	// PERSONAL EXPENSES
	mux.Post("/expenses", dynamicMiddleware.ThenFunc(app.expenseHandler.CreatePersonalExpense))      // Create a new expense
	mux.Get("/expenses", standardMiddleware.ThenFunc(app.expenseHandler.GetAllPersonalExpenses))     // Get all expenses
	mux.Get("/expenses/:id", standardMiddleware.ThenFunc(app.expenseHandler.GetPersonalExpenseByID)) // Get expense by ID
	mux.Put("/expenses/:id", standardMiddleware.ThenFunc(app.expenseHandler.UpdatePersonalExpense))  // Update expense by ID
	mux.Del("/expenses/:id", standardMiddleware.ThenFunc(app.expenseHandler.DeletePersonalExpense))  // Delete expense by ID

	// REPORTS
	// company month
	mux.Get("/reports/company/month/global", standardMiddleware.ThenFunc(app.transactionHandler.GetMonthlyAmountsByGlobal))               //global - company - month
	mux.Get("/reports/company/month/year", standardMiddleware.ThenFunc(app.transactionHandler.GetMonthlyAmountsByYear))                   //year - company - month
	mux.Get("/reports/company/month/company", standardMiddleware.ThenFunc(app.transactionHandler.GetMonthlyAmountsByCompany))             //company - company - month 1
	mux.Get("/reports/company/month/year/company", standardMiddleware.ThenFunc(app.transactionHandler.GetMonthlyAmountsByYearAndCompany)) //year and company - company - month

	// users month
	mux.Get("/reports/users/month/global", standardMiddleware.ThenFunc(app.transactionHandler.GetMonthlyAmountsGroupedByYear))                      //global - users - month
	mux.Get("/reports/users/month/user", standardMiddleware.ThenFunc(app.transactionHandler.GetMonthlyAmountsGroupedByYearForUser))                 //user - users - month
	mux.Get("/reports/users/month/user/year", standardMiddleware.ThenFunc(app.transactionHandler.GetMonthlyAmountsForUserByYear))                   //user and year - users - month
	mux.Get("/reports/users/month/user/year/company", standardMiddleware.ThenFunc(app.transactionHandler.GetMonthlyAmountsForUserByYearAndCompany)) //user and year and company - users - month

	// company by company
	mux.Get("/reports/companies/company/global", standardMiddleware.ThenFunc(app.transactionHandler.GetTotalAmountGroupedByCompany))             //global - company - company
	mux.Get("/reports/companies/company/year", standardMiddleware.ThenFunc(app.transactionHandler.GetTotalAmountByCompanyForYear))               //year - company - company
	mux.Get("/reports/companies/company/month", standardMiddleware.ThenFunc(app.transactionHandler.GetTotalAmountByCompanyForMonth))             //month - company - company 2
	mux.Get("/reports/companies/company/year/month", standardMiddleware.ThenFunc(app.transactionHandler.GetTotalAmountByCompanyForYearAndMonth)) //year and month - company - company

	//company users
	mux.Get("/reports/users/company/global", standardMiddleware.ThenFunc(app.transactionHandler.GetTotalAmountGroupedByCompanyForUsers))              //global - users - month
	mux.Get("/reports/users/company/user", standardMiddleware.ThenFunc(app.transactionHandler.GetTotalAmountByCompanyForUser))                        //user - users - month
	mux.Get("/reports/users/company/user/month", standardMiddleware.ThenFunc(app.transactionHandler.GetTotalAmountByCompanyForUserAndMonth))          //user and month - users - month 3
	mux.Get("/reports/users/company/user/year", standardMiddleware.ThenFunc(app.transactionHandler.GetTotalAmountByCompanyForUserAndYear))            //user and year - users - month
	mux.Get("/reports/users/company/user/year/month", standardMiddleware.ThenFunc(app.transactionHandler.GetTotalAmountByCompanyForUserYearAndMonth)) //user and year and company - users - month

	// NOTIFY
	mux.Post("/notify", dynamicMiddleware.ThenFunc(app.fcmHandler.NotifyChange))
	mux.Post("/notify/token/create", dynamicMiddleware.ThenFunc(app.fcmHandler.CreateToken))

	return standardMiddleware.Then(mux)
}
