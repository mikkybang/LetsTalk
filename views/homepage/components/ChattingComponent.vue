(%define "chattingComponent"%)
<template>
  <div>
    <v-row>
      <v-col cols="12">
        <v-container fluid>
          <v-row no-gutters>
            <v-col cols="mx-auto">
              <span>
                <b>{{currentviewedroomname}}</b>
              </span>
            </v-col>
            <v-col cols="auto">
              <v-dialog scrollable v-model="addUserDialog" width="600px">
                <template v-slot:activator="{ on }">
                  <v-btn fab depressed v-on="on">
                    <v-icon>mdi-account-multiple-plus-outline</v-icon>
                  </v-btn>
                </template>

                <v-card>
                  <v-card-title>Add Users</v-card-title>
                  <v-divider></v-divider>
                  <v-card-text style="height: 500px;">
                    <v-row>
                      <v-col cols="12">
                        <v-text-field
                          placeholder="Search Users"
                          rounded
                          v-model="searchText"
                          outlined
                          append-icon="mdi-magnify"
                          @click:append="getUsers"
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12">
                        <v-checkbox
                          v-for="(user,i) in usersFound"
                          :key="i"
                          :label="user"
                          :value="user"
                          v-model="addedUsers"
                        ></v-checkbox>
                      </v-col>
                    </v-row>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="green darken-1" text @click="requestUsersToJoinRoom">Add User(s)</v-btn>
                    <v-btn color="green darken-1" text @click="closeSearchDialog">Close</v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <v-btn fab depressed>
                <v-icon>mdi-phone</v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-container>
        <v-divider></v-divider>
      </v-col>

      <v-col cols="12">
        <v-container class="overflow-y-auto" style="height: 75vh;" fluid>
          <div v-if="currentchatcontentsloaded===true" align="center">
            <v-dialog max-width="300px" persistent v-model="currentchatcontentsloaded">
              <v-card>
                <v-card-text>
                  <div class="text-center" align="center" justify="center">
                    <v-row>
                      <v-col cols="12">
                        <v-progress-circular indeterminate color="green"></v-progress-circular>
                      </v-col>
                      <v-col cols="12">
                        <span>Fetching Content</span>
                      </v-col>
                    </v-row>
                  </div>
                </v-card-text>
              </v-card>
            </v-dialog>
          </div>

          <div v-else>
            <v-row>
              <v-col cols="12" v-for="(chat,i) in currentchatcontent" :key="i">
                <div align="center" justify="center" v-if="chat['type']==='info'">
                  <v-card class="justify-center" outlined>
                    <v-card-text>{{chat['message']}}</v-card-text>
                  </v-card>
                </div>

                <div v-else-if="chat['userID']==='(%.Email%)'" align="right">
                  <v-card outlined max-width="200px" class="d-inline-block mx-auto">
                    <v-card-title class="text--secondary">
                      <h6>{{chat['name']}}</h6>
                      <v-spacer></v-spacer>
                      <v-card-actions>
                        <v-menu absolute bottom>
                          <template v-slot:activator="{ on }">
                            <v-btn icon v-on="on">
                              <v-icon>mdi-chevron-down</v-icon>
                            </v-btn>
                          </template>

                          <v-list>
                            <v-list-item v-for="i in 5" :key="i">
                              <v-list-item-title>{{i}}</v-list-item-title>
                            </v-list-item>
                          </v-list>
                        </v-menu>
                      </v-card-actions>
                    </v-card-title>
                    <v-card-text>
                      <span>{{chat['message']}}</span>
                    </v-card-text>
                    <v-card-subtitle>{{chat['time']}}</v-card-subtitle>
                  </v-card>
                </div>

                <div v-else align="left">
                  <v-card max-width="300px" outlined class="d-inline-block mx-auto">
                    <v-card-title class="text--secondary">
                      <h6>{{chat['name']}}</h6>
                      <v-spacer></v-spacer>
                      <v-card-actions>
                        <v-menu absolute bottom left>
                          <template v-slot:activator="{ on }">
                            <v-btn icon v-on="on">
                              <v-icon>mdi-chevron-down</v-icon>
                            </v-btn>
                          </template>

                          <v-list>
                            <v-list-item v-for="i in 5" :key="i">
                              <v-list-item-title>{{i}}</v-list-item-title>
                            </v-list-item>
                          </v-list>
                        </v-menu>
                      </v-card-actions>
                    </v-card-title>
                    <v-card-text>
                      <span>{{chat['message']}}</span>
                    </v-card-text>
                    <v-card-subtitle>{{chat['time']}}</v-card-subtitle>
                  </v-card>
                </div>
              </v-col>
            </v-row>
          </div>
        </v-container>
      </v-col>

      <v-col cols="12">
        <v-container fluid>
          <v-textarea
            height="10"
            v-model="messageContent"
            prepend-inner-icon="mdi-emoticon"
            prepend-icon="mdi-paperclip"
            append-outer-icon="mdi-send"
            @click:append-outer="sendMessage"
            @keyup.enter.exact="sendMessage"
            solo
            hide-details="auto"
            dense
            auto-grow
            rows="1"
            rounded
            clearable
          ></v-textarea>
        </v-container>
      </v-col>
    </v-row>
  </div>
</template>
(%end%)