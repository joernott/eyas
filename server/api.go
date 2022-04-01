package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joernott/eyas/eyaml"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type ErrorMessage struct {
	Message string `json:"message"`
	Err     error  `json:"error"`
}

type ErrorResponse struct {
	Error *ErrorMessage `json:"error"`
}

type KeyResponse struct {
	Error *ErrorMessage `json:"error"`
	Data  []string      `json:"data"`
}

type EncryptResponse struct {
	Error *ErrorMessage     `json:"error"`
	Data  map[string]string `json:"data"`
}

func errorOut(w http.ResponseWriter, code int, id string, message string, err error) {
	var response ErrorResponse
	response.Error = &ErrorMessage{message, err}
	w.WriteHeader(code)
	log.Error().Err(err).Str("id", id).Msg(message)
	json, _ := json.Marshal(response)
	fmt.Fprintf(w, string(json[:]))
}

func PKCS7Keys(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response KeyResponse

	list, err := os.ReadDir(KeyPath)
	if err != nil {
		errorOut(w, 500, "ERR00010", "Could not get key list", err)
		return
	}
	a := zerolog.Arr()
	for _, d := range list {
		if d.IsDir() {
			if _, err := os.Stat(KeyPath + "/" + d.Name() + "/public_key.pkcs7.pem"); errors.Is(err, os.ErrNotExist) {
				log.Warn().Str("id", "WRN00010").Str("path", KeyPath+"/"+d.Name()).Msg("Ignoring " + KeyPath + "/" + d.Name() + " as it does not contain a readable public_key.pkcs7.pem")
			} else {
				response.Data = append(response.Data, d.Name())
				a.Str(d.Name())
			}
		}
	}
	log.Debug().Array("keys", a).Str("id", "DBG00001").Msg("Get PKCS7 keys")
	json, err := json.Marshal(response)
	if err != nil {
		errorOut(w, 500, "ERR00011", "Could not transform key list to json", err)
		return
	}
	fmt.Fprintf(w, string(json[:]))
}

func EncryptSingle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response EncryptResponse
	var output eyaml.OutputType
	response.Data = make(map[string]string)

	r.ParseForm()
	key := r.Form.Get("key")
	if key == "" {
		errorOut(w, 422, "ERR00020", "No Key provided", nil)
		return
	}
	password := r.Form.Get("password1")
	if password == "" {
		errorOut(w, 422, "ERR00021", "No Password provided", nil)
		return
	}
	o := r.Form.Get("output")
	switch o {
	case "block":
		{
			output = eyaml.Block
		}
	case "string":
		{
			output = eyaml.String
		}
	default:
		{
			errorOut(w, 422, "ERR00022", "Illegal output format "+o, nil)
			return
		}
	}
	use := r.Form["use"]
	for _, pkcs7 := range use {
		encrypted, err := eyaml.Encrypt(password, key, output, KeyPath+"/"+pkcs7+"/public_key.pkcs7.pem")
		if err != nil {
			errorOut(w, 500, "ERR00023", "Could not encrypt with key "+pkcs7, err)
			return
		}
		log.Debug().Str("id", "DBG00020").Str("pkcs7", pkcs7).Str("key", key).Str("password", password).Str("encrypted", encrypted).Str("type", "single").Msg("Encrypted")
		response.Data[pkcs7] = encrypted
	}
	json, err := json.Marshal(response)
	if err != nil {
		errorOut(w, 500, "ERR00024", "Could not transform encrypted keys to json", err)
		return
	}
	fmt.Fprintf(w, string(json[:]))
}

