import API from '../../apicaller.js?version=101'
import RaceItem from './race-item.js?version=100'

export default {
  components: {
    RaceItem,
  },
  data() {
    return {
      loadingData: false,
      search: '',
      showRaceItem: false,
      newRace: {
        Name: '',
        Title: '',
        Distance: '',
        Date: '',
      },
      headers: [
        { text: 'ID', value: 'Id' },
        { text: 'Name', value: 'Name' },
        { text: 'Title', value: 'Title' },
        { text: 'Distance', value: 'Distance' },
        { text: 'Date', value: 'Date' },
        { text: 'Actions', value: 'actions', sortable: false },
      ],
      races: [],
      dialogDelete: false,
      raceToDelete: null,
    }
  },
  created() {
    console.log('race-editor created')
  },
  computed: {
  },
  methods: {
    addRace() {
      this.newRace = { Name: '', Title: '', Distance: '', Date: '' }
      this.showRaceItem = true
    },
    onSaveRace(race) {
      race.Id = this.races.length + 1
      this.races.push(race)
      this.showRaceItem = false
      console.log('Race saved', race)
    },
    onCancelRace() {
      this.showRaceItem = false
    },
    deleteRace(item) {
      this.raceToDelete = item
      this.dialogDelete = true
    },
    confirmDelete() {
      const index = this.races.indexOf(this.raceToDelete)
      if (index > -1) {
        this.races.splice(index, 1)
        console.log('Race deleted', this.raceToDelete)
      }
      this.dialogDelete = false
      this.raceToDelete = null
    },
  },
  template: `
  <v-container>
    <v-card>
      <v-card-subtitle>Race Editor</v-card-subtitle>

      <v-divider></v-divider>

      <v-card-actions>
        <v-btn color="primary" @click="addRace">
          <v-icon left>mdi-plus</v-icon>
          Add Race
        </v-btn>
      </v-card-actions>

      <v-data-table
        :headers="headers"
        :items="races"
        :loading="loadingData"
        item-key="Id"
        class="elevation-1"
        :search="search"
        :footer-props="{
          showFirstLastPage: true,
          firstIcon: 'mdi-arrow-collapse-left',
          lastIcon: 'mdi-arrow-collapse-right',
          prevIcon: 'mdi-minus',
          nextIcon: 'mdi-plus',
        }"
      >
        <template v-slot:item.actions="{ item }">
          <v-icon small color="red" @click="deleteRace(item)">mdi-delete</v-icon>
        </template>
      </v-data-table>
    </v-card>
    <RaceItem
      v-if="showRaceItem"
      :value="newRace"
      @save="onSaveRace"
      @cancel="onCancelRace"
    />
    <v-dialog v-model="dialogDelete" persistent max-width="350">
      <v-card>
        <v-card-title class="headline">Confirm Delete</v-card-title>
        <v-card-text>
          Are you sure you want to delete this race?
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" text @click="confirmDelete">OK</v-btn>
          <v-btn color="green darken-1" text @click="dialogDelete = false">Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
`
}
