import RaceItem from './race-item.js?version=100'
import RaceLapEditor from './race_lap_editor.js?version=100'

export default {
  components: {
    RaceItem,
    RaceLapEditor
  },
  props: {
    races: {
      type: Array,
      default: () => []
    },
    loadingData: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      search: '',
      showRaceItem: false,
      showLapEditor: false,
      selectedRaceForLaps: null,
      editedIndex: -1,
      editedRace: {
        Name: '',
        Title: '',
        meter_length: '',
        ascending_meter: 0,
        descending_meter: 0,
        rank_global: 0,
        rank_gender: 0,
        rank_class: 0,
        class_name: '',
        race_start_datetime: '',
        sport_type_id: 0,
        result_time: '',
        pace_kmh: 0,
        pace_minkm: '',
        comment: '',
        race_subtype_id: 0,
        km_length: 0,
        loop_number: 0,
        loop_length: 0,
      },
      headers: [
        { text: 'ID', value: 'Id' },
        { text: 'Name', value: 'Name' },
        { text: 'Title', value: 'Title' },
        { text: 'Distance (Km)', value: 'meter_length' },
        { text: 'Date', value: 'race_start_datetime' },
        { text: 'Actions', value: 'actions', sortable: false },
      ],
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
      this.editedIndex = -1
      this.editedRace = { Id: 0, Name: '', Title: '', meter_length: '', ascending_meter: 0, descending_meter: 0, rank_global: 0, rank_gender: 0, rank_class: 0, class_name: '', race_start_datetime: '', sport_type_id: 0, result_time: '', pace_kmh: 0, pace_minkm: '', comment: '', race_subtype_id: 0, km_length: 0, loop_number: 0, loop_length: 0 }
      this.showRaceItem = true
      this.showLapEditor = false
    },
    editLaps(item) {
      this.selectedRaceForLaps = item
      this.showLapEditor = true
      this.showRaceItem = false
    },
    closeLapEditor() {
      this.showLapEditor = false
      this.selectedRaceForLaps = null
    },
    editRace(item) {
      this.editedIndex = this.races.indexOf(item)
      this.editedRace = Object.assign({}, item)
      this.showRaceItem = true
      this.showLapEditor = false
    },
    onSaveRace(race) {
      this.$emit('save-race', race)
      this.showRaceItem = false
      this.editedIndex = -1
    },
    onCancelRace() {
      this.showRaceItem = false
      this.editedIndex = -1
    },
    deleteRace(item) {
      this.raceToDelete = item
      this.dialogDelete = true
    },
    formatDate(isoStr) {
      if (!isoStr) return ''
      const d = new Date(isoStr)
      const dd = String(d.getDate()).padStart(2, '0')
      const mm = String(d.getMonth() + 1).padStart(2, '0')
      const yyyy = d.getFullYear()
      return `${dd}-${mm}-${yyyy}`
    },
    confirmDelete() {
      if (this.raceToDelete && this.raceToDelete.Id) {
        this.$emit('delete-race', this.raceToDelete.Id)
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
        <template v-slot:item.meter_length="{ item }">
          <span>{{ item.meter_length ? Number(item.meter_length).toLocaleString() : '' }}</span>
        </template>
        <template v-slot:item.race_start_datetime="{ item }">
          <span>{{ item.race_start_datetime ? formatDate(item.race_start_datetime) : '' }}</span>
        </template>
        <template v-slot:item.actions="{ item }">
          <v-icon small class="mr-2" color="blue" @click="editLaps(item)">mdi-format-list-numbered</v-icon>
          <v-icon small class="mr-2" @click="editRace(item)">mdi-pencil</v-icon>
          <v-icon small color="red" @click="deleteRace(item)">mdi-delete</v-icon>
        </template>
      </v-data-table>
    </v-card>
    <div v-if="showRaceItem" class="mb-12">
      <RaceItem
        :value="editedRace"
        @save="onSaveRace"
        @cancel="onCancelRace"
      />
    </div>
    <div v-if="showLapEditor" class="mb-12">
      <RaceLapEditor
        :race="selectedRaceForLaps"
        @close="closeLapEditor"
      />
    </div>
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
