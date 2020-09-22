package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nikvkov/database-stub-api/models"
	"github.com/nikvkov/database-stub-api/swagger"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(w, "Database stub api")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SupplementaryDataHandler(w http.ResponseWriter, r *http.Request) {
	templateSchema, err := getSchema(w, r)
	if err != nil {
		log.Println(err.Error())
		return
	}
	supplementary := templateSchema.SupplementaryData
	result, _ := json.Marshal(supplementary)
	log.Println(string(result))
	w.WriteHeader(http.StatusOK)
	w.Write(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func DisclosureIdentificationHandler(w http.ResponseWriter, r *http.Request) {
	templateSchema, err := getSchema(w, r)
	if err != nil {
		log.Println(err.Error())
		return
	}
	accounts := templateSchema.DisclosureInformation
	result, _ := json.Marshal(accounts)
	log.Println(string(result))
	w.WriteHeader(http.StatusOK)
	w.Write(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getSchema(w http.ResponseWriter, r *http.Request) (*swagger.DisclosureResponse, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code":500,"message":"internal server error"}`))
		return nil, errors.New(`{"code":500,"message":"internal server error"}`)
	}
	log.Println(string(body))
	identifier := models.Identifier{}
	json.Unmarshal(body, &identifier)
	log.Println("Identifier :", identifier.Lei)
	if identifier.Lei == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code":400,"message":"identifier is empty"}`))
		return nil, errors.New(`{"code":400,"message":"identifier is empty"}`)
	}

	templateSchema := readJsonFile("json/schema.json")
	return &templateSchema, nil
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// some identical actions
		log.Println(r.URL.Path)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		next(w, r)
	}
}

func readJsonFile(path string) swagger.DisclosureResponse {
	var schema swagger.DisclosureResponse
	data, err := ioutil.ReadFile(path)
	if err != nil {
		schema = swagger.DisclosureResponse{}
	}
	json.Unmarshal(data, &schema)

	return schema
}
