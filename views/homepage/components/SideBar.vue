(%define "sidebar" %)
<template>
  <v-container>
    <v-row justify="center" align="center">
      <v-col cols="12">
        <v-img src="./assets/unilag.svg" align="left" contain height="100"></v-img>
      </v-col>

      <v-col cols="12">
        <v-menu bottom offset-y>
          <template v-slot:activator="{ on }">
            <v-text-field
              v-on="on"
              v-model="searchText"
              rounded
              filled
              clearable
              placeholder="search for contacts"
              :change="getUsers"
            ></v-text-field>
          </template>

          <v-list>
            <v-flex fluid style="max-height: 300px" class="overflow-y-auto">
              <v-list-item v-for="(user,i) in users " :key="i" :click="createRoom">
                <v-list-item-title>{{ user }}</v-list-item-title>
              </v-list-item>
            </v-flex>
          </v-list>
        </v-menu>
      </v-col>

      <v-row justify="center" align="center">
        <v-col md="auto">
          <v-btn outlined height="50" width="50">
            <v-icon>mdi-phone</v-icon>
          </v-btn>
        </v-col>
        <v-col md="auto" offset-md="1">
          <v-btn outlined height="50" width="50">
            <v-icon>mdi-bell</v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </v-row>

    <v-expansion-panels>
      <v-expansion-panel>
        <v-expansion-panel-header>Chats</v-expansion-panel-header>
        <v-expansion-panel-content>
          <v-flex style="height: 55vh; max-width: 400px" class="overflow-y-auto">
            <v-list tile dense three-line>
              <v-list-item-group color="black">
                <v-list-item v-for="(chat, i) in chats" :key="i">
                  <v-list-item-avatar>
                    <v-icon large>mdi-account-circle</v-icon>
                  </v-list-item-avatar>
                  <v-list-item-content>
                    <v-list-item-title>matric@live.unilag.edu.ng</v-list-item-title>
                    <v-list-item-subtitle>{{chat}}</v-list-item-subtitle>
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
