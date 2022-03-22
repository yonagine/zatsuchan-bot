package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const uri = "https://kittoco.com/kyuukyoku-situmon/"

var (
	Token string
)

type Entry struct {
	Title    string `json:"title"`
	Question string `json:"question"`
}

const entryLimit int = 60

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	fmt.Println("Hello, 世界。")

	// scraper.Scrape(uri)

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("An error was encountered when creating a Discord session.", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("An error was encountered when opening the connection.", err)
		return
	}

	fmt.Println("ザツちゃん登場します！ Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!zatsu" {
		jsonFile, err := os.Open("entries.json")
		if err != nil {
			fmt.Println(err)
		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var entries []Entry

		err = json.Unmarshal(byteValue, &entries)
		if err != nil {
			fmt.Println(err)
		}

		var randomNumber = rand.Intn(entryLimit)
		var message = entries[randomNumber].Title + "\n" + entries[randomNumber].Question

		_, err = s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
