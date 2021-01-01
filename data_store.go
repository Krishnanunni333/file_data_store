package main

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"strings"
	"sync"
	"time"
)

var mutex sync.Mutex

//Checking and printing all errors
func check(e error) {

	if e != nil {

		panic(e)

	}
}

//Create a key value pair
func create() {

	mutex.Lock()
	in := bufio.NewReader(os.Stdin)
	var p string

	//Reading absolute path from the user
	fmt.Println("Enter the file path with filename (Only CSV files are accepted) or press enter key for default path")
	fmt.Scanf("%s", &p)

	//If user prefers default location
	if len(p) <= 1 {
		//Getting current username
		usr, err := user.Current()
		check(err)
		//Current user home directory
		p = string(usr.HomeDir) + "/data.csv"

	}
	//Splitting the path recieved
	split := strings.Split(p, "/")
	p = path.Join(split...)

	//Checking the size of the specified file
	fi, err := os.Stat("/" + p)
	if err == nil {
		// Get file size
		size := fi.Size()
		if size >= 1073741824 {
			fmt.Println("Sorry! No space in this file")
			fmt.Println()
		}

	}

	//Getting key value
	fmt.Println("Enter the Key")
	key, err := in.ReadString('\n')
	check(err)
	key = strings.TrimSpace(key)
	csvfile, err := os.Open("/" + p)
	if err == nil {
		r := csv.NewReader(csvfile)

		// Iterate through the records
		for {
			// Read each record from csv
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Sorry")
				break
			}

			//If key is already present in the file
			if record[0] == key {

				fmt.Println("Sorry, key Already present")
				fmt.Println()
				return

			}

		}
		//Closing the file
		if err := csvfile.Close(); err != nil {

			fmt.Println("Error in closing the file")
			defer mutex.Unlock()
			return

		}
	}
	//Getting JSON Object and cleaning it to write into the file
	fmt.Println("Enter the JSON")
	decoder := json.NewDecoder(os.Stdin)
	decoder.UseNumber()
	js := make(map[string]interface{})
	decoder.Decode(&js)
	//Converting JSON to string to write it into the file
	jsonString, _ := json.Marshal(js)
	var value string = string(jsonString)
	if strings.ReplaceAll(value, "{}", "") == "" {
		fmt.Println("Invalid JSON format")
		fmt.Println()
		defer mutex.Unlock()
		return
	}
	//Replacing double quotes to 2 double quotes
	value = string(jsonString)
	value = strings.ReplaceAll(value, "\"", "\"\"")

	//Getting the lifetime of the data
	var sec int
	var newT string = "nil"
	t := time.Now()
	fmt.Println()
	fmt.Println("Enter the time to live in seconds")
	fmt.Scanf("%d", &sec)
	if sec > 0 {
		newT = t.Add(time.Second * time.Duration(sec)).String()
	}

	f, err := os.OpenFile("/"+p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {

		fmt.Println("File not found or path not reachable")
		defer mutex.Unlock()
		return

	}

	k := []byte(key)
	//Checking key length
	if len(k) > 32 {

		fmt.Println("Key is too big")
		fmt.Println()
		defer mutex.Unlock()
		return

	}
	//Checking value size in bytes
	v := []byte(fmt.Sprintf("%v", value))
	if binary.Size(v) >= 16000 {

		fmt.Println("Value is too large")
		fmt.Println()
		defer mutex.Unlock()
		//go expire(p)
		return

	}

	//Writing or Appending Data to the file
	if _, err := f.Write([]byte(key + ",\"" + value + "\"," + newT + "\n")); err != nil {

		fmt.Println("File error")
		fmt.Println()
		defer mutex.Unlock()
		go expire(p)
		return

	}

	if err := f.Close(); err != nil {

		fmt.Println("Error in closing the file")
		defer mutex.Unlock()
		//go expire(p)

	}
	defer mutex.Unlock()
	//go expire(p)

}

//Read data from file
func read() {

	mutex.Lock()
	in := bufio.NewReader(os.Stdin)
	var p string
	usr, err := user.Current()
	check(err)
	//Getting file path from user
	fmt.Println("Enter the file path with file name or press enter key")
	fmt.Scanf("%s", &p)

	if len(p) <= 1 {

		p = string(usr.HomeDir) + "/data.csv"

	}
	split := strings.Split(p, "/")
	p = path.Join(split...)

	fmt.Println()
	fmt.Println("Enter a Key to search")
	key, err := in.ReadString('\n')
	check(err)
	key = strings.TrimSpace(key)

	csvfile, err := os.Open("/" + p)
	if err == nil {
		r := csv.NewReader(csvfile)

		// Iterate through the records
		var k int
		//Loop to find if a key is present in the csv file
		for {
			// Read each record from csv

			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Sorry")
				break
			}

			if record[0] == key {

				fmt.Println(record[1])
				k = 1
				break

			}
		}
		if k == 0 {

			fmt.Println()
			fmt.Println("Key not found")
		}

		//Closing the file
		if err := csvfile.Close(); err != nil {

			fmt.Println("Error in closing the file")
		}
		defer mutex.Unlock()
		//go expire(p)
		return

	} else {

		fmt.Println("File Not found")
		defer mutex.Unlock()
		//go expire(p)
		return
	}

}

