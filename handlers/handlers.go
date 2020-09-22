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

const jsonDisclosureInfo = `{"safekeepingAccountAndHolding":[{"accountServicer":{"anyBIC":"DIRUNGG1242"},"shareholdingBalanceOnOwnAccount":{"unit":1.0},"shareholdingBalanceOnClientAccount":{"unit":1.0},"totalShareholdingBalance":{"unit":1.0},"accountSubLevel":{"belowThresholdShareholdingQuantity":{"unit":1.0},"disclosure":[{"safekeepingAccount":"21216635","accountHolder":{},"shareholdingBalance":[{"shareholdingType":"BENE","quantity":{"unit":1.0}}]},{"safekeepingAccount":"16279221","accountHolder":{},"shareholdingBalance":[{"shareholdingType":"BENE","quantity":{"unit":1.0}}]},{"safekeepingAccount":"75080116","accountHolder":{},"shareholdingBalance":[{"shareholdingType":"BENE","quantity":{"unit":1.0}}]},{"safekeepingAccount":"07912963","accountHolder":{"legalPerson":{"nameAndAddress":{"name":"Elda","address":{"addressType":"ADDR","addressLine":["2467 Reilly Stream","85984 Vandervort Estates"],"streetName":"Giovanny Camp","buildingNumber":"HNO:127","postBox":"P.O. Box 60","postCode":"12375","townName":"Port Howellshire","countrySubDivision":"Bedfordshire","country":{}}},"emailAddress":{},"identification":{"anyBIC":"NLATEVARXXX"},"countryOfIncorporation":{},"activityIndicator":"I4756","investorType":{},"ownership":{"ownershipType":{"code":"USUF"}}}},"shareholdingBalance":[{"shareholdingType":"NOMI","forwardRequestDetails":{},"quantity":{"unit":1.0}}]},{"safekeepingAccount":"45327269","accountHolder":{"legalPerson":{"nameAndAddress":{"name":"Willard","address":{"addressType":"ADDR","addressLine":["77272 Vladimir Inlet","712 Kuhn Skyway"],"streetName":"Thompson Drives","buildingNumber":"HNO:127","postBox":"P.O. Box 60","postCode":"12375","townName":"East Ali","countrySubDivision":"Buckinghamshire","country":{}}},"emailAddress":{},"identification":{"anyBIC":"FIDTNGLAXXX"},"countryOfIncorporation":{},"activityIndicator":"R7196","investorType":{},"ownership":{"ownershipType":{"code":"USUF"}}}},"shareholdingBalance":[{"shareholdingType":"NOMI","forwardRequestDetails":{},"quantity":{"unit":1.0}}]},{"safekeepingAccount":"08752411","accountHolder":{"legalPerson":{"nameAndAddress":{"name":"Pink","address":{"addressType":"ADDR","addressLine":["411 Monserrat River","564 Marie Run"],"streetName":"Cormier Isle","buildingNumber":"HNO:127","postBox":"P.O. Box 60","postCode":"12375","townName":"West Ceceliachester","countrySubDivision":"Berkshire","country":{}}},"emailAddress":{},"identification":{"anyBIC":"SLAMCHZZXXX"},"countryOfIncorporation":{},"activityIndicator":"U0765","investorType":{},"ownership":{"ownershipType":{"code":"USUF"}}}},"shareholdingBalance":[{"shareholdingType":"NOMI","forwardRequestDetails":{},"quantity":{"unit":1.0}}]},{"safekeepingAccount":"42849829","accountHolder":{"legalPerson":{"nameAndAddress":{"name":"Ines","address":{"addressType":"ADDR","addressLine":["516 Ziemann Trail","1126 Dibbert Valleys"],"streetName":"Michel Island","buildingNumber":"HNO:127","postBox":"P.O. Box 60","postCode":"12375","townName":"Kodyton","countrySubDivision":"Bedfordshire","country":{}}},"emailAddress":{},"identification":{"anyBIC":"SLZMAMNZXXX"},"countryOfIncorporation":{},"activityIndicator":"I0274","investorType":{},"ownership":{"ownershipType":{"code":"USUF"}}}},"shareholdingBalance":[{"shareholdingType":"NOMI","forwardRequestDetails":{},"quantity":{"unit":1.0}}]}]}}]}`

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
	_, err := getSchema(w, r)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//accounts := templateSchema.DisclosureInformation
	result := []byte(jsonDisclosureInfo) //json.Marshal(accounts)
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
