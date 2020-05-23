package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/metaclips/LetsTalk/values"
)

type messageBytes []byte

func (msg messageBytes) handleCreateNewRoom() {
	var newRoom NewRoomRequest
	if err := json.Unmarshal(msg, &newRoom); err != nil {
		log.Println("Could not convert to required New Room Request struct")
		return
	}

	roomID, err := newRoom.CreateNewRoom()
	if err != nil {
		log.Println("Unable to create a new room for user:", newRoom.Email, "err:", err.Error())
		return
	}

	// Broadcast a joined message.
	userJoinedMessage := Joined{
		RoomID:      roomID,
		Email:       newRoom.Email,
		RoomName:    newRoom.RoomName,
		MessageType: "UserJoinedRoom",
	}

	jsonByte, err := json.Marshal(userJoinedMessage)
	if err != nil {
		log.Println("Could not marshal to jsonByte while creating room", err.Error())
		return
	}

	m := WSMessage{jsonByte, newRoom.Email}
	HubConstruct.Broadcast <- m
}

func (msg messageBytes) handleRequestUserToJoinRoom() {
	var request JoinRequest
	if err := json.Unmarshal(msg, &request); err != nil {
		log.Println("Could not convert to required Joined Request struct")
		return
	}

	for _, user := range request.Users {
		roomRegisteredUser, err := request.RequestUserToJoinRoom(user)
		if err != nil {
			log.Println("Error while requesting to room", err)
			continue
		}

		data := map[string]interface{}{
			"requesterID":   request.RequestingUserID,
			"requesterName": request.RequestingUserName,
			"userRequested": user,
			"roomID":        request.RoomID,
			"roomName":      request.RoomName,
			"msgType":       "RequestUsersToJoinRoom",
		}

		jsonContent, err := json.Marshal(data)
		if err != nil {
			log.Println("could not marshal to RequestUsersToJoinRoom, err:", err)
			continue
		}

		// Send back RequestUsersToJoinRoom signal to everyone registered in room.
		for _, roomRegisteredUser := range roomRegisteredUser {
			m := WSMessage{jsonContent, roomRegisteredUser}
			HubConstruct.Broadcast <- m
		}

		m := WSMessage{jsonContent, user}
		HubConstruct.Broadcast <- m
	}
}

func (msg messageBytes) handleUserAcceptRoomRequest(joiner string) {
	var roomRequest Joined
	if err := json.Unmarshal(msg, &roomRequest); err != nil {
		log.Println("Could not convert to required Join Room Request struct")
		return
	}

	if roomRequest.Email != joiner {
		return
	}

	users, err := roomRequest.AcceptRoomRequest()
	if err != nil {
		log.Println("could not join room", err)
		return
	}

	for _, user := range users {
		m := WSMessage{msg, user}
		HubConstruct.Broadcast <- m
	}
}

func handleRequestAllMessages(roomID, requester string) {
	messages, roomName, err := GetAllMessageInRoom(roomID)
	if err != nil {
		log.Println("could not get all messages in room, err:", err)
		return
	}

	data := map[string]interface{}{
		"messages": messages,
		"msgType":  "RequestAllMessages",
		"roomName": roomName,
		"roomID":   roomID,
	}

	jsonContent, err := json.Marshal(data)
	if err != nil {
		log.Println("could not marshal images, err:", err)
		return
	}

	m := WSMessage{jsonContent, requester}
	HubConstruct.Broadcast <- m
}

func (msg messageBytes) handleNewMessage(requester string) {
	var newMessage Message
	if err := json.Unmarshal(msg, &newMessage); err != nil {
		log.Println("Could not convert to required New Message struct")
		return
	}

	if requester != newMessage.UserID {
		return
	}
	newMessage.Time = time.Now().Format(values.TimeLayout)
	// Message is sent back to all users including sender.
	registeredUsers, err := newMessage.SaveMessageContent()
	if err != nil {
		log.Println("Error saving msg to db", err, requester)
		return
	}

	jsonContent, err := json.Marshal(newMessage)
	if err != nil {
		log.Println("Error converted message to json content", err)
		return
	}

	for _, registeredUser := range registeredUsers {
		m := WSMessage{jsonContent, registeredUser}
		HubConstruct.Broadcast <- m
	}
}