//Delete a key value pair
func del() {

	mutex.Lock()
	//Struct to store each row values
	type com struct {
		id, json, time string
	}

	//Slice of struct pointers
	data := []*com{}

	in := bufio.NewReader(os.Stdin)
	var p string
	usr, err := user.Current()
	check(err)
	fmt.Println("Enter the file path with file name or press enter key")
	fmt.Scanf("%s", &p)

	if len(p) <= 1 {

		p = string(usr.HomeDir) + "/data.csv"

	}
	split := strings.Split(p, "/")
	p = path.Join(split...)

	//Getting the key to get deleted from user
	fmt.Println("Enter a Key to delete")
	key, err := in.ReadString('\n')
	fmt.Println()
	check(err)
	key = strings.TrimSpace(key)
	var value string
	csvfile, err := os.Open("/" + p)
	if err == nil {
		r := csv.NewReader(csvfile)

		// Iterate through the records
		var k int
		for {
			// Read each record from csv
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Sorry")
			}
			//Appending the structure data to the slice of pointers
			if record[0] != key {

				c := new(com)
				c.id = record[0]
				value = strings.ReplaceAll(record[1], "\"", "\"\"")
				c.json = value
				c.time = record[2]
				data = append(data, c)
			} else {
				k = 1
			}
		}
		//Closing the file
		if err := csvfile.Close(); err != nil {

			fmt.Println("Error in closing the file")
			defer mutex.Unlock()
			return

		}
		if k == 0 {

			fmt.Println("Key not found")
			fmt.Println()
			defer mutex.Unlock()
			return
		}

		f, err := os.Create("/" + p)
		check(err)
		for i := range data {

			if _, err := f.Write([]byte(data[i].id + ",\"" + data[i].json + "\"," + data[i].time + "\n")); err != nil {

				fmt.Println("File error")
				defer mutex.Unlock()
				return
			}

		}
		if err := f.Close(); err != nil {

			fmt.Println("Error in closing the file")
			defer mutex.Unlock()
			return
		}

	} else {

		fmt.Println("File Not found")
		defer mutex.Unlock()
		return
	}

	go expire(p)
	return
}

//Deleting files that has expired
func expire(p string) {
	layout := "2021-01-01 13:22:57.162902067 +0530 IST m=+85.473189500"
	type com struct {
		id, json, time string
	}
	var value string
	csvfile, err := os.Open("/" + p)
	if err == nil {
		r := csv.NewReader(csvfile)

		// Iterate through the records
		data := []*com{}
		for {
			// Read each record from csv
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Sorry")
			}
			//Appending the structure data to the slice of pointers
			t, err := time.Parse(layout, record[2])
			n := time.Now()
			if err == nil {

				if n.Before(t) {

					c := new(com)
					c.id = record[0]
					value = strings.ReplaceAll(record[1], "\"", "\"\"")
					c.json = value
					c.time = record[2]
					data = append(data, c)
				}
			} else {
				c := new(com)
				c.id = record[0]
				value = strings.ReplaceAll(record[1], "\"", "\"\"")
				c.json = value
				c.time = record[2]
				data = append(data, c)
			}
		}
		//Closing the file
		if err := csvfile.Close(); err != nil {

			fmt.Println("Error in closing the file")
			return

		}

		f, err := os.Create("/" + p)
		check(err)
		for i := range data {

			if _, err := f.Write([]byte(data[i].id + ",\"" + data[i].json + "\"," + data[i].time + "\n")); err != nil {

				fmt.Println("File error")
				defer mutex.Unlock()
				return
			}

		}
		if err := f.Close(); err != nil {

			fmt.Println("Error in closing the file")
			return
		}

	} else {

		fmt.Println("File Not found")
		return
	}
	return

}
func main() {

	for {
		var c int
		fmt.Println("Enter your choice:")
		fmt.Println()
		fmt.Println("1.Create 2.Read 3.Delete 4.Exit")
		fmt.Scanf("%d", &c)
		switch c {

		case 1:
			create()
		case 2:
			read()
		case 3:
			del()
		case 4:
			return

		}

	}

}
