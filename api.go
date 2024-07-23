package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w,r); err != nil{
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct{
	listenAddr string
	store Storage
}

func NewApiServer(listenAddr string, store Storage) *APIServer{
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

func (s *APIServer) Run(){
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountById))

	log.Printf("JSON API RUNNING ON PORT : %s", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r* http.Request) error{
	switch (r.Method){
	case "GET":
		return s.handleGetAccount(w,r)
	case "POST":
		return s.handleCreateAccount(w,r)
	case "DELETE":
		return s.handleDeleteAccount(w,r)
	}

	return fmt.Errorf("METHOD NOT ALLOWED %s",r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r* http.Request) error{
	accounts, err := s.store.GetAccounts()
	if err != nil{
		return err
	}

	return WriteJSON(w,http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r* http.Request) error{
	id := mux.Vars(r)["id"]

	fmt.Println(id)
	
	return WriteJSON(w,http.StatusOK, &Account{})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r* http.Request) error{
	accRequest := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(accRequest); err != nil{
		return err
	}

	account := NewAccount(accRequest.FirstName, accRequest.LastName)
	if err := s.store.CreateAccount(account); err!= nil{
		return err
	}
	return WriteJSON(w, http.StatusCreated, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r* http.Request) error{
	
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r* http.Request) error{
	
	return nil
}