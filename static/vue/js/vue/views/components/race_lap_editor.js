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
        { text: 'Lap Pace (min/km)', value: 'lap_pace_minkm' },
        { text: 'Tot Meter Race', value: 'tot_meter_race' },
        { text: 'Tot Race (min/km)', value: 'tot_race_minkm' },
        { text: 'Race ID', value: 'race_id' },
        { text: 'Tot Race Time', value: 'tot_race_time' },
        { text: 'Tot KM Race', value: 'tot_km_race' },
      ],
      laps: []
    }
  },
  methods: {
    close() {
      this.$emit('close')
    }
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
      <v-data-table
        :headers="headers"
        :items="laps"
        item-key="id"
        class="elevation-1"
        :search="search"
      >
        <template v-slot:top>
          <v-toolbar flat>
            <v-toolbar-title>Lap List</v-toolbar-title>
            <v-spacer></v-spacer>
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
