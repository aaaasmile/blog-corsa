import API from '../apicaller.js?version=101'


export default {
    data() {
        return {
            Password: '',
            User: ''
        }
    },
    created() {
    },
    Computed: {
        ...Vuex.mapState({
            ResLogs: state => {
                return state.gen.reslog
            }
        })
    },
    methods: {
        doLogin() {
            console.log('submit login')
            let para = {
                user: this.User,
                password: this.Password,
            }
            API.DoLogin(this, para)
        }
    },
    template: `
  <v-card>
    <v-card-title class="subheading grey--text">Login</v-card-title>

    <v-divider></v-divider>
    <v-tooltip bottom>
      <template v-slot:activator="{ on }">
        <v-text-field
          v-model="User"
          label="User"
        ></v-text-field>
        <v-text-field
          v-model="Password"
          label="Password"
        ></v-text-field>

        <v-btn icon @click="doLogin" v-on="on">
          <v-icon>mdi-file</v-icon>
        </v-btn>
      </template>
      <span>Login</span>
    </v-tooltip>
  </v-card>
`
}