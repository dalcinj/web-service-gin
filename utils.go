package main

import (
	"path/filepath"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func CheckExtention(fileName string) bool {
	extension := filepath.Ext(fileName)
	// acceptedExtentions := []string{".csv", ".xls", ".xlsx"}
	// if contains(acceptedExtentions, extension) == false {
	// 	return fmt.Sprintf("invalid file extention, please upload one of the following: '%v'", acceptedExtentions)
	// }
	if extension == ".csv" {
		return true
	} else {
		return false
	}
}
