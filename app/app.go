package app

import (
	"flag"
	"os"

	"github.com/gianebao/smsx/shorten"
	"github.com/gianebao/smsx/sms"
)

var (
	port string

	nexmo = sms.Nexmo{}
	bitly = shorten.Bitly{}
)

// Init initializes parameters in starting the application
func Init() {
	flag.StringVar(&port, "port", "",
		"listening port")
	flag.StringVar(&nexmo.APIKey, "nexmoapikey", os.Getenv("NEXMOAPIKEY"),
		"(NEXMOAPIKEY) Nexmo service API key")
	flag.StringVar(&nexmo.APISecret, "nexmoapisecret", os.Getenv("NEXMOAPISECRET"),
		"(NEXMOAPISECRET) Nexmo service API secret")
	flag.StringVar(&bitly.Username, "bitlyusername", os.Getenv("BITLYUSERNAME"),
		"(BITLYUSERNAME) Bit.ly account username")
	flag.StringVar(&bitly.Password, "bitlypassword", os.Getenv("BITLYPASSWORD"),
		"(BITLYPASSWORD) Bit.ly account password")
}
