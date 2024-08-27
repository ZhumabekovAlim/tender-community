package main

import (
	"context"
	"database/sql"
	"firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"tender/internal/handlers"
	"tender/internal/repositories"
	"tender/internal/services"
)

type application struct {
	errorLog                *log.Logger
	infoLog                 *log.Logger
	userHandler             *handlers.UserHandler
	permissionHandler       *handlers.PermissionHandler
	companyHandler          *handlers.CompanyHandler
	transactionHandler      *handlers.TransactionHandler
	expenseHandler          *handlers.PersonalExpenseHandler
	extraTransactionHandler *handlers.ExtraTransactionHandler
	fcmHandler              *handlers.FCMHandler
}

func initializeApp(db *sql.DB, errorLog, infoLog *log.Logger) *application {

	ctx := context.Background()
	sa := option.WithCredentialsFile("/root/go/src/tender/cmd/tender/serviceAccountKey.json")

	firebaseApp, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: "tendercommunity-17cd5"}, sa)
	if err != nil {
		errorLog.Fatalf("Ошибка в нахождении приложения: %v\n", err)
	}

	fcmClient, err := firebaseApp.Messaging(ctx)
	if err != nil {
		errorLog.Fatalf("Ошибка при неверном ID устройства: %v\n", err)
	}

	fcmHandler := handlers.NewFCMHandler(fcmClient, db)

	userRepo := &repositories.UserRepository{Db: db}
	userService := &services.UserService{Repo: userRepo}
	userHandler := &handlers.UserHandler{Service: userService}

	permissionRepo := &repositories.PermissionRepository{Db: db}
	permissionService := &services.PermissionService{Repo: permissionRepo}
	permissionHandler := &handlers.PermissionHandler{Service: permissionService}

	companyRepo := &repositories.CompanyRepository{Db: db}
	companyService := &services.CompanyService{Repo: companyRepo}
	companyHandler := &handlers.CompanyHandler{Service: companyService}

	expenseRepo := &repositories.PersonalExpenseRepository{Db: db}
	expenseService := &services.PersonalExpenseService{Repo: expenseRepo}
	expenseHandler := &handlers.PersonalExpenseHandler{Service: expenseService}

	extraTransactionRepo := &repositories.ExtraTransactionRepository{Db: db}
	extraTransactionService := &services.ExtraTransactionService{Repo: extraTransactionRepo}
	extraTransactionHandler := &handlers.ExtraTransactionHandler{Service: extraTransactionService}

	transactionRepo := &repositories.TransactionRepository{Db: db}
	transactionService := &services.TransactionService{Repo: transactionRepo}
	transactionHandler := &handlers.TransactionHandler{
		Service:                 transactionService,
		ExtraTransactionService: extraTransactionService,
	}

	return &application{
		errorLog:                errorLog,
		infoLog:                 infoLog,
		userHandler:             userHandler,
		permissionHandler:       permissionHandler,
		companyHandler:          companyHandler,
		transactionHandler:      transactionHandler,
		expenseHandler:          expenseHandler,
		extraTransactionHandler: extraTransactionHandler,
		fcmHandler:              fcmHandler,
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("%v", err)
		panic("failed to connect to database")
		return nil, err
	}
	db.SetMaxIdleConns(35)
	if err = db.Ping(); err != nil {
		log.Printf("%v", err)
		panic("failed to ping the database")
		return nil, err
	}
	fmt.Println("successfully connected")

	return db, nil
}

func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}
