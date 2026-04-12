import RaceEditor from './components/race-editor.js?version=100'

export default {
  components: {
    RaceEditor,
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
  },
  methods: {
  },
  template: `
  <v-card>
    <v-card-title class="subheading grey--text">Races</v-card-title>
    <v-card-subtitle>Race Management</v-card-subtitle>
    <v-divider></v-divider>
    <RaceEditor />
  </v-card>
`
}
