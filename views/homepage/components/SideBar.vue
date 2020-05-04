(%define "sidebar" %)
<template>
  <v-container>
    <v-row style="height: 30vh; max-width: 400px" justify="center" align="center">
      <v-col cols="12">
        <v-img src="./assets/unilag.svg" align="left" contain height="100"></v-img>
      </v-col>

      <v-col cols="12">
        <v-row justify="center" align="center">
          <v-col cols="auto">
            <v-dialog v-model="createRoomDialog" width="600px">
              <template v-slot:activator="{ on }">
                <v-btn outlined height="50" width="50" v-on="on">
                  <v-icon>mdi-plus</v-icon>
                </v-btn>
              </template>
              <v-card>
                <v-card-text>
                  <v-container fluid>
                    <v-row>
                      <v-col cols="12">
                        <span>Create New Room</span>
                      </v-col>
                      <v-col cols="12">
                        <v-text-field label="Specify Room Name" v-model="createNewRoomName"></v-text-field>
                      </v-col>
                      <v-col cols="12">
                        <v-spacer></v-spacer>
                        <v-btn color="green darken-1" text @click="createRoom">Create Room</v-btn>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-card-text>
              </v-card>
            </v-dialog>
          </v-col>

          <v-col cols="auto">
            <v-badge
              @click.native="openMessageDialog=!openMessageDialog"
              :content="notificationcount"
              :value="notificationcount"
              color="red"
              overlap
            >
              <v-icon large>mdi-email</v-icon>
            </v-badge>

            <v-dialog v-model="openMessageDialog" width="600px">
              <v-card>
                <v-card-title>Messages</v-card-title>
                <v-card-text>
                  <span>You have no messages</span>
                </v-card-text>
              </v-card>
            </v-dialog>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <v-expansion-panels>
      <v-expansion-panel>
        <v-expansion-panel-header>Chats</v-expansion-panel-header>
        <v-expansion-panel-content>
          <v-flex style="height: 60vh; max-width: 400px" class="overflow-y-auto">
            <v-list tile dense three-line>
              <v-list-item-group color="black">
                <v-list-item v-for="(chatID,i) in chats" :key="i" @click="loadChatContent(chatID)">
                  <v-list-item-avatar>
                    <v-icon large>mdi-account-circle</v-icon>
                  </v-list-item-avatar>
                  <v-list-item-content>
                    <v-list-item-title>{{chatspreview[chatID]["roomName"]}}</v-list-item-title>
                    <v-list-item-subtitle>{{chatspreview[chatID]["message"]}}</v-list-item-subtitle>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-flex>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
  </v-container>
</template>
(%end%)
