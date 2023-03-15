package config

import (
	"flag"
	"fmt"
	"os"
)

func ParseArgs() (s StreamSettings, m MessageBrokerSettings) {
	var videoResolution Resolution
	var videoEncoder VideoEncoder = H264
	var videoBaseFramerate uint
	var videoBaseBitrate uint
	var videoShowCursor bool
	var audioBaseBitrate uint
	var audioBasePacketLossPct uint

	var host string
	var port PortNumber = 5672
	var vHost string
	var username string
	var password string

	flag.Var(&videoResolution, "res", "The resolution to use (required). Should be in the format [WIDTH]x[HEIGHT].")
	flag.Var(&videoEncoder, "venc", "The video encoder to use.")
	flag.UintVar(&videoBaseFramerate, "vfr", 60, "Video base framerate.")
	flag.UintVar(&videoBaseBitrate, "vbr", 52000, "Video base bitrate in kbit/sec.")
	flag.BoolVar(&videoShowCursor, "vsc", true, "Whether to show cursor in recorded screen.")
	flag.UintVar(&audioBaseBitrate, "abr", 64000, "Audio base bitrate in bps.")
	flag.UintVar(&audioBasePacketLossPct, "apl", 5, "Audio base packet loss percentage.")

	flag.StringVar(&host, "rmqhost", "localhost", "RabbitMQ message broker host.")
	flag.Var(&port, "rmqport", "RabbitMQ message broker port. Should be in the range 0-65535.")
	flag.StringVar(&vHost, "rmqvhost", "/", "RabbitMQ virtual host.")
	flag.StringVar(&username, "rmquser", "", "RabbitMQ username (required).")
	flag.StringVar(&password, "rmqpass", "", "RabbitMQ password (required).")

	flag.Parse()

	// check required fields
	if videoResolution.Width == 0 || videoResolution.Height == 0 {
		fmt.Println("Error: the flag -res is required.")
		flag.Usage()
		os.Exit(1)
	}
	if username == "" {
		fmt.Println("Error: the flag -rmquser is required.")
		flag.Usage()
		os.Exit(1)
	}
	if password == "" {
		fmt.Println("Error: the flag -rmqpass is required.")
		flag.Usage()
		os.Exit(1)
	}

	s.AudioBaseBitrate = audioBaseBitrate
	s.AudioBasePacketLossPct = audioBasePacketLossPct
	s.VideoBaseBitrate = videoBaseBitrate
	s.VideoBaseFramerate = videoBaseFramerate
	s.VideoEncoder = videoEncoder
	s.VideoResolution = videoResolution
	s.VideoShowCursor = videoShowCursor
	s.VideoResolution = videoResolution

	m.Host = host
	m.Password = password
	m.Port = port
	m.Username = username
	m.VHost = vHost

	return
}