func EncryptYaml(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response EncryptResponse
	var output eyaml.OutputType
	response.Data = make(map[string]string)

	r.ParseForm()
	yamldata := r.Form.Get("yaml")
	if yamldata == "" {
		errorOut(w, 422, "ERR00031", "No Yaml data provided", nil)
		return
	}
	o := r.Form.Get("output")
	switch o {
	case "block":
		{
			output = eyaml.Block
		}
	case "string":
		{
			output = eyaml.String
		}
	default:
		{
			errorOut(w, 422, "ERR00032", "Illegal output format "+o, nil)
			return
		}
	}
	use := r.Form["use"]

	data := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(yamldata), &data)
	if err != nil {
		errorOut(w, 422, "ERR00033", "Could not parse yaml data", err)
		return
	}
	for _, pkcs7 := range use {
		for key, value := range data {
			switch password := value.(type) {
			case string:
				{
					encrypted, err := eyaml.Encrypt(string(password), key, output, "keys/"+pkcs7+"/public_key.pkcs7.pem")
					if err != nil {
						errorOut(w, 500, "ERR00034", "Could not encrypt with key "+pkcs7, err)
						return
					}
					log.Debug().Str("id", "DBG00030").Str("pkcs7", pkcs7).Str("key", key).Str("password", string(password)).Str("encrypted", encrypted).Str("type", "yaml").Msg("Encrypted")
					response.Data[pkcs7] = response.Data[pkcs7] + " \n " + encrypted
				}
			default:
				{
					log.Warn().Str("key", key).Msg("Skipping non-string")
				}
			}
		}
	}
	json, err := json.Marshal(response)
	if err != nil {
		errorOut(w, 500, "ERR00035", "Could not transform encrypted keys to json", err)
		return
	}
	fmt.Fprintf(w, string(json[:]))
}

func EncryptCSV(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response EncryptResponse
	var output eyaml.OutputType
	response.Data = make(map[string]string)

	r.ParseForm()
	v := r.Form.Get("keycol")
	keycol, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		errorOut(w, 422, "ERR00040", "Could not parse key column number "+v, err)
		return
	}
	if keycol < 0 {
		errorOut(w, 422, "ERR00041", "Key column number must be >= 0", nil)
		return
	}
	v = r.Form.Get("passcol")
	passcol, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		errorOut(w, 422, "ERR00042", "Could not parse password column number "+v, err)
		return
	}
	if passcol < 0 {
		errorOut(w, 422, "ERR00043", "Password column number must be >= 0", nil)
		return
	}
	separator := r.Form.Get("separator")
	if separator == "" {
		errorOut(w, 422, "ERR00044", "Separator must not be empty", err)
		return
	}
	o := r.Form.Get("output")
	switch o {
	case "block":
		{
			output = eyaml.Block
		}
	case "string":
		{
			output = eyaml.String
		}
	default:
		{
			errorOut(w, 422, "ERR00045", "Illegal output format "+o, nil)
			return
		}
	}
	lines := strings.Split(strings.ReplaceAll(r.Form.Get("csv"), "\r\n", "\n"), "\n")
	use := r.Form["use"]
	for _, pkcs7 := range use {
		for _, line := range lines {
			fields := strings.Split(line, separator)
			if int(keycol) >= len(fields) {
				errorOut(w, 422, "ERR00046", "Key column out of bounds for line "+line, nil)
				return
			}
			if int(passcol) >= len(fields) {
				errorOut(w, 422, "ERR00047", "Password column out of bounds for line "+line, nil)
				return
			}
			encrypted, err := eyaml.Encrypt(fields[passcol], fields[keycol], output, "keys/"+pkcs7+"/public_key.pkcs7.pem")
			if err != nil {
				errorOut(w, 500, "ERR00048", "Could not encrypt key "+fields[keycol]+" with "+pkcs7, err)
				return
			}
			log.Debug().Str("id", "DBG00040").Str("pkcs7", pkcs7).Str("key", fields[keycol]).Str("password", fields[passcol]).Str("encrypted", encrypted).Str("type", "csv").Msg("Encrypted")
			response.Data[pkcs7] = response.Data[pkcs7] + "\n" + encrypted
		}

	}
	json, err := json.Marshal(response)
	if err != nil {
		errorOut(w, 500, "ERR00049", "Could not transform encrypted keys to json", err)
		return
	}
	fmt.Fprintf(w, string(json[:]))
}
