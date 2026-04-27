import RaceLapItem from './race_lap_item.js?version=100'

export default {
  components: {
    RaceLapItem
  },
  props: {
    race: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      search: '',
      showLapItem: false,
      editedIndex: -1,
      editedLap: {},
      laps: [],
      dialogDelete: false,
      lapToDelete: null,
    }
  },
  computed: {
    isBackyard() {
      return this.race && this.race.race_subtype_id === 8
    },
    headers() {
      if (this.isBackyard) {
        return [
          { text: 'Lap Number', value: 'lap_number' },
          { text: 'Start Time', value: 'start_time' },
          { text: 'Checkpoint', value: 'checkpoint' },
          { text: 'Backyard Time (mm:ss)', value: 'backyard_time' },
          { text: 'Actions', value: 'actions', sortable: false },
        ]
      }
      return [
        { text: 'ID', value: 'id' },
        { text: 'Lap Number', value: 'lap_number' },
        { text: 'Lap Time', value: 'lap_time' },
        { text: 'Lap Meter', value: 'lap_meter' },
        { text: 'Tot Race (min/km)', value: 'tot_race_minkm' },
        { text: 'Tot Race Time', value: 'tot_race_time' },
        { text: 'Tot KM Race', value: 'tot_km_race' },
        { text: 'Actions', value: 'actions', sortable: false },
      ]
    },
    defaultLap() {
      if (this.isBackyard) {
        return {
          id: 0,
          lap_number: 0,
          start_time: '',
          backyard_time: '',
          checkpoint: '',
        }
      }
      return {
        id: 0,
        lap_number: 0,
        lap_time: '',
        lap_meter: 0,
        tot_race_minkm: '',
        tot_race_time: '',
        tot_km_race: 0,
      }
    },
  },
  methods: {
    close() {
      this.$emit('close')
    },
    addLap() {
      this.editedIndex = -1
      this.editedLap = Object.assign({}, this.defaultLap)
      this.showLapItem = true
    },
    editLap(item) {
      this.editedIndex = this.laps.indexOf(item)
      this.editedLap = Object.assign({}, item)
      this.showLapItem = true
    },
    onSaveLap(lap) {
      if (this.editedIndex > -1) {
        Object.assign(this.laps[this.editedIndex], lap)
        console.log('Lap updated', lap)
      } else {
        lap.id = this.laps.length + 1
        this.laps.push(lap)
        console.log('Lap saved', lap)
      }
      this.showLapItem = false
      this.editedIndex = -1
    },
    onCancelLap() {
      this.showLapItem = false
      this.editedIndex = -1
    },
    deleteLap(item) {
      this.lapToDelete = item
      this.dialogDelete = true
    },
    confirmDelete() {
      const index = this.laps.indexOf(this.lapToDelete)
      if (index > -1) {
        this.laps.splice(index, 1)
        console.log('Lap deleted', this.lapToDelete)
      }
      this.dialogDelete = false
      this.lapToDelete = null
    },
  },
  template: `
  <v-card class="mt-4">
    <v-card-title>
      Laps for Race: {{ race.Name || race.Title }}
      <v-spacer></v-spacer>
      <v-btn icon @click="close">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-text>
      <v-btn color="primary" class="mb-4" @click="addLap">
        <v-icon left>mdi-plus</v-icon>
        Add Lap
      </v-btn>

      <v-data-table
        :headers="headers"
        :items="laps"
        item-key="id"
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
        <template v-slot:top>
          <v-toolbar flat>
            <v-text-field
              v-model="search"
              append-icon="mdi-magnify"
              label="Search"
              single-line
              hide-details
            ></v-text-field>
          </v-toolbar>
        </template>
        <template v-slot:item.actions="{ item }">
          <v-icon small class="mr-2" @click="editLap(item)">mdi-pencil</v-icon>
          <v-icon small color="red" @click="deleteLap(item)">mdi-delete</v-icon>
        </template>
      </v-data-table>
    </v-card-text>

    <div v-if="showLapItem" class="pa-4">
      <RaceLapItem
        :value="editedLap"
        :race-subtype-id="race.race_subtype_id"
        @save="onSaveLap"
        @cancel="onCancelLap"
      />
    </div>

    <v-dialog v-model="dialogDelete" persistent max-width="350">
      <v-card>
        <v-card-title class="headline">Confirm Delete</v-card-title>
        <v-card-text>
          Are you sure you want to delete this lap?
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" text @click="confirmDelete">OK</v-btn>
          <v-btn color="green darken-1" text @click="dialogDelete = false">Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
  `
}
