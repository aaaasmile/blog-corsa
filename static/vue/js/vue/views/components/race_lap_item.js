export default {
  props: {
    value: {
      type: Object,
      default: () => ({
        id: 0,
        lap_number: 0,
        lap_time: '',
        lap_meter: 0,
        tot_race_minkm: '',
        tot_race_time: '',
        tot_km_race: 0,
      })
    }
  },
  data() {
    return {
      item: Object.assign({
        id: 0,
        lap_number: 0,
        lap_time: '',
        lap_meter: 0,
        tot_race_minkm: '',
        tot_race_time: '',
        tot_km_race: 0,
      }, this.value),
      timeRules: [
        v => !v || /^\d+:[0-5]\d:[0-5]\d$/.test(v) || 'Format must be hh:mm:ss'
      ],
    }
  },
  watch: {
    value(newVal) {
      this.item = Object.assign({}, newVal)
    },
  },
  methods: {
    save() {
      this.$emit('save', Object.assign({}, this.item))
    },
    cancel() {
      this.$emit('cancel')
    },
  },
  template: `
  <v-card class="mt-4">
    <v-card-title>Lap Details</v-card-title>
    <v-card-text>
      <v-form>
        <v-row>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.lap_number"
              label="Lap Number"
              type="number"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="item.lap_time"
              label="Lap Time (hh:mm:ss)"
              :rules="timeRules"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.lap_meter"
              label="Lap Meter"
              type="number"
              outlined
              dense
            ></v-text-field>
          </v-col>
        </v-row>

        <v-divider class="mb-4"></v-divider>

        <v-row>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="item.tot_race_minkm"
              label="Tot Race (min/km)"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="item.tot_race_time"
              label="Tot Race Time"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.tot_km_race"
              label="Tot KM Race"
              type="number"
              step="0.1"
              outlined
              dense
            ></v-text-field>
          </v-col>
        </v-row>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn text @click="cancel">Cancel</v-btn>
      <v-btn color="primary" @click="save">Save</v-btn>
    </v-card-actions>
  </v-card>
`
}
