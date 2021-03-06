/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type UpsertUsers struct {
	Users []User `json:"users"`
}

type User struct {
	Id        string `json:"id"`
	Email     string `json:"email,omitempty"`
	Telephone string `json:"telephone,omitempty"`
	Language  string `json:"language,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	PhotoUrl  string `json:"photoURL,omitempty"`
	Country   string `json:"country,omitempty" `
}

var id string
var firstName string
var lastName string
var email string
var mobile string
var country string
var language string
var timezone string
var photoUrl string
var filePath string

// upsertUserCmd represents the upsertUser command
var upsertUserCmd = &cobra.Command{
	Use:   "upsertUser",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if filePath != "" {
			fmt.Println("The configuration will be only read from file ", filePath)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upsertUser called")

		if filePath == "" {

			user := new(User)
			user.Email = email
			user.Telephone = mobile
			user.Country = country
			user.FirstName = firstName
			user.LastName = lastName
			user.Id = id
			user.PhotoUrl = photoUrl

			dto := new(UpsertUsers)
			dto.Users = []User{*user}

			httpRequest := upsertUsers(dto)
			runRequest(httpRequest)

		} else {
			confFile, err := ioutil.ReadFile(filePath)
			cobra.CheckErr(err)

			var users *UpsertUsers
			err = json.Unmarshal(confFile, &users)
			cobra.CheckErr(err)

			var dto *UpsertUsers

			for _, user := range users.Users {
				dto = new(UpsertUsers)
				dto.Users = []User{user}
				req:= upsertUsers(dto)
				runRequest(req)
			}

			// Batch doesn't work
			/* if amountUsers > 2 {
				cursor := 0
				var dto *UpsertUsers
				for over := true; over; {
					
					next := cursor + 2
					if next > amountUsers {
						next = amountUsers
						over = false
					}

					dto = new(UpsertUsers)
					dto.Users = users.Users[cursor:next-1]
					req := upsertUsers(dto)
					runRequest(req)

					cursor = next
				}

			} else {
				req := upsertUsers(users)
				runRequest(req)
			} */
		}

	},
}

func upsertUsers(dto *UpsertUsers) *http.Request {

	json_data, err := json.Marshal(dto)
	cobra.CheckErr(err)

	httpRequest, err := http.NewRequest("POST", "https://api.notifuse.com/users.upsert", bytes.NewBuffer(json_data))
	cobra.CheckErr(err)

	apiKey := viper.GetString("NOTIFUSE_APIKEY")
	if apiKey == "" {
		cobra.CompError("no API key is supplied")
	}

	httpRequest.Header.Add("Authorization", "Bearer "+apiKey)
	req, _ := httputil.DumpRequest(httpRequest, true)
	fmt.Println(string(req))

	return httpRequest

}

func runRequest(httpRequest *http.Request) {
	client := new(http.Client)

	response, err := client.Do(httpRequest)

	cobra.CheckErr(err)
	defer response.Body.Close()

	res, err := httputil.DumpResponse(response, true)
	cobra.CheckErr(err)
	fmt.Println(string(res))
}

func init() {
	rootCmd.AddCommand(upsertUserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upsertUserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upsertUserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	upsertUserCmd.Flags().StringVar(&id, "id", "", "User id")
	upsertUserCmd.Flags().StringVar(&firstName, "first-name", "", "User firstname")
	upsertUserCmd.Flags().StringVar(&lastName, "last-name", "", "User last name")
	upsertUserCmd.Flags().StringVar(&email, "email", "", "User email")
	upsertUserCmd.Flags().StringVar(&timezone, "tz", "", "User timezone")
	upsertUserCmd.Flags().StringVar(&language, "lang", "", "User language")
	upsertUserCmd.Flags().StringVar(&country, "country", "", "User country")
	upsertUserCmd.Flags().StringVar(&mobile, "phone", "", "User mobile(international format)")
	upsertUserCmd.Flags().StringVar(&photoUrl, "profile-picture", "", "User profile picture")
	upsertUserCmd.Flags().StringVar(&filePath, "from-file", "", "Upsert users from file")
}
