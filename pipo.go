package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
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
	botName         = "pipo"
	botID           = "U3RD48GMC"
	botAvatar       = "https://avatars.slack-edge.com/2017-01-12/126139559856_47ebe28f7381fdbb392d_original.png"
	cmdBook         = "book"
	cmdCancel       = "cancel"
	cmdBookings     = "bookings"
	cmdLeaderboards = "leaderboards"
	cmdStatus       = "status"
	cmdHelp         = "help"
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

		if game.StartTime.Add(gameDuration).Before(time.Now()) {
			game.InProgress = false

			// remove it
			copy(games[i:], games[i+1:])
			games[len(games)-1] = nil // or the zero value of T
			games = games[:len(games)-1]
		}

		if !game.InProgress && (game.StartTime.Equal(time.Now().Add(3*time.Minute)) || game.StartTime.Before(time.Now().Add(3*time.Minute))) {
			Notify(game)
		}
	}
}

func piporun() {

	slk := slack.New(token)

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
			regexStr := "^@?pipo(?:\\s(help|bookings))?$"
			regex := regexp.MustCompile("(?i)" + regexStr)
			captureGroups := regex.FindAllStringSubmatch(ev.Text, -1)
			// log.Printf("CAPTUREGROUPS = |%q|", captureGroups)

			if captureGroups != nil {
				command := captureGroups[0][1]

				if command == cmdHelp {
					showHelpCommands(ev.Channel)
				} else if command == cmdBookings {
					listBookings(ev.Channel)
				}
			} else {
				regexStr2 := "^@?pipo\\s(book|cancel)\\s(<@\\w+>)\\s((?:[0-9]|0[0-9]|1[0-9]|2[0-3]):?(?:[0-5][0-9])?\\s?(?:A\\.?M\\.?|P\\.?M\\.?)?)$"
				regex2 := regexp.MustCompile("(?i)" + regexStr2)
				captureGroups2 := regex2.FindAllStringSubmatch(ev.Text, -1)
				// log.Printf("CAPTUREGROUPS2 = |%q|", captureGroups2)

				if captureGroups2 != nil {
					command2 := captureGroups2[0][1]
					target := captureGroups2[0][2]
					time := captureGroups2[0][3]

					if command2 == cmdBook {
						createBooking(ev.Channel, ev.Msg.User, target, time)
					} else if command2 == cmdCancel {
						cancelBooking(ev.Channel, ev.Msg.User, target, time)
					}
				}
			}

			// if strings.Contains(ev.Text, botName) || strings.Contains(ev.Text, botID) {
			// 	player1, err := slk.GetUserInfo(ev.Msg.User)
			// 	if err != nil {
			// 		log.Println(err)
			// 		continue
			// 	}
			//
			// 	tokens := strings.Split(ev.Text, " ")
			// 	if tokens[0] == botName || tokens[0] == "<@"+botID+">" {
			// 		if len(tokens) > 1 {
			// 			command := tokens[1]
			//
			// 			if len(tokens) == 2 {
			// 				if command == cmdHelp {
			// 					showHelpCommands(ev.Channel)
			// 				} else if command == cmdBookings {
			// 					listBookings(ev.Channel)
			// 				} else if command == cmdLeaderboards {
			// 					showLeaderboard(ev.Channel)
			// 				} else if command == cmdStatus {
			// 					checkTableStatus(ev.Channel)
			// 				} else {
			// 					showErrorReponse(ev.Channel)
			// 				}
			// 			} else if len(tokens) == 4 && command == cmdBook {
			// 				startTime, err := parseTime(tokens[3])
			// 				if err != nil {
			// 					log.Println(err)
			// 					showErrorReponse(ev.Channel)
			// 					continue
			// 				}
			//
			// 				player2ID := tokens[2][2 : len(tokens[2])-1]
			// 				player2, err := slk.GetUserInfo(player2ID)
			// 				if err != nil {
			// 					log.Println(err)
			// 					continue
			// 				}
			//
			// 				createBooking(ev.Channel, player1, player2, startTime)
			// 			} else {
			// 				showErrorReponse(ev.Channel)
			// 			}
			// 		} else {
			// 			showHelpCommands(ev.Channel)
			// 		}
			// 	}
			// }
		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())
		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return
		}
	}
}

// Notify - notify users that their game is beginning soon
func Notify(game *Game) {
	player1 := *game.Player1
	player2 := *game.Player2
	p1Message := "Hey " + player1.Name + "! You have a scheduled ping pong match against " + player2.Name + " at " + formatTime(game.StartTime) + ". Don't be late!"
	p2Message := "Hey " + player2.Name + "! You have a scheduled ping pong match against " + player1.Name + " at " + formatTime(game.StartTime) + ". Don't be late!"

	postMessage("@"+player1.ID, p1Message, true)
	postMessage("@"+player2.ID, p2Message, true)
}

func showHelpCommands(channel string) {
	message := "To view all bookings, just say: \n`pipo bookings`\n\n\n" +
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

	now := time.Now().UTC().Add(-6 * time.Hour) // Hack to get times in UTC and account for 6 hour difference
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
}

func parseTime(timeStr string) (time.Time, error) {
	suffix := ""
	periodSuffix := ""
	pmSuffix := false

	// Check if AM / PM exists
	if len(timeStr) > 2 {
		suffix = strings.ToUpper(timeStr[len(timeStr)-2:])
	}

	// Check if A.M. / P.M. exists
	if len(timeStr) > 4 {
		periodSuffix = strings.ToUpper(timeStr[len(timeStr)-4:])
	}

	if suffix == "AM" || suffix == "PM" {
		timeStr = timeStr[0 : len(timeStr)-2]
	} else if periodSuffix == "A.M." || periodSuffix == "P.M." {
		timeStr = timeStr[0 : len(timeStr)-4]
	}
	pmSuffix = (suffix == "PM" || periodSuffix == "P.M.")

	if strings.Count(timeStr, ":") == 0 {
		timeInt, err := strconv.Atoi(timeStr)
		if err != nil {
			return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), err
		}

		if timeStr[0] != '0' && timeInt >= 1 && timeInt <= 9 {
			timeStr = "0" + timeStr
		}

		timeStr += ":00:00"
	} else if strings.Count(timeStr, ":") == 1 && len(timeStr) > 1 {
		timeStr += ":00"
	}

	// Get and format the current date
	currentDate := time.Now().Local()
	dateStr := currentDate.Format("2006-01-02")

	// Combine the date and time
	dateTimeStr := dateStr + "T" + timeStr + "Z"
	dateTime, err := time.Parse(time.RFC3339, dateTimeStr)
	if err != nil {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), err
	}

	if pmSuffix {
		dateTime = dateTime.Add(12 * time.Hour)
	}

	return dateTime, nil
}

func formatTime(rawTime time.Time) string {
	return rawTime.Format("15:04")
}

func postMessage(channel, message string, asUser bool) {
	rtm.PostMessage(channel, message, slack.PostMessageParameters{
		Username: botName,
		IconURL:  botAvatar,
		AsUser:   asUser,
	})
}
