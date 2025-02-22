package device

import (
	"github.com/AlexxIT/go2rtc/internal/api"
	"github.com/AlexxIT/go2rtc/pkg/core"
	"net/url"
	"os/exec"
	"regexp"
)

func queryToInput(query url.Values) string {
	video := query.Get("video")
	audio := query.Get("audio")

	if video == "" && audio == "" {
		return ""
	}

	// https://ffmpeg.org/ffmpeg-devices.html#dshow
	input := "-f dshow"

	if video != "" {
		video = indexToItem(videos, video)

		for key, value := range query {
			switch key {
			case "resolution":
				input += " -video_size " + value[0]
			case "video_size", "framerate", "pixel_format":
				input += " -" + key + " " + value[0]
			}
		}
	}

	if audio != "" {
		audio = indexToItem(audios, audio)

		for key, value := range query {
			switch key {
			case "sample_rate", "sample_size", "channels", "audio_buffer_size":
				input += " -" + key + " " + value[0]
			}
		}
	}

	if video != "" {
		input += ` -i video="` + video + `"`

		if audio != "" {
			input += `:audio="` + audio + `"`
		}
	} else {
		input += ` -i audio="` + audio + `"`
	}

	return input
}

func deviceInputSuffix(video, audio string) string {
	switch {
	case video != "" && audio != "":
		return `video="` + video + `":audio=` + audio + `"`
	case video != "":
		return `video="` + video + `"`
	case audio != "":
		return `audio="` + audio + `"`
	}
	return ""
}

func initDevices() {
	cmd := exec.Command(
		Bin, "-hide_banner", "-list_devices", "true", "-f", "dshow", "-i", "",
	)
	b, _ := cmd.CombinedOutput()

	re := regexp.MustCompile(`"([^"]+)" \((video|audio)\)`)
	for _, m := range re.FindAllStringSubmatch(string(b), -1) {
		name := m[1]
		kind := m[2]

		stream := api.Stream{
			Name: name, URL: "ffmpeg:device?" + kind + "=" + name,
		}

		switch kind {
		case core.KindVideo:
			videos = append(videos, name)
			stream.URL += "#video=h264#hardware"
		case core.KindAudio:
			audios = append(audios, name)
		}

		streams = append(streams, stream)
	}
}
