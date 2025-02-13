package main

import (
	"discord-dex-screener-bot/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const token = "TOKEN"

type Command string

const (
	Dex Command = "!dex"
)

var MessageTemplate string = "🤖 **%s** | %s \n 🔗 <%s>\n\nCreated: %s\n\n📊 MarketCap: **$%s** | Price: $%s | Liquidity: **$%s**\nVolume 5mins: **$%s** | Volume 1hour: **$%s**\nBuys/Sells 5mins: 🟩 **%d**    🟥 **%d**\n\n📊 Bullx : <https://neo.bullx.io/terminal?chainId=1399811149&address=%s>\n"

func FormatAmount(amount int) string {
	if amount > 1000000 {
		i := amount / 1000000
		str := strconv.Itoa(i)
		str += "M"
		return str
	}
	if amount > 1000 {
		i := amount / 1000
		str := strconv.Itoa(i)
		str += "K"
		return str
	}
	return strconv.Itoa(amount)
}

type Logger struct {
	*log.Logger
}

func NewLogger() Logger {
	return Logger{log.New(os.Stderr, "", log.Ldate|log.Ltime)}
}

func (log Logger) Error(msg string, err error) {
	log.Print("ERROR: " + msg + ": " + err.Error())
}

func (log Logger) Info(msg string) {
	log.Print("INFO: " + msg)
}

func router(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}
	log := NewLogger()

	switch {
	case strings.Contains(message.Content, string(Dex)):
		ca := strings.Trim(message.Content, string(Dex))
		ca = strings.Trim(ca, " ")

		log.Info("received dex command for ca " + ca)

		path := "https://api.dexscreener.com" + "/token-pairs/v1/solana/" + ca

		resp, err := http.Get(path)
		if err != nil {
			log.Error("could not get CA on dex screener", err)
			return
		}
		defer resp.Body.Close()

		log.Info("successfully got CA for dexscreener")

		var response model.Chains
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Error("could not decode json", err)
			return
		}

		var volumeM5 string
		var volumeH1 string
		if response[0].Volume.M5.IsZero() {
			volumeM5 = "0"
		} else {
			volumeM5 = FormatAmount(int(response[0].Volume.M5.IntPart()))
		}

		if response[0].Volume.H1.IsZero() {
			volumeH1 = "0"
		} else {
			volumeH1 = FormatAmount(int(response[0].Volume.H1.IntPart()))
		}

		var output string
		if len(response) > 0 {
			output = fmt.Sprintf(
				MessageTemplate,
				response[0].BaseToken.Symbol,
				response[0].BaseToken.Name,
				response[0].URL,
				time.UnixMilli(int64(response[0].PairCreatedAt)).Format(time.RFC3339),
				FormatAmount(int(response[0].MarketCap.IntPart())),
				response[0].PriceUsd,
				FormatAmount(int(response[0].Liquidity.Usd.IntPart())),
				volumeM5,
				volumeH1,
				response[0].Txns.M5.Buys,
				response[0].Txns.M5.Sells,
				response[0].BaseToken.Address,
			)
		} else {
			output = "CA not found 😞"
		}

		_, err = discord.ChannelMessageSend(message.ChannelID, output)
		if err != nil {
			log.Error("could not send message to channel", err)
			return
		}

		log.Info("sent message to channel")
	}
}

func main() {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	discord.UserAgent = "DiscordBot"

	discord.AddHandler(router)
	discord.Open()
	defer discord.Close()
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
