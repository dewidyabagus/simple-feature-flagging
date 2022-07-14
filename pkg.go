package main

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func getEnv[T typ](key string, def T) T {
	var res interface{} = def

	switch reflect.TypeOf(def).Kind() {
	case reflect.String:
		if !isEmpty(os.Getenv(key)) {
			res = os.Getenv(key)
		}

	case reflect.Uint32:
		if !isEmpty(os.Getenv(key)) {
			rs, err := strconv.Atoi(os.Getenv(key))
			if err != nil {
				log.Println("key env:", key, "error convert value, to set default value")
			} else {
				res = uint32(rs)
			}
		}

	}

	return res.(T)
}

func isEmpty(val string) bool {
	return strings.TrimSpace(val) == ""
}
