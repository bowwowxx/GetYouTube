package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type YoutubeVideo struct {
	Id       string
	Info     url.Values
	Filename string
	Url      string
}

var format = map[string]string{
	"5":   "flv",
	"6":   "flv",
	"13":  "3gp",
	"17":  "3gp",
	"18":  "mp4",
	"22":  "mp4",
	"34":  "flv",
	"35":  "flv",
	"36":  "3gp",
	"37":  "mp4",
	"38":  "mp4",
	"43":  "webm",
	"44":  "webm",
	"45":  "webm",
	"46":  "webm",
	"82":  "mp4",
	"83":  "mp4",
	"84":  "mp4",
	"85":  "mp4",
	"100": "webm",
	"101": "webm",
	"102": "webm",
	"92":  "mp4",
	"93":  "mp4",
	"94":  "mp4",
	"95":  "mp4",
	"96":  "mp4",
	"132": "mp4",
	"151": "mp4",
	"133": "mp4",
	"134": "mp4",
	"135": "mp4",
	"136": "mp4",
	"137": "mp4",
	"138": "mp4",
	"160": "mp4",
	"264": "mp4",
	"139": "m4a",
	"140": "m4a",
	"141": "m4a",
	"167": "webm",
	"168": "webm",
	"169": "webm",
	"170": "webm",
	"218": "webm",
	"219": "webm",
	"242": "webm",
	"243": "webm",
	"244": "webm",
	"245": "webm",
	"246": "webm",
	"247": "webm",
	"248": "webm",
	"171": "webm",
	"172": "webm",
}

const baseVideoInfoUrl string = "http://www.youtube.com/get_video_info"
const videoInformationValid string = "ok"

func downloadVideoInformation(video *YoutubeVideo) string {

	params := url.Values{}
	params.Add("el", "vevo")
	params.Add("video_id", video.Id)
	var videoInfoUrl string = baseVideoInfoUrl + "?" + params.Encode()

	response, err := http.Get(videoInfoUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	rawVideoInfo, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(rawVideoInfo)
}

func parseRawVideoInfo(rawData string) url.Values {
	values, err := url.ParseQuery(rawData)
	if err != nil {
		log.Fatal(err)
	}
	return values
}

func (video *YoutubeVideo) GetVideoInformation() {
	log.Print("Getting video information of: ", video.Id)

	rawVideoInfo := downloadVideoInformation(video)
	video.Info = parseRawVideoInfo(rawVideoInfo)

	if video.Info["status"][0] != "ok" {
		log.Fatal("Status ko ", video.Info["reason"][0])
	}
}

func getFileName(fileName string, itag string) string {
	fileName = strings.Replace(fileName, "(", "", -1)
	fileName = strings.Replace(fileName, ")", "", -1)
	fileName = strings.Trim(fileName, " ")
	fileName = fileName + "." + format[itag]
	return fileName
}

func (video *YoutubeVideo) SelectBestLink() {
	var choosedVideoFormat url.Values
	rawVideoStreamMap := strings.Split(video.Info["url_encoded_fmt_stream_map"][0], ",")
	for _, element := range rawVideoStreamMap {
		parsedVideoStreamMap := parseRawVideoInfo(string(element))

		// fmt.Println(parsedVideoStreamMap["itag"][0] + " " + parsedVideoStreamMap["quality"][0])
		// fmt.Println(parsedVideoStreamMap["type"][0])

		if parsedVideoStreamMap["itag"][0] == "18" {
			choosedVideoFormat = parsedVideoStreamMap
		}

	}

	video.Filename = getFileName(video.Info["title"][0], choosedVideoFormat["itag"][0])
	video.Url = choosedVideoFormat["url"][0]
}
