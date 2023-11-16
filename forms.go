package shuffle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func HandleCreateForms(resp http.ResponseWriter, request *http.Request) {
	cors := HandleCors(resp, request)
	if cors {
		return
	}

	user, err := HandleApiAuthentication(resp, request)
	if err != nil {
		log.Printf("[AUDIT] INITIAL Api authentication failed in form creation: %s", err)
		resp.WriteHeader(401)
		resp.Write([]byte(`{"success": false}`))
		return
	}

	if user.Role != "admin" {
		log.Printf("[AUTH] User isn't admin")
		resp.WriteHeader(401)
		resp.Write([]byte(`{"success": false, "reason": "Need to be admin to create forms"}`))
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println("Failed reading body")
		resp.WriteHeader(401)
		resp.Write([]byte(`{"success": false, "reason": "Failed to read data"}`))
		return
	}

	log.Printf("[DEBUG] Body: %s", body)

	var curform FormStructure
	err = json.Unmarshal(body, &curform)
	if err != nil {
		log.Printf("[ERROR] Failed unmarshaling: %s", err)
		resp.WriteHeader(401)
		resp.Write([]byte(`{"success": false, "reason": "Failed to unmarshal data"}`))
		return
	}

	if len(curform.Name) < 1 {
		curform.Name = "Untitled form"
	}

	ctx := GetContext(request)

	log.Printf("[DEBUG] Body: %s", curform)

	// generate a new form id
	curform.Id = uuid.NewV4().String()

	// add the form to the database
	SetForm(ctx, curform, curform.Id)

	resp.WriteHeader(200)
	resp.Write([]byte(fmt.Sprintf(`{"success": true, "id": "%s"}`, curform.Id)))
}

func HandleGetForms(resp http.ResponseWriter, request *http.Request) {
	cors := HandleCors(resp, request)
	if cors {
		return
	}

	_, err := HandleApiAuthentication(resp, request)
	if err != nil {
		log.Printf("[AUDIT] INITIAL Api authentication failed in form creation: %s", err)
		resp.WriteHeader(401)
		resp.Write([]byte(`{"success": false}`))
		return
	}

	ctx := GetContext(request)

	forms, err := GetForms(ctx)
	if err != nil {
		log.Printf("[ERROR] Failed getting forms: %s", err)
		resp.WriteHeader(401)
		resp.Write([]byte(`{"success": false, "reason": "Failed to get forms"}`))
		return
	}

	formsJson, err := json.Marshal(forms)
	if err != nil {
		log.Printf("[ERROR] Failed marshaling forms: %s", err)
		resp.WriteHeader(401)
		resp.Write([]byte(`{"success": false, "reason": "Failed to marshal forms"}`))
		return
	}

	resp.WriteHeader(200)
	resp.Write([]byte(fmt.Sprintf(`{"success": true, "forms": %s}`, formsJson)))
}

func HandleGetForm(resp http.ResponseWriter, request *http.Request) {
	cors := HandleCors(resp, request)
	if cors {
		return
	}

	// _, err := HandleApiAuthentication(resp, request)
	// if err != nil {
	// 	log.Printf("[AUDIT] INITIAL Api authentication failed in form get: %s", err)
	// 	resp.WriteHeader(401)
	// 	resp.Write([]byte(`{"success": false}`))
	// 	return
	// }

	var formId string
	location := strings.Split(request.URL.String(), "/")
	if location[1] == "api" {
		if len(location) <= 4 {
			log.Printf("[INFO] Path too short: %d", len(location))
			resp.WriteHeader(401)
			resp.Write([]byte(`{"success": false}`))
			return
		}

		formId = location[4]
	}

	log.Printf("[DEBUG] Form ID: %s", formId)

	if len(formId) != 36 {
		resp.WriteHeader(400)
		resp.Write([]byte(`{"success": false, "message": "ID not valid"}`))
		return
	}

	ctx := GetContext(request)
	form, err := GetForm(ctx, formId)
	// print form
	log.Printf("[DEBUG] Form: %s", form)
	if err != nil {
		log.Printf("[ERROR] Failed getting form: %s", err)
		resp.WriteHeader(400)
		resp.Write([]byte(`{"success": false, "reason": "Failed to get form"}`))
		return
	}

	formJson, err := json.Marshal(form)
	if err != nil {
		log.Printf("[ERROR] Failed marshaling form: %s", err)
		resp.WriteHeader(400)
		resp.Write([]byte(`{"success": false, "reason": "Failed to marshal form"}`))
		return
	}

	resp.WriteHeader(200)
	resp.Write([]byte(formJson))
	return
	// resp.Write([]byte(fmt.Sprintf(`{"success": true, "form": %s}`, formJson)))
}
