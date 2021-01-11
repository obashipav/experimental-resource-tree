package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func httpAbortWithMessage(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err := w.Write([]byte(`{"error":"` + message + `"}`))
	if err != nil {
		log.Println("failed to write http response: ", err)
	}
}

func httpResponse(writer http.ResponseWriter, content interface{}, status int) {
	//write the json to a temporary buffer
	buffer := bytes.NewBuffer([]byte(""))
	iter := allowAPI.BorrowStream(buffer)
	defer allowAPI.ReturnStream(iter)

	//write
	iter.WriteVal(content)

	//internal encoding error
	if iter.Error != nil {
		httpAbortWithMessage(writer, "failed to write content buffer: ", http.StatusInternalServerError)
		return
	}

	//flush
	if err := iter.Flush(); err != nil {
		httpAbortWithMessage(writer, "failed to clear iter buffer: ", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, err := io.Copy(writer, buffer)
	if err != nil {
		log.Println("failed to write content back to response: ", err)
	}
}
