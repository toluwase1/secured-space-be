package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pusher/pusher-http-go"
	"io/ioutil"
	"log"
	"os"
)

type Chat struct {
	Client	pusher.Client
}

type user struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
	ApartmentID	string	`json:"apartment_id"`
}

func NewChat() *Chat {
    client := pusher.Client{
        AppID:   os.Getenv("PUSHER_APP_ID"),
        Key:     os.Getenv("PUSHER_APP_KEY"),
        Secret:  os.Getenv("PUSHER_APP_SECRET"),
        Cluster: os.Getenv("PUSHER_APP_CLUSTER"),
        Secure:  true,
    }

    return &Chat{
        Client: client,
    }
}

func (s *Server) SendMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		message := struct {
			Message string `json:"message"`
			Channel string  `json:"channel"`
		}{}
		err := c.ShouldBindJSON(&message)
		if err != nil {
			log.Printf("Error: %v", err.Error())
		}
		client := NewChat()
		err = client.Client.Trigger(message.Channel, "message", message)
	    if err != nil {
			log.Printf("Error: %v", err.Error())
		}
		c.JSON(200, message)
	}
}

func (s *Server) registerNewUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)

		if err != nil {
			panic(err)
		}

		var newUser user
		err = json.Unmarshal(body, &newUser)

		if err != nil {
			panic(err)
		}
		client := NewChat()
		client.Client.Trigger("update", "new-user", newUser)

		c.JSON(200, newUser)
	}

}
func (s *Server) CreateChat() gin.HandlerFunc{
	return func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)

		if err != nil {
			panic(err)
		}

		chat := struct {
            ChannelName string `json:"channel_name"`
			InitiatedBy string `json:"initiated_by"`
			ChatWith string `json:"chat_with"`
        }{}
		err = json.Unmarshal(body, &chat)

		if err != nil {
			panic(err)
		}
		client := NewChat()
		client.Client.Trigger(fmt.Sprintf("private-notification-%s", chat.InitiatedBy), "one-to-one-chat-request", chat)
		client.Client.Trigger(fmt.Sprintf("private-notification-%s", chat.ChatWith), "one-to-one-chat-request", chat)

		c.JSON(200, chat)
	}

}
func (s *Server) pusherAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := NewChat()
		params, _ := ioutil.ReadAll(c.Request.Body)

		response, err := client.Client.AuthenticatePrivateChannel(params)

		if err != nil {
			panic(err)
		}

		c.JSON(200, response)
	}

}



func (s *Server) SendNewMessage() gin.HandlerFunc {
    return func(c *gin.Context) {
		payload := struct {
			Message string `json:"message"`
			Username string  `json:"username"`
		}{}
		err := c.ShouldBindJSON(&payload)
		if err != nil {
			log.Printf("Error: %v", err.Error())
		}
		client := NewChat()
		err = client.Client.Trigger("chat", "message", payload)
		if err != nil {
			log.Printf("Error: %v", err.Error())
		}
		c.JSON(200, payload)
    }
}