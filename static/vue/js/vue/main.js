

import Navbar from './views/components/Navbar.js?version=53'
import store from './store/index.js?version=112'
import routes from './routes.js?version=109'


export const app = new Vue({
  el: '#app',
  router: new VueRouter({ routes }),
  components: { Navbar },
  vuetify: new Vuetify(),
  store,
  data() {
    return {
      Buildnr: "",
    }
  },
  computed: {
    ...Vuex.mapState({
      ServiceServer: state => {
        return state.cmt.server_name
      }
    })
  },
  created() {
    this.Buildnr = window.myapp.buildnr
    this.$store.commit('serviceServer', window.myapp.servername)
  },
  methods: {

  },
  template: `
  <v-app class="grey lighten-4">
    <Navbar />
    <v-content class="mx-4 mb-4">
      <router-view></router-view>
    </v-content>
    <v-footer absolute>
      <v-col class="text-center caption" cols="12">
        {{ new Date().getFullYear() }} —
        <span>Buildnr: {{Buildnr}} - Server: {{ServiceServer}}</span>
      </v-col>
    </v-footer>
  </v-app>
`
})

console.log('Main is here!')