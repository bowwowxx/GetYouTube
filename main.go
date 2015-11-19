package main

import (
	"log"
	"os"

	"./util"
)

func handleVideo(videoId string) {
	log.Print("Youtube Id: ", videoId)

	video := util.YoutubeVideo{Id: videoId}
	video.GetVideoInformation()
	video.SelectBestLink()
	util.DownloadFile(video.Url, video.Filename)

	log.Print("ok ....")
}

func main() {
	if len(os.Args) == 2 {
		handleVideo(os.Args[1])
	} else {
		log.Fatal("please input youtube id ex:1GJA4c_SBT0")
	}
}
