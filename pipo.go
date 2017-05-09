package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

var (
	slk          *slack.Client
	rtm          *slack.RTM
	games        GameList
	gameDuration = 20 * time.Minute
	token        = os.Getenv("BOT_TOKEN")
)

const (
	botName     = "pipo"
	botID       = "U3RD48GMC"
	botAvatar   = "https://avatars.slack-edge.com/2017-01-12/126139559856_47ebe28f7381fdbb392d_original.png"
	cmdBook     = "book"
	cmdCancel   = "cancel"
	cmdBookings = "bookings"
	cmdStatus   = "status"
	cmdHelp     = "help"
)

func runCleanup() {
	t := time.NewTicker(1 * time.Minute)

	for range t.C {
		sweepGames()
	}
}

func sweepGames() {
	for i, game := range games {

		if game.StartTime.Equal(time.Now()) || game.StartTime.Before(time.Now()) {
			game.InProgress = true
		}

		if game.StartTime.Add(gameDuration).Equal(time.Now()) || game.StartTime.Add(gameDuration).Before(time.Now()) {
			game.InProgress = false

			// remove it
			copy(games[i:], games[i+1:])
			games[len(games)-1] = nil
			games = games[:len(games)-1]
		}

		if !game.InProgress && (time.Now().Add(3 * time.Minute).Equal(game.StartTime)) {
			notify(game)
		}
	}
}

func piporun() {

	slk = slack.New(token)

	_, err := slk.AuthTest()
	if err != nil {
		log.Fatal(err)
	}

	rtm = slk.NewRTM()
	go rtm.ManageConnection()

	go runCleanup()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			pipoStr := "^(?:<@" + botID + ">|pipo)$"
			pipoRegex := regexp.MustCompile("(?i)" + pipoStr)
			pipoCaptureGroups := pipoRegex.FindAllStringSubmatch(ev.Text, -1)
			if pipoCaptureGroups != nil {
				showHelpCommands(ev.Channel)
				continue
			}

			infoStr := "^(?:<@" + botID + ">|pipo)\\s(help|bookings?)$"
			infoRegex := regexp.MustCompile("(?i)" + infoStr)
			infoCaptureGroups := infoRegex.FindAllStringSubmatch(ev.Text, -1)
			if infoCaptureGroups != nil {
				command := infoCaptureGroups[0][1]

				if command == cmdHelp {
					showHelpCommands(ev.Channel)
				} else if command == cmdBookings || command == cmdBookings[:len(cmdBookings)-1] {
					listBookings(ev.Channel)
				}
				continue
			}

			commandStr := "^(?:<@" + botID + ">|pipo)\\s(book|cancel)\\s(<@\\w+>)\\s((?:[0-9]|0[0-9]|1[0-9]|2[0-3])(?:[0-5][0-9]|:[0-5][0-9])?\\s?(?:AM|PM)?)$"
			commandRegex := regexp.MustCompile("(?i)" + commandStr)
			commandCaptureGroups := commandRegex.FindAllStringSubmatch(ev.Text, -1)
			if commandCaptureGroups != nil {
				command := commandCaptureGroups[0][1]
				target := commandCaptureGroups[0][2]
				time := commandCaptureGroups[0][3]

				if command == cmdBook {
					createBooking(ev.Channel, ev.Msg.User, target, time)
				} else if command == cmdCancel {
					cancelBooking(ev.Channel, ev.Msg.User, target, time)
				}
				continue
			}

			errorStr := "^(?:<@" + botID + ">|pipo)\\s\\w+"
			errorRegex := regexp.MustCompile("(?i)" + errorStr)
			errorCaptureGroups := errorRegex.FindAllStringSubmatch(ev.Text, -1)
			if errorCaptureGroups != nil {
				showErrorReponse(ev.Channel)
			}
		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())
		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return
		}
	}
}

func notify(game *Game) {
	player1 := *game.Player1
	player2 := *game.Player2
	p1Message := fmt.Sprintf("Hey %s! You have a scheduled ping pong match against %s at %s. Don't be late!", player1.Name, player2.Name, formatTime(game.StartTime))
	p2Message := fmt.Sprintf("Hey %s! You have a scheduled ping pong match against %s at %s. Don't be late!", player2.Name, player1.Name, formatTime(game.StartTime))

	postMessage("@"+player1.ID, p1Message, true)
	postMessage("@"+player2.ID, p2Message, true)
}

func showHelpCommands(channel string) {
	message := "Each game booking is 20 minutes.\n\n\n" +
		"To view all bookings, just say: \n`pipo bookings`\n\n\n" +
		"To make a booking, just say: \n`pipo book [@opponent] [time]`\n" +
		"`EXAMPLE: pipo book @pipo 3:15 PM`\n\n\n" +
		"To cancel a booking, just say: \n`pipo cancel [@opponent] [time]`\n" +
		"`EXAMPLE: pipo cancel @pipo 3:15 PM`"
	postMessage(channel, message, false)
}

func showErrorReponse(channel string) {
	message := "Sorry, I don't understand. For a list of commands, just say:\n `pipo help`"
	postMessage(channel, message, false)
}

