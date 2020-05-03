(%define "chattingComponent"%)
<template>
  <v-container style="height: 100vh;" fluid>
    <v-row no-gutters>
      <v-col cols="12">
        <v-row no-gutters>
          <v-col>
            <v-col md="12"></v-col>
            <h2>
              <b>Name Here (Matric@live.unilag.edu.ng)</b>
            </h2>
          </v-col>

          <v-col md="auto" align="right">
            <v-btn fab outlined depressed>
              <v-icon>mdi-account-multiple-plus-outline</v-icon>
            </v-btn>
            <v-btn fab outlined depressed>
              <v-icon>mdi-video-outline</v-icon>
            </v-btn>
            <v-btn fab outlined depressed>
              <v-icon>mdi-phone</v-icon>
            </v-btn>
          </v-col>
        </v-row>
        <v-divider></v-divider>
        <v-col cols="12"></v-col>
      </v-col>

      <v-container fluid style="height: 80vh; " class="overflow-auto">
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
            <v-col cols="12" align="right">
              <v-card outlined class="d-inline-block mx-auto">
                <v-card-title class="text--secondary">
                  <h6>11:44PM 12/04/2020</h6>
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
                <v-card-text class="flex">
                  <span>Hello there my name is michael</span>
                </v-card-text>
              </v-card>
            </v-col>

            <v-col cols="12" align="left">
              <v-card outlined class="d-inline-block mx-auto">
                <v-card-title class="text--secondary">
                  <h6>Matric 11:44PM 12/04/2020</h6>
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
                <v-card-text class="flex">
                  <span>Hi there</span>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </div>
      </v-container>

      <v-row no-gutters>
        <v-col md="auto" align="right">
          <v-btn fab outlined depressed>
            <v-icon>mdi-microphone-settings</v-icon>
          </v-btn>

          <v-btn fab outlined depressed>
            <v-icon>mdi-paperclip</v-icon>
          </v-btn>
        </v-col>

        <v-textarea rows="1" filled rounded clearable></v-textarea>
      </v-row>
    </v-row>
  </v-container>
</template>
(%end%)