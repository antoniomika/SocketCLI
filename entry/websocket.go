// Package entry maintains the main socketcli event loop
package entry

import (
	"encoding/json"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/antoniomika/socketcli/utils"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	requests   = []utils.EventData{}
	words      = map[string]*utils.WordPlace{}
	wordplaces = []*utils.WordPlace{}
)

// Start initializes the websocket handler.
func Start() {
	// Allow us to catch SIGINT's
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	logrus.Printf("Intitializing websocket connection to: %s", viper.GetString("websocket-address"))

	c, _, err := websocket.DefaultDialer.Dial(viper.GetString("websocket-address"), nil)
	if err != nil {
		logrus.Fatal("error connecting to websocket:", err)
	}
	defer c.Close()

	// done exists to allow us to control the separate websocket listener and event loop routines.
	done := make(chan struct{})

	// Setup a regex to extract only the words sans punctuation and fail if it doesn't compile.
	reg, err := regexp.Compile(`\b[-.,()&$#!\[\]{}"']+\B|\B[-.,()&$#!\[\]{}"']+\b`)
	if err != nil {
		logrus.Fatal("error compiling regex:", err)
	}

	// Main routine to load websocket data.
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}

			parsedMessage := utils.JSONResponse{}
			err = json.Unmarshal(message, &parsedMessage)
			if err != nil {
				return
			}

			requests = append(requests, utils.EventData{
				Response: parsedMessage,
				RealTime: time.Now(),
			})

			logrus.Debugln(parsedMessage)

			w := strings.Fields(strings.ToLower(reg.ReplaceAllString(parsedMessage.Message, "")))

			for _, v := range w {
				word := v
				if _, ok := words[word]; ok {
					words[word].Count++
				} else {
					words[word] = &utils.WordPlace{
						Word:     word,
						RealWord: v,
						Count:    1,
					}
					wordplaces = append(wordplaces, words[word])
				}
			}
		}
	}()

	// Print our aggregate states when the function returns.
	defer func() {
		dur := requests[len(requests)-1].RealTime.Sub(requests[0].RealTime)

		logrus.Println("Top 10 words:")

		sort.Slice(wordplaces, func(i, j int) bool {
			return wordplaces[i].Count > wordplaces[j].Count
		})

		for i := 0; i < 10; i++ {
			logrus.Printf("%s: %d", wordplaces[i].Word, wordplaces[i].Count)
		}

		logrus.Println("Average events per min:", float64(len(requests))/dur.Minutes())
	}()

	// Main loop to control operation.
	for {
		select {
		case <-time.After(viper.GetDuration("stop-after-time")):
			if viper.GetBool("stop-after") {
				close(interrupt)
			}
		case <-done:
			return
		case <-interrupt:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
