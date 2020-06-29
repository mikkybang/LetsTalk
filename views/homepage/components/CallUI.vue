(%define "callUI"%)
<template>
  <v-container
    fill-height
    align="center"
    justify="center"
    style="min-width: 100%; min-height: 100%; background-color: black;"
    @click="showVideoOptions=!showVideoOptions"
  >
    <!-- shown if video is disabled. -->
    <v-row
      v-if="!video_on"
      style="width: 100%; height: 100%;right: 0; top: 0; left:0; bottom: 0; position: fixed; z-index: 1; "
      align="center"
      justify="center"
      class="mx-auto my-auto"
    >
      <svg height="100" width="100">
        <circle cx="50" cy="50" r="40" stroke="black" fill="#56CCF2" />
        <text
          x="50%"
          y="60%"
          text-anchor="middle"
          fill="black"
          font-size="30px"
          font-family="Arial"
        >{{caller_name_initials}}</text>
      </svg>
      <span class="white--text headline">{{caller_names}}</span>
    </v-row>

    <video
      style="width: 100%; height: 100%;top: 0; left:0; right: 0; bottom: 0; position: absolute; object-fit: cover; "
      id="videoID"
      autoplay
    ></video>

    <v-expand-transition>
      <v-row
        justify="center"
        v-if="showVideoOptions"
        style="bottom: 1rem; padding: 20px; left: 50%; margin-left: -150px; position: fixed; z-index: 2;"
      >
        <v-btn
          class="ml-2"
          small
          fab
          dark
          @click="desktopShare()"
          :color="isDesktopShared? 'blue' : ''"
        >
          <v-icon>mdi-desktop-mac</v-icon>
        </v-btn>

        <v-btn class="ml-2" fab small dark @click="changeVideoStatus()">
          <v-icon>{{videoIcon}}</v-icon>
        </v-btn>

        <v-btn class="ml-2" color="red" small fab dark @click="endCall()">
          <v-icon>mdi-phone-hangup</v-icon>
        </v-btn>

        <v-btn class="ml-2" small fab dark @click="changeAudioStatus()">
          <v-icon>{{audioIcon}}</v-icon>
        </v-btn>

        <v-btn
          :color="fullScreen ? 'blue' : ''"
          class="ml-2"
          small
          fab
          dark
          @click="requestFullScreen()"
        >
          <v-icon>mdi-fullscreen</v-icon>
        </v-btn>
      </v-row>
    </v-expand-transition>
  </v-container>
</template>
(%end%)