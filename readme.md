# Lets Talk

A chatting application using Pion webRTC and gorilla websocket for text, video, voice and file transfer. [Test server](https://unilag-letstalk.herokuapp.com)


# Preface

Lets Talk is a web chatting platform proposal for the University of Lagos (unilag.edu.ng). Due to the covid-19 pandemic, neccesity of having an online chatting/learning platform between students and lecturers is vivid.

Lets Talk supports multi-room chats between users, file transfer (over websocket), video and voice calls support with desktop sharing.

Administrator can only register user **tentative**


# v1.0.0 Milestone

- [x] Multiple room support for users using gorilla websocket

- [x] Mobile UI

- [x] Desktop UI

- [x] Seamless websocket connection

- [x] Resumable file transfer

- [x] Seamless Desktop screen sharing support

- [x] Voice and Video call support

- [x] Low bandwidth consumption using selective video call transfer [#31](https://github.com/metaclips/LetsTalk/issues/31)

- [ ] Add logging system.

- [ ] Admin portal


# Dependencies

## Backend

 - [Golang](golang.org)
 - [pion/webrtc](https://github.com/pion/webrtc)
 - [mongodb](go.mongodb.org/mongo-driver)
 - [httprouter](github.com/julienschmidt/httprouter)
 - [gorilla secure cookies](github.com/gorilla/securecookie)
 - [gorilla websocket](github.com/gorilla/websocket)


See [go.mod](go.mod) for more information

## Frontend

 - Vue
 - Vuetify


# Installation
```bash
git clone https://github.com/metaclips/LetsTalk.git
cd LetsTalk

go run .
```


Login credentials are generated by administrators on the administrators portal `/admin`. Default email and password for administrator is `admin@email.com` and `admin` which can later be disable or changed.
Users will login through generated email addresses and a password of the first name used to register  `surname` with the first letter titled e.g. Kaduru


# Configuration

- To update when voice call is implemented


# Browser support

- To update when voice call is implemented

# License

[Apache 2.0](LICENSE)
