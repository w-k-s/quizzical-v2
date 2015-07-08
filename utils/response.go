package utils

import(
	"encoding/xml"
	"encoding/json"
	"net/http"
	"fmt"
)

func RespondWithJSON(w http.ResponseWriter, v interface{}) {

	json, err := json.MarshalIndent(v, "", "    ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(json))

}

func RespondWithXML(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "text/xml")
	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RespondWithText(w http.ResponseWriter,text string){
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, text)
}