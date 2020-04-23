(%define "tabs"%)
<template>
  <v-card>
    <v-row>
      <v-col cols="12">
        <v-img class="white--text align-end" width="400" height="400" src="/assets/unilag.svg" />
      </v-col>

      <v-col cols="12">
        <div align="center">
          <v-row>
            <v-spacer></v-spacer>
            <v-card-title>Welcome Admin</v-card-title>
            <v-spacer></v-spacer>
          </v-row>
        </div>
      </v-col>

      <v-col cols="12">
        <div>
          <v-tabs grow v-model="tabs" dark>
            <v-tab>ADD User(s)</v-tab>
            <v-tab-item>
              <v-card flat>(%template "addUsers" .%)</v-card>
            </v-tab-item>

            <v-tab>Block User</v-tab>
            <v-tab-item>
              <v-card flat>
                <v-card flat>(%template "blockUser" .%)</v-card>
              </v-card>
            </v-tab-item>

            <v-tab>Scan Messages</v-tab>
            <v-tab-item>
              <v-card flat>(%template "messageScanner" .%)</v-card>
            </v-tab-item>
          </v-tabs>
        </div>
      </v-col>
    </v-row>
  </v-card>
</template>
(%end%)

(%define "tabsData"%)
    data(){
        return{
           
        }
    },

    methods:{
       
    }

(%end%)