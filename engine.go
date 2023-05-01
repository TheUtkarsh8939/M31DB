package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func dbengine(args []string) string {
	var returnval string = ""
	switch args[1] {
	case "select":
		args[2] = strings.ReplaceAll(args[2], "/", "\\")
		file := strings.Replace(args[0], "\\main.exe", "\\db\\", -1)
		file = file + args[2] + ".json"

		dataCh := make(chan []byte)
		errCh := make(chan error)
		go func() {
			datab, err := os.ReadFile(file)
			ferr(err)
			mdata := formatJSON(string(datab))
			err = os.WriteFile(file, []byte(mdata), 0644)
			ferr(err)
		}()
		go func() {
			data, err := os.ReadFile(file)
			if err != nil {
				errCh <- err
				return
			}
			dataCh <- data
		}()

		select {
		case data := <-dataCh:
			returnval = string(data)
		case err := <-errCh:
			fmt.Println("Error Fetching Data:", err)
		}
		runtime.GC()
	case "insert":
		args[2] = strings.ReplaceAll(args[2], "/", "\\")
		file := strings.Replace(args[0], "\\main.exe", "\\db\\", -1)
		file = file + args[2] + ".json"

		dataCh := make(chan []byte)
		go func() {
			datab, err := os.ReadFile(file)
			ferr(err)
			mdata := formatJSON(string(datab))
			err = os.WriteFile(file, []byte(mdata), 0644)
			ferr(err)
		}()
		errCh := make(chan error)
		go func() {
			data, err := os.ReadFile(file)
			if err != nil {
				errCh <- err
				return
			}
			argsmod := convandrotojson(args[3])
			mddata := string(data)
			modifieddata := strings.ReplaceAll(mddata, "]", "") + ",\n" + argsmod + "]"
			dataCh <- []byte(modifieddata)
		}()

		go func() {
			err := os.WriteFile(file, <-dataCh, 00644)
			if err != nil {
				if strings.Contains(err.Error(), "file") {
					errCh <- errors.New("database does not exists")
				} else {
					errCh <- err
				}
				return
			}
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			if err == nil {
				returnval = "Done Successfuly"
			} else {
				fmt.Println("Error Inserting Data", err)
			}
		}
		runtime.GC()
	case "create":
		if args[2] != "table" {
			err := os.Mkdir("db\\"+args[2], 0755)
			if err == nil {
				returnval = "Done Successfuly"
			} else {
				fmt.Println("Database Already exists")
			}
		} else if args[2] == "table" {
			filename := strings.Replace(args[3], "/", "\\", -1)
			file, err := os.Create(filename)
			println(file, "\n")
			ferr(err)
			err = os.WriteFile(args[0]+"\\db\\"+args[3], []byte(formatJSON("[{\"name\":\""+args[3]+"\"}]")), 0666)
			ferr(err)
		} else {
			fmt.Println("\033[31mError:Can't recognize what you are trying to create!")
		}

	}
	return returnval
}