package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//brand list
// select id, name, url_us, url ridc from brand
type Brand struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Url_us   string `json:"url_us"`
	Url_ridc string `json:"url_ridc"`
}

//model list from a given brand
// select id, name, brand, url_us from model where brand = %d
// INSERT INTO model (name, brand, url_us) VALUES ('soul', '1', 'https://www.ultimatespecs.com/car-specs/Kia-models');
type Model struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Brand  int    `json:"brand"`
	Url_us string `json:"url_us"`
}

//ridc gen from a given model
// select id, full_name, name, url, model_id, year, b_l_f, b_l, b_w from GenRidc where model_id = %d
// INSERT INTO `gen_ridc` (`full_name`, `name`, `url`, `model_id`, `year`, `b_l_f`, `b_l`, `b_w`) VALUES ('soul 2007 sw', 'soul sw', 'https://www.ridc.org.uk/features-reviews/out-and-about/choosing-car/car/kia-e-soul-5dr-saloon-2019', '1', '2007', '12', '123', '1234');
type GenRidc struct {
	Id        int    `json:"id"`
	Full_name string `json:"full_name"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	Model_id  int    `json:"model_id"`
	Year      int    `json:"year"`
	B_l_f     int    `json:"b_l_f"`
	B_l       int    `json:"b_l"`
	B_w       int    `json:"b_w"`
}

type Configuration struct {
	MysqlHost string `json:"mysql-host"`
	MysqlPort int    `json:"mysql-port"`
	MysqlUser string `json:"mysql-username"`
	MysqlPass string `json:"mysql-password"`
	MysqlData string `json:"mysql-database"`
}

func Conf() (conf Configuration) {
	file, _ := os.Open("config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&conf)

	if err != nil {
		panic(err.Error())
	}
	return
}

func dbConnect() (db *sql.DB) {
	config := Conf()

	var openString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.MysqlUser, config.MysqlPass, config.MysqlHost, config.MysqlPort, config.MysqlData)
	db, err := sql.Open("mysql", openString)

	if err != nil {
		panic(err.Error())
	}

	return
}

func executeQuery(query string) {
	//call insert query
	db := dbConnect()
	defer db.Close()

	// query

	result, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()
}

func getBrands(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := "select id, name, url_us, url_ridc from brand"
	if id != "" {
		query = fmt.Sprintf("%s where id = %s", query, id)
	}

	db := dbConnect()
	defer db.Close()

	// query

	result, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	var brands []Brand

	for result.Next() {
		var brand Brand
		err = result.Scan(&brand.Id, &brand.Name, &brand.Url_us, &brand.Url_ridc)
		if err != nil {
			panic(err.Error())
		}
		brands = append(brands, brand)
	}
	obj, err := json.Marshal(brands)

	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, string(obj))
	defer result.Close()

}

func getModels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := "select id, name, brand, url_us from model"
	if id != "" {
		query = fmt.Sprintf("%s where brand = %s", query, id)
	}

	db := dbConnect()
	defer db.Close()

	// query

	result, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	var models []Model

	for result.Next() {
		var model Model
		err = result.Scan(&model.Id, &model.Name, &model.Brand, &model.Url_us)
		if err != nil {
			panic(err.Error())
		}
		models = append(models, model)
	}
	obj, err := json.Marshal(models)

	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, string(obj))
	defer result.Close()

}

func getGenRidc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := "select id, full_name, name, url, model_id, year, b_l_f, b_l, b_w from gen_ridc"
	if id != "" {
		query = fmt.Sprintf("%s where model_id = %s", query, id)
	}

	db := dbConnect()
	defer db.Close()

	// query

	result, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	var models []GenRidc

	for result.Next() {
		var model GenRidc
		err = result.Scan(&model.Id, &model.Full_name, &model.Name, &model.Url, &model.Model_id, &model.Year, &model.B_l_f, &model.B_l, &model.B_w)
		if err != nil {
			panic(err.Error())
		}
		models = append(models, model)
	}
	obj, err := json.Marshal(models)

	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, string(obj))
	defer result.Close()

}

func createBrand(w http.ResponseWriter, r *http.Request) {
	var newBrand Brand
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(reqBody, &newBrand)

	query := fmt.Sprintf("insert into brand (name, url_us, url_ridc) VALUES('%s', '%s', '%s')", newBrand.Name, newBrand.Url_us, newBrand.Url_ridc)
	executeQuery(query)

	obj, err := json.Marshal(newBrand)

	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, string(obj))
}

func createModel(w http.ResponseWriter, r *http.Request) {
	var newModel Model
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(reqBody, &newModel)

	query := fmt.Sprintf("insert into model (name, brand, url_us) VALUES('%s', %d, '%s')", newModel.Name, newModel.Brand, newModel.Url_us)
	executeQuery(query)

	obj, err := json.Marshal(newModel)

	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, string(obj))
}

func createGenRidc(w http.ResponseWriter, r *http.Request) {
	var newModel GenRidc
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(reqBody, &newModel)

	query := fmt.Sprintf("insert into `gen_ridc` (`full_name`, `name`, `url`, `model_id`, `year`, `b_l_f`, `b_l`, `b_w`)  VALUES('%s','%s','%s', %d, %d, %d, %d, %d)", newModel.Full_name, newModel.Name, newModel.Url, newModel.Model_id, newModel.Year, newModel.B_l_f, newModel.B_l, newModel.B_w)
	executeQuery(query)

	obj, err := json.Marshal(newModel)

	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, string(obj))
}

func main() {

	fmt.Println("*** Go MySql ***")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/brands", createBrand).Methods("POST")
	router.HandleFunc("/brands", getBrands)
	router.HandleFunc("/brands/{id}", getBrands)

	router.HandleFunc("/models", createModel).Methods("POST")
	router.HandleFunc("/models", getModels)
	router.HandleFunc("/models/{id}", getModels)

	router.HandleFunc("/genridc", createGenRidc).Methods("POST")
	router.HandleFunc("/genridc", getGenRidc)
	router.HandleFunc("/genridc/{id}", getGenRidc)

	log.Fatal(http.ListenAndServe(":3000", router))

}
