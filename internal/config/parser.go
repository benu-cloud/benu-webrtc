package config

import (
	"fmt"
	"os"
	"time"

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

	var rmqhost string
	var rmqport PortNumber = PortNumber(5672)
	var rmqvHost string
	var rmqusername string
	var rmqpassword string
	var rmqpublishTimeoutSeconds uint

	flag.Var(&videoResolution, "vresolution", "The resolution to use (required). Should be in the format [WIDTH]x[HEIGHT].")
	flag.Var(&videoEncoder, "vencoder", "The video encoder to use.")
	flag.UintVar(&videoBaseFramerate, "vframerate", 60, "Video base framerate.")
	flag.UintVar(&videoBaseBitrate, "vbitrate", 52000, "Video base bitrate in kbit/sec.")
	flag.BoolVar(&videoShowCursor, "vcursor", true, "Whether to show cursor in recorded screen.")
	flag.UintVar(&audioBaseBitrate, "abitrate", 64000, "Audio base bitrate in bps.")
	flag.UintVar(&audioBasePacketLossPct, "apacketlosspct", 5, "Audio base packet loss percentage. Should be in range 0-100.")

	flag.StringVar(&rmqhost, "rmqhost", "localhost", "RabbitMQ message broker host.")
	flag.Var(&rmqport, "rmqport", "RabbitMQ message broker port. Should be in the range 0-65535.")
	flag.StringVar(&rmqvHost, "rmqvhost", "", "RabbitMQ virtual host.")
	flag.StringVar(&rmqusername, "rmqusername", "", "RabbitMQ username (required).")
	flag.StringVar(&rmqpassword, "rmqpassword", "", "RabbitMQ password (required).")
	flag.UintVar(&rmqpublishTimeoutSeconds, "rmqtimeout", 5, "RabbitMQ publish timeout in seconds")

	flag.Parse()

	// check required fields
	if videoResolution.Width == 0 || videoResolution.Height == 0 {
		fmt.Println("Error: the flag -vresolution is required.")
		flag.Usage()
		os.Exit(1)
	}
	if audioBasePacketLossPct >= 100 {
		fmt.Println("Error: the flag -apacketlosspct is less than 100")
		flag.Usage()
		os.Exit(1)
	}
	if rmqusername == "" {
		fmt.Println("Error: the flag -rmqusername is required.")
		flag.Usage()
		os.Exit(1)
	}
	if rmqpassword == "" {
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

	m.Host = rmqhost
	m.Password = rmqpassword
	m.Port = rmqport
	m.Username = rmqusername
	m.VHost = rmqvHost
	m.PublishTimeout = time.Second * time.Duration(rmqpublishTimeoutSeconds)

	return
}