func listBookings(channel string) {
	var message string

	if len(games) > 0 {
		message = "I have the following games booked: ```"
		for _, game := range games {
			message += formatTime(game.StartTime) + " - " + game.Player1.Name + " vs " + game.Player2.Name + "\n"
		}
		message += "```"
	} else {
		message += "It doesn't look like I have any games booked right now."
	}

	postMessage(channel, message, false)
}

func createBooking(channel, player, opponent, gameTime string) {
	user1, err := slk.GetUserInfo(player)
	if err != nil {
		log.Println(err)
		return
	}

	opponentID := opponent[2 : len(opponent)-1]
	user2, err := slk.GetUserInfo(opponentID)
	if err != nil {
		log.Println(err)
		return
	}

	startTime, err := parseTime(gameTime)
	if err != nil {
		log.Println(err)
		showErrorReponse(channel)
		return
	}

	now := time.Now()
	if startTime.Before(now) {
		message := "Hey, you can't book a game in the past!"
		postMessage(channel, message, false)
		return
	}

	newEndTime := startTime.Add(gameDuration)

	for _, game := range games {
		gameEndTime := game.StartTime.Add(gameDuration)

		if (startTime.After(game.StartTime) && startTime.Before(gameEndTime)) ||
			(newEndTime.After(game.StartTime) && newEndTime.Before(gameEndTime)) ||
			(startTime.Equal(game.StartTime) && newEndTime.Equal(gameEndTime)) {
			message := "Unfortunately, there was a booking conflict. To see a list of bookings, just say:\n `pipo bookings`"
			postMessage(channel, message, false)
			return
		}
	}

	player1 := &Player{
		ID:     user1.ID,
		Name:   user1.RealName,
		Avatar: user1.Profile.ImageOriginal,
		Score:  0,
	}

	player2 := &Player{
		ID:     user2.ID,
		Name:   user2.RealName,
		Avatar: user2.Profile.ImageOriginal,
		Score:  0,
	}

	game := &Game{
		Player1:    player1,
		Player2:    player2,
		StartTime:  startTime,
		InProgress: false,
	}
	games = append(games, game)
	sort.Sort(games)

	message := "Okay! I've made a booking for " + player1.Name + " against " + player2.Name + " at " + formatTime(startTime)
	postMessage(channel, message, false)
}

func cancelBooking(channel, player, opponent, gameTime string) {
	startTime, err := parseTime(gameTime)
	if err != nil {
		log.Println(err)
		return
	}

	cancelled := false
	opponentID := opponent[2 : len(opponent)-1]
	for i, game := range games {
		if game.Player1.ID == player && game.Player2.ID == opponentID && game.StartTime == startTime {
			copy(games[i:], games[i+1:])
			games[len(games)-1] = nil
			games = games[:len(games)-1]
			cancelled = true
			break
		}
	}

	message := ""
	if cancelled {
		message = "Booking cancelled!"
	} else {
		message = "Sorry, I could find the game you mentioned."
	}

	postMessage(channel, message, false)
}

func parseTime(timeStr string) (time.Time, error) {
	timeRegexStr := "^((?:[0-9]|0[0-9]|1[0-9]|2[0-3])(?:[0-5][0-9]|:[0-5][0-9])?)\\s?(AM|PM)?$"
	timeRegex := regexp.MustCompile("(?i)" + timeRegexStr)
	captureGroups := timeRegex.FindAllStringSubmatch(timeStr, -1)

	now := time.Now()
	timeSegment := captureGroups[0][1]
	timeSuffix := strings.ToUpper(captureGroups[0][2])

	if !strings.Contains(timeSegment, ":") {
		if len(timeSegment) == 1 || len(timeSegment) == 2 {
			timeSegment += ":00"
		} else if len(timeSegment) == 3 {
			timeSegment = timeSegment[0:1] + ":" + timeSegment[1:]
		} else if len(timeSegment) == 4 {
			timeSegment = timeSegment[0:2] + ":" + timeSegment[2:]
		}
	}

	if timeSuffix == "" {
		if now.Hour() >= 12 {
			timeSuffix = "PM"
		} else {
			timeSuffix = "AM"
		}
	}

	fmtTimeStr := fmt.Sprintf("%d-%d-%d %s %s", now.Year(), now.Month(), now.Day(), timeSegment, timeSuffix)
	fmtTime, err := time.ParseInLocation("2006-1-2 3:04 PM", fmtTimeStr, time.Local)
	if err != nil {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), err
	}

	return fmtTime, nil
}

func formatTime(rawTime time.Time) string {
	return rawTime.Format("3:04 PM")
}

func postMessage(channel, message string, asUser bool) {
	rtm.PostMessage(channel, message, slack.PostMessageParameters{
		Username: botName,
		IconURL:  botAvatar,
		AsUser:   asUser,
	})
}

func printGames() {
	for _, game := range games {
		log.Printf("%+v", game.Player1)
		log.Printf("%+v", game.Player2)
		log.Printf("%+v", game.StartTime)
		log.Printf("%+v", game.InProgress)
	}
	log.Println("---------------------------------------")
}
