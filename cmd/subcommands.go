package main

import "log"

func CreateLink(appData AppData) {
	// GET a JSON object
	// id := 12
	// var post placeholder
	// resp, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id))
	// if err != nil {
	// 	fmt.Println("could not connect to jsonplaceholder.typicode.com:", err)
	// }
	// defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// json.Unmarshal(body, &post)
	// fmt.Println(post.Title)
	log.SetFlags(0)
	log.Fatalln("Create")
}

func DeleteLink(appData AppData) {
	log.SetFlags(0)
	log.Fatalln("Delete")
}

func ExpandLink(appData AppData) {
	log.SetFlags(0)
	log.Fatalln("Expand")
}

func GetAll(appData AppData) {
	log.SetFlags(0)
	log.Fatalln("Get all")
}
