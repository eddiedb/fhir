package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/intervention-engine/fhir/models"
	"github.com/intervention-engine/fhir/search"
	"gopkg.in/mgo.v2/bson"
)

func DeviceIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case search.UnsupportedError:
				http.Error(rw, x.Error(), http.StatusNotImplemented)
			case search.InvalidSearchError:
				http.Error(rw, x.Error(), http.StatusBadRequest)
			default:
				http.Error(rw, fmt.Sprintf("%s", x), http.StatusInternalServerError)
			}
		}
	}()

	var result []models.Device
	c := Database.C("devices")

	r.ParseForm()
	if len(r.Form) == 0 {
		iter := c.Find(nil).Limit(100).Iter()
		err := iter.All(&result)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else {
		searcher := search.NewMongoSearcher(Database)
		query := search.Query{Resource: "Device", Query: r.URL.RawQuery}
		err := searcher.CreateQuery(query).All(&result)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}

	var deviceEntryList []models.BundleEntryComponent
	for i := range result {
		var entry models.BundleEntryComponent
		entry.Resource = &result[i]
		deviceEntryList = append(deviceEntryList, entry)
	}

	var bundle models.Bundle
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Type = "searchset"
	var total = uint32(len(result))
	bundle.Total = &total
	bundle.Entry = deviceEntryList

	log.Println("Setting device search context")
	context.Set(r, "Device", result)
	context.Set(r, "Resource", "Device")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(&bundle)
}

func LoadDevice(r *http.Request) (*models.Device, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("devices")
	result := models.Device{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting device read context")
	context.Set(r, "Device", result)
	context.Set(r, "Resource", "Device")
	return &result, nil
}

func DeviceShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadDevice(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Device"))
}

func DeviceCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	device := &models.Device{}
	err := decoder.Decode(device)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("devices")
	i := bson.NewObjectId()
	device.Id = i.Hex()
	err = c.Insert(device)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting device create context")
	context.Set(r, "Device", device)
	context.Set(r, "Resource", "Device")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Device/"+i.Hex())
}

func DeviceUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	device := &models.Device{}
	err := decoder.Decode(device)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("devices")
	device.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, device)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting device update context")
	context.Set(r, "Device", device)
	context.Set(r, "Resource", "Device")
	context.Set(r, "Action", "update")
}

func DeviceDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("devices")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting device delete context")
	context.Set(r, "Device", id.Hex())
	context.Set(r, "Resource", "Device")
	context.Set(r, "Action", "delete")
}
