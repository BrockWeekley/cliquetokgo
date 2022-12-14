package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type VideoResponse struct {
	Urls []string `json:"urls"`
}

func main() {
	StartServer()
}

func StartServer() {
	apiServerMux := http.NewServeMux()
	apiServer := http.Server{
		Addr:    fmt.Sprintf(":%v", 5001),
		Handler: apiServerMux,
	}
	apiServerMux.HandleFunc("/api/v2", handler)
	apiServerMux.HandleFunc("/api/v2/getVideos", getVideos)
	CheckForError(apiServer.ListenAndServe())
}

func getVideos(w http.ResponseWriter, r *http.Request) {
	setWriter(w)
	url := buildUrl(r.URL.Query().Get("tag"))
	resp, err := http.Get(url)
	CheckForError(err)

	fmt.Println("Response code from TikTok: ")
	fmt.Println(resp.StatusCode)
	w.WriteHeader(resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		CheckForError(err)
		videoUrls := filterVideos(string(bodyBytes))
		fmt.Println("Size of found video urls: ")
		fmt.Println(len(videoUrls))
		fmt.Println("First url in list: ")
		fmt.Println(videoUrls[0])
		response := VideoResponse{Urls: videoUrls}
		fmt.Println(response.Urls)
		CheckForError(json.NewEncoder(w).Encode(response))
	}
}

func filterVideos(body string) []string {
	rex := regexp.MustCompile("https:\\/\\/v16m-default.tiktokcdn-us.com\\/[-a-zA-Z0-9()@:%_\\+.~#?&//=]+ve[^\"]+")
	foundVideos := rex.FindAllString(body, -1)
	return removeDuplicates(foundVideos)
}

func removeDuplicates(foundVideos []string) []string {
	existing := make(map[string]bool)
	var filteredVideos []string

	for _, entry := range foundVideos {
		if _, value := existing[entry]; !value {
			existing[entry] = true
			filteredVideos = append(filteredVideos, entry)
		}
	}
	return filteredVideos
}

func buildUrl(tag string) string {
	fmt.Println("Received a new request with tag: ")
	fmt.Println(tag)
	return "https://us.tiktok.com/api/topic/item_list/" +
		"?aid=1988&" +
		"app_language=en&" +
		"app_name=tiktok_web&" +
		"battery_info=0.69&" +
		"browser_language=en-US&" +
		"browser_name=Mozilla&" +
		"browser_online=true&" +
		"browser_platform=MacIntel&" +
		"channel=tiktok_web&" +
		"cookie_enabled=true&" +
		"count=9&" +
		"device_id=7146004948433700398&" +
		"device_platform=web_pc&" +
		"focus_state=true&" +
		"from_page=topics_gaming&" +
		"history_len=9&" +
		"is_fullscreen=false&" +
		"is_page_visible=true&" +
		"language=en&" +
		"os=mac&" +
		"priority_region=&" +
		"referer=&" +
		"region=US&" +
		"screen_height=1120&" +
		"screen_width=1792&" +
		"topic=" + tag + "&" +
		"tz_name=America%2FChicago&" +
		"webcast_language=en"
}

func handler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println(w)
}

func setWriter(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func CheckForError(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
}
