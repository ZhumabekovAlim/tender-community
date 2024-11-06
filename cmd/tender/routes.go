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
	mux.Post("/transactions", dynamicMiddleware.ThenFunc(app.transactionHandler.CreateTransaction))                                                 // Create a new transaction
	mux.Get("/transactions", standardMiddleware.ThenFunc(app.transactionHandler.GetAllTransactions))                                                // Get all transactions
	mux.Get("/transactions/:id", standardMiddleware.ThenFunc(app.transactionHandler.GetTransactionByID))                                            // Get transaction by ID
	mux.Get("/transactions/user/:id", standardMiddleware.ThenFunc(app.transactionHandler.GetTransactionsByUser))                                    // Get transaction by user ID
	mux.Get("/transactions/company/:id", standardMiddleware.ThenFunc(app.transactionHandler.GetTransactionsByCompany))                              // Get transaction by company ID
	mux.Get("/transactions/user/:user_id/company/:company_id", standardMiddleware.ThenFunc(app.transactionHandler.GetTransactionsForUserByCompany)) // Get transaction by user and company ID
	mux.Get("/transactions/user/zakup/:id", standardMiddleware.ThenFunc(app.transactionHandler.GetTransactionsDebtZakup))                           // Get transaction by user and company ID
	mux.Get("/transactions/user/debt/:id", standardMiddleware.ThenFunc(app.transactionHandler.GetTransactionsDebt))                                 // Get transaction by user and company ID
	mux.Get("/transactions/realization/sum", standardMiddleware.ThenFunc(app.transactionHandler.GetAllTransactionsSum))                             // Get transaction by user and company ID
	mux.Get("/transactions/realization/count/:id", standardMiddleware.ThenFunc(app.transactionHandler.GetTransactionCountsByUserID))                // Get transaction by user and company ID
	mux.Get("/transactions/tranches/debt", standardMiddleware.ThenFunc(app.transactionHandler.GetCompanyDebtById))                                  // Get transaction by user and company ID
	mux.Get("/transactions/tranches/company/debt", standardMiddleware.ThenFunc(app.transactionHandler.GetCompanyDebt))                              // Get transaction by user and company ID
	mux.Get("/transactions/tranches/id/debt/:id", standardMiddleware.ThenFunc(app.transactionHandler.GetCompanyDebtId))                             // Get transaction by user and company ID
	mux.Put("/transactions/:id", standardMiddleware.ThenFunc(app.transactionHandler.UpdateTransaction))                                             // Update transaction by ID
	mux.Del("/transactions/:id", standardMiddleware.ThenFunc(app.transactionHandler.DeleteTransaction))                                             // Delete transaction by ID

	// EXTRA TRANSACTIONS
	mux.Post("/extra_transactions", dynamicMiddleware.ThenFunc(app.extraTransactionHandler.CreateExtraTransaction))                            // Create a new extra transaction
	mux.Get("/extra_transactions", standardMiddleware.ThenFunc(app.extraTransactionHandler.GetAllExtraTransactions))                           // Get all extra transactions
	mux.Get("/extra_transactions/:id", standardMiddleware.ThenFunc(app.extraTransactionHandler.GetExtraTransactionByID))                       // Get extra transaction by ID
	mux.Get("/extra_transactions/user/:id", standardMiddleware.ThenFunc(app.extraTransactionHandler.GetExtraTransactionsByUser))               // Get extra transactions by user ID
	mux.Put("/extra_transactions/:id", standardMiddleware.ThenFunc(app.extraTransactionHandler.UpdateExtraTransaction))                        // Update extra transaction by ID
	mux.Del("/extra_transactions/:id", standardMiddleware.ThenFunc(app.extraTransactionHandler.DeleteExtraTransaction))                        // Delete extra transaction by ID
	mux.Get("/extra_transactions/realization/:id", standardMiddleware.ThenFunc(app.extraTransactionHandler.GetExtraTransactionCountsByUserID)) // Get extra transactions by user ID

	// PERSONAL EXPENSES
	mux.Post("/expenses", dynamicMiddleware.ThenFunc(app.expenseHandler.CreatePersonalExpense))                        // Create a new expense
	mux.Get("/expenses", standardMiddleware.ThenFunc(app.expenseHandler.GetAllPersonalExpenses))                       // Get all expenses
	mux.Get("/expenses/:id", standardMiddleware.ThenFunc(app.expenseHandler.GetPersonalExpenseByID))                   // Get expense by ID
	mux.Get("/expenses/category/:id", standardMiddleware.ThenFunc(app.expenseHandler.GetPersonalExpensesByCategoryId)) // Get expense by ID
	mux.Put("/expenses/:id", standardMiddleware.ThenFunc(app.expenseHandler.UpdatePersonalExpense))                    // Update expense by ID
	mux.Del("/expenses/:id", standardMiddleware.ThenFunc(app.expenseHandler.DeletePersonalExpense))                    // Delete expense by ID

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
	mux.Del("/notify/token/:id", dynamicMiddleware.ThenFunc(app.fcmHandler.DeleteToken))
	mux.Post("/notify/history", dynamicMiddleware.ThenFunc(app.fcmHandler.ShowNotifyHistory))
	mux.Del("/notify/history/:id", dynamicMiddleware.ThenFunc(app.fcmHandler.DeleteNotifyHistory))

	// PASSWORD RECOVERY
	mux.Post("/password/recovery", dynamicMiddleware.ThenFunc(app.userHandler.SendRecoveryHandler))
	mux.Get("/password/recovery/mail", dynamicMiddleware.ThenFunc(app.userHandler.PasswordRecoveryHandler))

	// CATEGORY
	mux.Post("/categories", dynamicMiddleware.ThenFunc(app.categoryHandler.CreateCategory))              // Create a new category
	mux.Get("/categories", standardMiddleware.ThenFunc(app.categoryHandler.GetAllCategories))            // Get all categories
	mux.Get("/categories/parent/:id", standardMiddleware.ThenFunc(app.categoryHandler.GetAllCategories)) // Get all categories
	mux.Get("/categories/:id", standardMiddleware.ThenFunc(app.categoryHandler.GetCategoryByID))         // Get category by ID
	mux.Put("/categories/:id", standardMiddleware.ThenFunc(app.categoryHandler.UpdateCategory))          // Update category by ID
	mux.Del("/categories/:id", standardMiddleware.ThenFunc(app.categoryHandler.DeleteCategory))          // Delete category by ID

	// BALANCE HISTORY
	mux.Post("/balance-history", dynamicMiddleware.ThenFunc(app.balanceHistoryHandler.CreateBalanceHistory))                       // Create a new balance history record
	mux.Get("/balance-history/:id", standardMiddleware.ThenFunc(app.balanceHistoryHandler.GetBalanceHistoryByUserID))              // Get balance history record by user ID
	mux.Get("/balance-history/category/:id", standardMiddleware.ThenFunc(app.balanceHistoryHandler.GetBalanceHistoryByCategoryID)) // Get balance history record by user ID
	mux.Put("/balance-history/:id", standardMiddleware.ThenFunc(app.balanceHistoryHandler.UpdateBalanceHistory))                   // Update balance history record by ID
	mux.Del("/balance-history/:id", standardMiddleware.ThenFunc(app.balanceHistoryHandler.DeleteBalanceHistory))                   // Delete balance history record by ID

	// TENDERS ( GOIK and GOPP)
	mux.Post("/tenders", dynamicMiddleware.ThenFunc(app.tenderHandler.CreateTender))                                  // Create a new tender
	mux.Get("/tenders", standardMiddleware.ThenFunc(app.tenderHandler.GetAllTenders))                                 // Get all tenders
	mux.Get("/tenders/debt/company", standardMiddleware.ThenFunc(app.tenderHandler.GetTotalNetByCompany))             // Get all tenders
	mux.Get("/tenders/:id", standardMiddleware.ThenFunc(app.tenderHandler.GetTenderByID))                             // Get tender by ID
	mux.Get("/tenders/user/:id", standardMiddleware.ThenFunc(app.tenderHandler.GetTendersByUserID))                   // Get tender by user ID
	mux.Get("/tenders/company/:id", standardMiddleware.ThenFunc(app.tenderHandler.GetTendersByCompanyID))             // Get tender by user ID
	mux.Get("/tenders/realization/sum", standardMiddleware.ThenFunc(app.tenderHandler.GetAllTendersSum))              // Get tender by user ID
	mux.Get("/tenders/realization/count/:id", standardMiddleware.ThenFunc(app.tenderHandler.GetTenderCountsByUserID)) // Get tender by user ID
	mux.Put("/tenders/:id", standardMiddleware.ThenFunc(app.tenderHandler.UpdateTender))                              // Update tender by ID
	mux.Del("/tenders/:id", standardMiddleware.ThenFunc(app.tenderHandler.DeleteTender))                              // Delete tender by ID

	// SUMS ALL TABLES
	mux.Get("/sums/all/:id", standardMiddleware.ThenFunc(app.sumHandler.GetSumsByUserID))
	mux.Get("/debts", standardMiddleware.ThenFunc(app.sumHandler.GetDebtsByAccount))

	mux.Get("/sums/:id", standardMiddleware.ThenFunc(app.clientHandler.GetClientData))

	//REALIZATION
	mux.Post("/data/user/:user_id/status/:status", standardMiddleware.ThenFunc(app.transactionHandler.GetAllByUserIDAndStatus))
	mux.Post("/data/transactions/date", standardMiddleware.ThenFunc(app.transactionHandler.GetAllTransactionsByDateRange))
	mux.Post("/data/tenders/date", standardMiddleware.ThenFunc(app.tenderHandler.GetAllTendersByDateRange))
	mux.Post("/data/extra/date", standardMiddleware.ThenFunc(app.extraTransactionHandler.GetAllExtraTransactionsByDateRange))
	mux.Post("/data/transactions/date/company", standardMiddleware.ThenFunc(app.transactionHandler.GetAllTransactionsByDateRangeCompany))
	mux.Post("/data/tenders/date/company", standardMiddleware.ThenFunc(app.tenderHandler.GetAllTendersByDateRangeCompany))
	mux.Post("/data/extra/date/company", standardMiddleware.ThenFunc(app.extraTransactionHandler.GetAllExtraTransactionsByDateRangeCompany))

	// TRANCHE
	mux.Post("/tranches", dynamicMiddleware.ThenFunc(app.trancheHandler.CreateTranche))                                             // Create a new tranche
	mux.Get("/tranches/:id", standardMiddleware.ThenFunc(app.trancheHandler.GetTrancheByID))                                        // Get tranche by ID
	mux.Put("/tranches", standardMiddleware.ThenFunc(app.trancheHandler.UpdateTranche))                                             // Update tranche by ID
	mux.Del("/tranches/:id", standardMiddleware.ThenFunc(app.trancheHandler.DeleteTranche))                                         // Delete tranche by ID
	mux.Get("/tranches/transaction/:transaction_id", standardMiddleware.ThenFunc(app.trancheHandler.GetAllTranchesByTransactionID)) // Get all tranches by transaction_id

	// CHANGE
	mux.Post("/changes", dynamicMiddleware.ThenFunc(app.changeHandler.CreateChange))                                             // Create a new change
	mux.Get("/changes/:id", standardMiddleware.ThenFunc(app.changeHandler.GetChangeByID))                                        // Get change by ID
	mux.Put("/changes", standardMiddleware.ThenFunc(app.changeHandler.UpdateChange))                                             // Update change by ID
	mux.Del("/changes/:id", standardMiddleware.ThenFunc(app.changeHandler.DeleteChange))                                         // Delete change by ID
	mux.Get("/changes/transaction/:transaction_id", standardMiddleware.ThenFunc(app.changeHandler.GetAllChangesByTransactionID)) // Get all changes by transaction_id

	// BALANCE CATEGORY
	mux.Post("/balance_categories", dynamicMiddleware.ThenFunc(app.balanceCategoryHandler.CreateBalanceCategory))      // Create a new balance category
	mux.Get("/balance_categories/:id", standardMiddleware.ThenFunc(app.balanceCategoryHandler.GetBalanceCategoryByID)) // Get balance category by ID
	mux.Put("/balance_categories", standardMiddleware.ThenFunc(app.balanceCategoryHandler.UpdateBalanceCategory))      // Update balance category by ID
	mux.Del("/balance_categories/:id", standardMiddleware.ThenFunc(app.balanceCategoryHandler.DeleteBalanceCategory))  // Delete balance category by ID
	mux.Get("/balance_categories", standardMiddleware.ThenFunc(app.balanceCategoryHandler.GetAllBalanceCategories))    // Get all balance categories

	// PERSONAL DEBTS
	mux.Post("/personal_debts", dynamicMiddleware.ThenFunc(app.personalDebtHandler.CreatePersonalDebt))      // Create a new personal debt
	mux.Get("/personal_debts/:id", standardMiddleware.ThenFunc(app.personalDebtHandler.GetPersonalDebtByID)) // Get personal debt by ID
	mux.Put("/personal_debts", standardMiddleware.ThenFunc(app.personalDebtHandler.UpdatePersonalDebt))      // Update personal debt by ID
	mux.Del("/personal_debts/:id", standardMiddleware.ThenFunc(app.personalDebtHandler.DeletePersonalDebt))  // Delete personal debt by ID
	mux.Get("/personal_debts", standardMiddleware.ThenFunc(app.personalDebtHandler.GetAllPersonalDebts))     // Get all personal debts

	return standardMiddleware.Then(mux)
}
