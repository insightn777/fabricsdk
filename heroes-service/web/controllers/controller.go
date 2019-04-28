package controllers

import (
//      "fmt"
        "github.com/chainHero/heroes-service/blockchain"
//      "html/template"
//      "net/http"
//      "os"
//      "path/filepath"
)

type Application struct {
        Fabric *blockchain.FabricSetup
}

type user struct {
        ObjectType string `json:"DocType"`
        ID string `json:"ID"`
        Name string `json:"Name"`
        Password string `json:"Password"`
        SocialNumber string `json:"SocialNumber"`
        Location string `json:"Location"`
        VoteResult []votevoteresult `json:"VoteResult"`
} // user information

type uservoteresult struct {
        ID string `json:"ID"`
        Location string `json:"Location"`
        Result []int `json:"Result"`
} // user's vote result

type votevoteresult struct {
        Votename string `json:"Votename"`
        Result []int `json:"Result"`
}

type vote struct {
        ObjectType string `json:"DocType"`
	Location string `json:"Location"`
        Votename string `json:"Votename"`
        StartDate string `json:"StartDate"`
        EndDate string `json:"EndDate"`
        Question []string `json:"Question"`
        UserResult []uservoteresult `json:"UserResult"`
} //have question and its result per user

