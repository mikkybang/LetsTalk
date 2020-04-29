(%define "sidebar" %)
<template>
  <v-container style="height: 100vh; max-width: 300px" fluid>
    <v-row justify="center" align="center">
      <v-col cols="12">
        <v-img src="./assets/unilag.svg" align="left" contain height="100"></v-img>
      </v-col>

      <v-col md="auto">
        <v-menu bottom offset-y>
          <template v-slot:activator="{ getUsers,on }">
            <v-text-field
              v-on="on"
              v-model="searchText"
              rounded
              filled
              clearable
              placeholder="search for contacts"
            ></v-text-field>
          </template>
          <v-list>
            <v-list-item v-for="(user,i) in users " :key="i" @click="() => {}">
              <v-list-item-title>{{ user }}</v-list-item-title>
            </v-list-item>
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

    <v-flex fluid style="height: 60vh; max-width: 300px" class="overflow-y-auto">
      <v-list tile dense three-line>
        <v-list-item-group v-model="item" color="black">
          <v-list-item v-for="i in 10" :key="i">
            <v-list-item-avatar>
              <v-icon large>mdi-account-circle</v-icon>
            </v-list-item-avatar>

            <v-list-item-content>
              <v-list-item-title>matric@live.unilag.edu.ng</v-list-item-title>
              <v-list-item-subtitle>Text Here</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-flex>
  </v-container>
</template>
(%end%)

(%define "sidebarData"%)
  data: () => ({
    showSearch: false,
    users: [],
    searchText: '',
    id: "(%.Email%)",
    uuid: "(%.UUID%)",
    item: 0,
  }),

  method: {
    getUsers: function () {
      if (this.searchText.length > 6){
        this.showSearch = true
        var url = location.protocol + "//"+ document.location.host +"/search/" + this.id + "/" + this.uuid + "/" + this.searchText
        console.log(url)
        axios.get(url)
          .then((response) => {
            var obj = JSON.parse(JSON.stringify(response.data));
            console.log(obj);
            console.log(obj.Users);
            this.users = obj.Users;
        });
      } else {
            this.users = []
      }
      return this.users
    }
  }
(%end%)