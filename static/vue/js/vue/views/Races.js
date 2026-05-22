import API from '../apicaller.js?version=102'
import RaceEditor from './components/race-editor.js?version=101'

export default {
  components: {
    RaceEditor,
  },
  data() {
    return {
      races: [],
      loadingData: false,
      errorMsg: '',
    }
  },
  created() {
    this.$store.commit('clearGenAll')
    this.$store.commit('tokenFromCache')
    const token = this.$store.state.admin.token
    if (!token) {
      console.log('token is missed')
      this.$router.push('/login')
      return
    }
    this.fetchRaces()
  },
  computed: {
  },
  methods: {
    fetchRaces() {
      this.loadingData = true
      this.errorMsg = ''
      API.GetRaces(this, (races) => {
        this.races = races || []
        this.loadingData = false
      })
    },
    onSaveRace(race) {
      API.SaveRace(this, race, (races) => {
        this.races = races || []
      })
    },
    onDeleteRace(id) {
      API.DeleteRace(this, id, (races) => {
        this.races = races || []
      })
    },
  },
  template: `
  <v-card>
    <v-card-title class="subheading grey--text">Races</v-card-title>
    <v-card-subtitle>Race Management</v-card-subtitle>
    <v-divider></v-divider>
    <v-alert v-if="errorMsg" type="error" dismissible>{{ errorMsg }}</v-alert>
    <RaceEditor
      :races="races"
      :loadingData="loadingData"
      @save-race="onSaveRace"
      @delete-race="onDeleteRace"
    />
  </v-card>
`
}
