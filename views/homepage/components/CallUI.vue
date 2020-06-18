(%define "callUI"%)
<template>
  <v-container
    fill-height
    align="center"
    justify="center"
    style="min-width: 100%; min-height: 100%; background-color: black;"
    @click="show_video_options=!show_video_options"
  >
    <!-- shown if video is disabled. -->
    <div v-if="!video_on" align="center" justify="center" class="mx-auto my-auto">
      <v-row>
        <v-col cols="12">
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
        </v-col>
        <v-col cols="12">
          <span class="white--text headline">{{caller_names}}</span>
        </v-col>
      </v-row>
    </div>

    <video
      v-else
      style="min-width: 100%; min-height: 100%;right: 0;bottom: 0;position: fixed; "
      id="videoID"
      muted="muted"
      autoplay
      playsinline
    ></video>

    <v-expand-transition>
      <v-row
        justify="center"
        v-if="show_video_options"
        style="bottom: 1rem; padding: 20px; left: 50%; margin-left: -150px; position: absolute;"
      >
        <v-btn
          class="ml-2"
          small
          fab
          dark
          @click="desktop_share()"
          :color="is_desktop_shared? 'blue' : ''"
        >
          <v-icon>mdi-desktop-mac</v-icon>
        </v-btn>

        <v-btn class="ml-2" fab small dark @click="change_video_status()">
          <v-icon>{{video_icon}}</v-icon>
        </v-btn>

        <v-btn class="ml-2" color="red" small fab dark @click="end_call()">
          <v-icon>mdi-phone-hangup</v-icon>
        </v-btn>

        <v-btn class="ml-2" small fab dark @click="change_audio_status()">
          <v-icon>{{audio_icon}}</v-icon>
        </v-btn>

        <v-btn
          :color="full_screen ? 'blue' : ''"
          class="ml-2"
          small
          fab
          dark
          @click="request_full_screen()"
        >
          <v-icon>mdi-fullscreen</v-icon>
        </v-btn>
      </v-row>
    </v-expand-transition>
  </v-container>
</template>
(%end%)