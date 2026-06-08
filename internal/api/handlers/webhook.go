package handlers

import (
	"fmt"
	"net/http"
	"strings"
)


func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	defer w.WriteHeader(http.StatusOK)

	if err:= r.ParseForm(); err != nil {
		fmt.Println("failed to parse form: ",err)
		return
	}
	sender:=r.FormValue("From")
	messageBody:=strings.TrimSpace(r.FormValue("Body"))

	phoneNumber:=strings.Replace(sender,"whatsapp:","",1)

	fmt.Println("received msg from %s: %s \n",phoneNumber,messageBody)
}