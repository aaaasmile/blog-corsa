export default {
  props: {
    race: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      search: '',
      headers: [
        { text: 'ID', value: 'id' },
        { text: 'Lap Number', value: 'lap_number' },
        { text: 'Lap Time', value: 'lap_time' },
        { text: 'Lap Meter', value: 'lap_meter' },
        { text: 'Tot Race (min/km)', value: 'tot_race_minkm' },
        { text: 'Tot Race Time', value: 'tot_race_time' },
        { text: 'Tot KM Race', value: 'tot_km_race' },
      ],
      laps: []
    }
  },
  methods: {
    close() {
      this.$emit('close')
    },
    addLap() {
      const newLap = {
        id: this.laps.length + 1,
        lap_number: this.laps.length + 1,
        lap_time: '',
        lap_meter: 0,
        tot_race_minkm: 0,
        tot_race_time: '',
        tot_km_race: 0
      }
      this.laps.push(newLap)
    }
  },
  template: `
  <v-card class="mt-4">
    <v-card-title>
      Laps for Race: {{ race.Name || race.Title }}
      <v-spacer></v-spacer>
      <v-btn color="primary" @click="addLap">
        <v-icon left>mdi-plus</v-icon>
        Add Lap
      </v-btn>
      <v-btn icon @click="close">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-text>
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
      </v-data-table>
    </v-card-text>
  </v-card>
  `
}
