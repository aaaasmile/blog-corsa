import API from '../apicaller.js?version=101'
import CmtToModerate from './components/cmt-to-mod.js?version=100'

export default {
  components: {
    CmtToModerate,
  },
  data() {
    return {
    }
  },
  created() {
    this.$store.commit('clearGenAll')
    this.$store.commit('tokenFromCache')
    const token = this.$store.state.admin.token
    if (!token){
      console.log('token is missed')
      this.$router.push('/login')
    }
  },
  computed: {
    ...Vuex.mapState({
      ResLogs: state => {
        return state.gen.reslog
      }
    })
  },
  methods: {
    approveCmt() {
      console.log('approve comment')
      let para = {  }
      API.ApproveCmt(this, para)
    }
  },
  template: `
  <v-card>
    <v-card-title class="subheading grey--text">Dashboard</v-card-title>
    <v-card-subtitle>Blog Comments</v-card-subtitle>
    <v-divider></v-divider>
    <v-sheet border="md" class="pa-6 text-white mx-auto" max-width="800">
      <h4 class="text-h5 font-weight-bold mb-4">Console</h4>
      <v-list dense>
        <v-list-item v-for="item in ResLogs" :key="item.key">
          <v-list-item-content>
            <div>{{ item.text }}</div>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-sheet>
  </v-card>
`
}