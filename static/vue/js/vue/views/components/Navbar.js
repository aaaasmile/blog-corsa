import Toast from './toast.js?version=106'

export default {
  components: { Toast },
  data() {
    return {
      drawer: false,
      AppTitle: "Blog Admin",
      links: [{ path: '/', icon: 'dashboard', title: 'Dashboard'},
      ],
    }
  },
  template: `
  <nav>
    <v-app-bar dense flat>
      <v-btn text color="grey" @click="drawer = !drawer">
        <v-icon>mdi-menu</v-icon>
      </v-btn>
      <v-toolbar-title class="text-uppercase grey--text">
        <span class="font-weight-light">{{AppTitle}}</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
    </v-app-bar>
    <Toast></Toast>
    <v-navigation-drawer app v-model="drawer">
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="title">{{AppTitle}}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-divider></v-divider>

      <v-list dense nav>
        <v-list-item v-for="item in links" :key="item.title" :to="item.path" link>
          <v-list-item-icon>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
       <v-divider></v-divider>
      <v-container>
        <v-row justify="center">
          <v-col cols="6">
            <v-btn icon text @click.stop="drawer = false"
              ><v-icon>close</v-icon> Close
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-navigation-drawer>
  </nav>
`
}