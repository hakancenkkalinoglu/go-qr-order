package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong! Gorder sunucusu ayakta.")
	})

	fmt.Println("🚀 Sunucu 8080 portunda başlatılıyor...")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Sunucu başlatılamadı:", err)
	}
}
