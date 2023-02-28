package server

import (
	"encoding/json"
	"example.com/supplyChainManagement/models"
	"example.com/supplyChainManagement/providers/dbhelpers/adminprovider/admin"
	"example.com/supplyChainManagement/scmerrors"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"

	"example.com/supplyChainManagement/utils"
)

type Response struct {
	Persons []Person `json:"persons"`
}

type Person struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (srv *Server) createNewUser(resp http.ResponseWriter, req *http.Request) {
	//var newUserReq models.CreateNewUserRequest
	var data Person
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		scmerrors.RespondClientErr(resp, err, http.StatusBadRequest)
		return
	}
	utils.EncodeJSONBody(resp, http.StatusCreated, data)
}

func (srv *Server) login(resp http.ResponseWriter, req *http.Request) {

}

func (srv *Server) home(resp http.ResponseWriter, req *http.Request) {
	//vars := mux.Vars(req)
	//title := "title"
	//page := "page"
	type Bio struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var bio Bio
	err := json.NewDecoder(req.Body).Decode(&bio)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	bio.Age = 26
	bio.Name = "shiv"

	//fmt.Fprintf(resp, "You've requested the book: %s on page %s\n", title, page)
	//data := fmt.Sprintf("You've requested the book: %s on page %s\n", title, page)
	//resp.Write([]byte("hi there "))
	//fmt.Fprintf(w, string("response added again "))

	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{
		"data": bio,
	})
}

func (srv *Server) HealthCheck(resp http.ResponseWriter, r *http.Request) {
	//specify status code
	resp.WriteHeader(http.StatusOK)

	//update response writer
	fmt.Fprintf(resp, "API is up and running")

	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{
		"data": "health is good !",
	})
}

func (srv *Server) CreateUser(resp http.ResponseWriter, r *http.Request) {
	//ctx := context.Background()
	var newUserJson models.CreateUserParams
	var newUser admin.CreateUserParams

	err := json.NewDecoder(r.Body).Decode(&newUserJson)
	if err != nil {
		logrus.Error("NewFunction : unable to decode request body ", err)
	}
	newUser.Name = newUserJson.Name
	newUser.Age = newUserJson.Age
	newUser.Password = newUserJson.Password
	newUser.Address = newUserJson.Address
	newUser.CountryCode = newUserJson.CountryCode
	newUser.Phone = newUserJson.Phone
	newUser.Email = newUserJson.Email

	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	// create an author
	//insertedAuthor, err := srv.AdminQueries.CreateUser(ctx, newUser)
	//if err != nil {
	//	logrus.Error("error creating user in database", err)
	//}

	utils.EncodeJSONBody(resp, http.StatusCreated, map[string]interface{}{
		"message": "success",
		//"data":    insertedAuthor,
	})
}

func (srv *Server) FetchAllUser(resp http.ResponseWriter, r *http.Request) {
	//ctx := context.Background()
	var newUser admin.CreateUserParams

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		logrus.Error("NewFunction : unable to decode request body ", err)
	}

	//fetchedAuthor, err := srv.AdminQueries.ListUsers(ctx)
	//if err != nil {
	//	logrus.Error("error creating user in database", err)
	//}
	//
	//utils.EncodeJSON200Body(resp, fetchedAuthor)

}
