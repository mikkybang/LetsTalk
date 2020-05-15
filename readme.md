# Lets Talk

A chatting application using Pion webRTC and gorilla websocket for text, video, voice and file transfer.

# Preface

Lets Talk is a web chatting platform proposal for the University of Lagos unilag.edu.ng. Due to the covid-19 pandemic, neccesity of having an online chatting/learning platform between students and lecturers is vivid.

Lets Talk supports multi-room chats between users, desktop sharing (during class sessions), file transfer (over websocket and RTC), video and voice calls.

# Privacy Features

For users of minor on the platform, a machine learning model is to be created to scan minors chats and report duly to the universities authority if any form of sexual harrasment is found. This would not breach privacy as only the model scans text and block faulting user/room. Machine learning isn't perfect so faulting user has the choice to fair hearing from administrators.

Users also can not register to the platform. Registration is done by the administrator who generates login details for students using their students email address.

# v1.0.0 Milestone

- [x] Text chat support using gorilla websocket

- [x] Mobile UI

- [x] Desktop UI

- [ ] Code optimization (Ok)

- [ ] Code optimization (Getting there)

- [ ] Code optimization (Perfect. Relativey 🤪)

- [ ] File transfer over websocket for file < 1gb and over webRTC for > 1gb

- [ ] Desktop screen sharing support

- [ ] Seamless websocket connection

- [ ] Voice and Video call support

- [ ] Low bandwidth consumption using selective video call transfer

- [ ] Test coverage....

- [ ] Add logging system.

- [ ] Admin portal


# Future goals (This are all considered and might not be implemented)

- Http2 support. Investigate more on server push.

- Use github.com/cloudwebrtc/go-protoo for room handling.

- If necessary, remove gorilla websocket to reduce unneccessary overheads.

- Environment variable support.



# Stack

## Backend

 - [Golang][go]
 - [pion/webrtc][pion]
 - [mongodb][mongo]
 - github.com/julienschmidt/httprouter 
 - github.com/gorilla/securecookie
 - github.com/gorilla/websocket

[go]: golang.org

[mongo]: go.mongodb.org/mongo-driver

[pion]: https://github.com/pion/webrtc

See [go.mod](go.mod) for more information

## Frontend

 - Vue
 - Vuetify

# Installation
```bash
git clone github.com/metaclips/LetsTalk
cd LetsTalk
go run .
```

# Configuration

- To update when voice call is implemented

# Browser support

- To update when voice call is implemented