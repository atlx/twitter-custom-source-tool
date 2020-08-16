package main

import (
	"bufio"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

const (
	OAuthConsumerKey = "OAUTH_CONSUMER_KEY"
	OAuthConsumerSecret = "OAUTH_CONSUMER_SECRET"
	OAuthToken = "OAUTH_TOKEN"
	OAuthTokenSecret = "OAUTH_TOKEN_SECRET"
)

func end() {
	fmt.Println("Program finished. Press ENTER to close this window...")
	_, _ = fmt.Scanln()
	os.Exit(0)
}

func waitForInput(msg string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg)
	text, _ := reader.ReadString('\n')
	return text
}

func checkResponse(str string, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(str), prefix)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load .env file; it might be missing. Add it to your project root.")
		end()
		return
	}
	for _, v := range [4]string{
		OAuthConsumerKey,
		OAuthConsumerSecret,
		OAuthToken,
		OAuthTokenSecret,
	} {
		if os.Getenv(v) == "" {
			fmt.Println(v + " is missing from your environment configuration. Ensure it is set, then try again.")
			end()
			return
		}
	}
	config := oauth1.NewConfig(os.Getenv(OAuthConsumerKey), os.Getenv(OAuthConsumerSecret))
	token := oauth1.NewToken(os.Getenv(OAuthToken), os.Getenv(OAuthTokenSecret))
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	fmt.Println("=========\nTwitter Custom Client Tool\n=========")

	tweetText := waitForInput("Enter tweet text, then press ENTER: ")
	wantsToSetLoc := waitForInput("Do you want to set a custom Tweet location (latitude/longitude)? [y/n] ")
	var lat *float64 = nil
	var long *float64 = nil
	if checkResponse(wantsToSetLoc, "y") {
		wantsPremadeList := waitForInput("Do you want to select from a pre-made list of locations (n for custom location)? [y/n] ")
		if checkResponse(wantsPremadeList, "y") {
			var placesStr []string
			for i, p := range Places {
				placesStr = append(placesStr, fmt.Sprintf("[%d] %s (lat: %f, long: %f)", i + 1, p.Name, p.Lat, p.Long))
			}
			fmt.Println(strings.Join(placesStr, "\n"))
			var resAsInt int
			fmt.Println("Type a number then press ENTER: ")
			_, e := fmt.Scanln(&resAsInt)
			if e != nil {
				fmt.Println("Could not parse number.")
				end()
				return
			}

			if resAsInt < 0 || resAsInt > len(Places) {
				fmt.Println("Out of bounds.")
				end()
				return
			}
			lat = &Places[resAsInt - 1].Lat
			long = &Places[resAsInt - 1].Long
		} else {
			args := waitForInput("Input desired latitude and longitude number, separated by a space (5.000 6.000): ")
			latLongArray := strings.Split(args, " ")
			if len(latLongArray) != 2 {
				fmt.Println("Latitude/longitude not correctly formatted.")
				end()
				return
			}
			latNum, err := strconv.ParseFloat(latLongArray[0], 64)
			if err != nil {
				fmt.Println("Failed to parse latitude.")
				end()
				return
			}
			longNum, err := strconv.ParseFloat(strings.TrimRight(latLongArray[1], "\n\r"), 64)
			if err != nil {
				fmt.Println("Failed to parse longitude.")
				end()
				return
			}
			fmt.Println(latNum)
			fmt.Println(longNum)
			lat = &latNum
			long = &longNum
		}
	}

	fmt.Println("Sending tweet.")

	tweet, _, err := client.Statuses.Update(tweetText, &twitter.StatusUpdateParams{
		Status:             "",
		InReplyToStatusID:  0,
		PossiblySensitive:  nil,
		Lat:                lat,
		Long:               long,
		PlaceID:            "",
		TrimUser:           nil,
		MediaIds:           nil,
		TweetMode:          "",
	})
	if err != nil {
		fmt.Println("Failed to send tweet. " + err.Error())
		end()
		return
	}
	fmt.Println("Done! Access your tweet at:\n" + fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr))
	end()
	return
}
