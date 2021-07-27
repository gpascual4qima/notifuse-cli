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
	"net/http"
	"net/http/httputil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userId string
var notificationId string

type SendMessageDTO struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	UserId         string            `json:"userId"`
	NotificationId string            `json:"notificationId"`
	UserPhotoURL   string            `json:"userPhotoURL"`
	Data           map[string]string `json:"data"`
}

// sendMessageCmd represents the sendMessage command
var sendMessageCmd = &cobra.Command{
	Use:   "sendMessage",
	Short: "Send message to user",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sendMessage called")

		client := new(http.Client)

		message := new(Message)
		message.Data = make(map[string]string)
		message.NotificationId = notificationId
		message.UserId = userId
		photoUrl, err := cmd.Flags().GetString("photo-url")
		if err == nil {
			message.UserPhotoURL = photoUrl
		}

		dto := new(SendMessageDTO)
		dto.Messages = []Message{*message}

		json_data, err := json.Marshal(dto)
		cobra.CheckErr(err)
		httpRequest, err := http.NewRequest("POST", "https://api.notifuse.com/messages.send", bytes.NewBuffer(json_data))
		cobra.CheckErr(err)

		apiKey := viper.GetString("NOTIFUSE_APIKEY")
		if apiKey == "" {
			cobra.CompError("no API key is supplied")
		}
		httpRequest.Header.Add("Authorization", "Bearer "+apiKey)
		req, err := httputil.DumpRequest(httpRequest, true)
		cobra.CheckErr(err)
		fmt.Println(string(req))

		response, err := client.Do(httpRequest)
		cobra.CheckErr(err)
		res, err := httputil.DumpResponse(response, true)
		cobra.CheckErr(err)
		fmt.Println(string(res))

	},
}

func init() {
	rootCmd.AddCommand(sendMessageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendMessageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendMessageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sendMessageCmd.Flags().StringVar(&userId, "user-id", "", "User ID")
	sendMessageCmd.MarkFlagRequired("user-id")
	sendMessageCmd.Flags().StringVar(&notificationId, "notification-id", "", "Notification ID")
	sendMessageCmd.MarkFlagRequired("notification-id")
	sendMessageCmd.Flags().String("photo-url", "", "User profile picture")
}
