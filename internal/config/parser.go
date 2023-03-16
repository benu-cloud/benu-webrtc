package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"
)

func ParseArgs() (s StreamSettings, m MessageBrokerSettings) {
	// try to load env variables if they exist
	godotenv.Load()

	var videoResolution Resolution
	var videoEncoder VideoEncoder = H264
	var videoBaseFramerate uint
	var videoBaseBitrate uint
	var videoShowCursor bool
	var audioBaseBitrate uint
	var audioBasePacketLossPct uint

	var host string
	var port PortNumber = PortNumber(5672)
	var vHost string
	var username string
	var password string

	flag.Var(&videoResolution, "vresolution", "The resolution to use (required). Should be in the format [WIDTH]x[HEIGHT].")
	flag.Var(&videoEncoder, "vencoder", "The video encoder to use.")
	flag.UintVar(&videoBaseFramerate, "vframerate", 60, "Video base framerate.")
	flag.UintVar(&videoBaseBitrate, "vbitrate", 52000, "Video base bitrate in kbit/sec.")
	flag.BoolVar(&videoShowCursor, "vcursor", true, "Whether to show cursor in recorded screen.")
	flag.UintVar(&audioBaseBitrate, "abitrate", 64000, "Audio base bitrate in bps.")
	flag.UintVar(&audioBasePacketLossPct, "apacketlosspct", 5, "Audio base packet loss percentage.")

	flag.StringVar(&host, "rmqhost", "localhost", "RabbitMQ message broker host.")
	flag.Var(&port, "rmqport", "RabbitMQ message broker port. Should be in the range 0-65535.")
	flag.StringVar(&vHost, "rmqvhost", "/", "RabbitMQ virtual host.")
	flag.StringVar(&username, "rmqusername", "", "RabbitMQ username (required).")
	flag.StringVar(&password, "rmqpassword", "", "RabbitMQ password (required).")

	flag.Parse()

	// check required fields
	if videoResolution.Width == 0 || videoResolution.Height == 0 {
		fmt.Println("Error: the flag -vresolution is required.")
		flag.Usage()
		os.Exit(1)
	}
	if username == "" {
		fmt.Println("Error: the flag -rmqusername is required.")
		flag.Usage()
		os.Exit(1)
	}
	if password == "" {
		fmt.Println("Error: the flag -rmqpassword is required.")
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
