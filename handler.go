package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aiteung/atapi"
	"github.com/aiteung/atmessage"
	"github.com/aiteung/module/model"
	"github.com/whatsauth/wa"
)

func PostBalasan(w http.ResponseWriter, r *http.Request) {
	var msg model.IteungMessage
	var resp atmessage.Response
	json.NewDecoder(r.Body).Decode(&msg)
	link := "https://medium.com/@imanulhaq5/tugas-chapt-7-whatsauth-free-2fa-otp-notif-whatsapp-gateway-api-gratis-d51198a9e9a7" + " dan " + "https://youtu.be/4aOVpoJaI1M"
	if r.Header.Get("Secret") == os.Getenv("SECRET") {
		if msg.Message == "loc" || msg.Message == "Loc" || msg.Message == "lokasi" || msg.LiveLoc {
			location, err := ReverseGeocode(msg.Latitude, msg.Longitude)
			if err != nil {
				// Handle the error (e.g., log it) and set a default location name
				location = "Unknown Location"
			}

			reply := fmt.Sprintf("Halo kamu, pasti lagi di %s \nKoordinatnya : %s - %s\nBerikut ini cara penggunaan WhatsAuth ada di link berikut ini ya kak..."+
				"yaa kak %s\n", location,
				strconv.Itoa(int(msg.Longitude)), strconv.Itoa(int(msg.Latitude)), link)
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: reply,
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
		} else {
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: "Halo kamu, " + msg.Alias_name + "\nMaulana nya lagi gaadaa \nAku Maulana, Salam kenal ya... \nBerikut ini cara penggunaan WhatsAuth ada di link berikut ini ya kak...\n" + link,
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
		}
	} else {
		resp.Response = "Secret Salah"
	}
	fmt.Fprintf(w, resp.Response)
}

func ReverseGeocode(latitude, longitude float64) (string, error) {
	// OSM Nominatim API endpoint
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", latitude, longitude)

	// Make a GET request to the API
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	// Extract the place name from the response
	displayName, ok := result["display_name"].(string)
	if !ok {
		return "", fmt.Errorf("unable to extract display_name from the API response")
	}

	return displayName, nil
}

func Liveloc(w http.ResponseWriter, r *http.Request) {
	var msg model.IteungMessage
	var resp atmessage.Response
	json.NewDecoder(r.Body).Decode(&msg)

	// Reverse geocode to get the place name
	location, err := ReverseGeocode(msg.Latitude, msg.Longitude)
	if err != nil {
		// Handle the error (e.g., log it) and set a default location name
		location = "Unknown Location"
	}

	reply := fmt.Sprintf("Halo, kamu pasti lagi di... %s \nKoordinatnya : %s - %s\n", location,
		strconv.Itoa(int(msg.Longitude)), strconv.Itoa(int(msg.Latitude)))

	if r.Header.Get("Secret") == os.Getenv("SECRET") {
		dt := &wa.TextMessage{
			To:       msg.Phone_number,
			IsGroup:  false,
			Messages: reply,
		}
		resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://cloud.wa.my.id/api/send/message/text")
	} else {
		resp.Response = "Secret Salah"
	}
	fmt.Fprintf(w, resp.Response)
}
