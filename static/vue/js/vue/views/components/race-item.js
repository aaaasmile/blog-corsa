export default {
  props: {
    value: {
      type: Object,
      default: () => ({
        Name: '',
        Title: '',
        Distance: '',
        Date: '',
        ascending_meter: 0,
        descending_meter: 0,
        rank_global: 0,
        rank_gender: 0,
        rank_class: 0,
        class_name: '',
        sport_type_id: 0,
        result_time: '',
        pace_kmh: 0,
        pace_minkm: '',
        comment: '',
        race_subtype_id: 0,
        km_length: 0,
      })
    }
  },
  data() {
    return {
      item: Object.assign({
        Name: '',
        Title: '',
        Distance: '',
        Date: '',
        ascending_meter: 0,
        descending_meter: 0,
        rank_global: 0,
        rank_gender: 0,
        rank_class: 0,
        class_name: '',
        sport_type_id: 0,
        result_time: '',
        pace_kmh: 0,
        pace_minkm: '',
        comment: '',
        race_subtype_id: 0,
        km_length: 0,
      }, this.value),
      dateMenu: false,
      timeMenu: false,
      sportTypes: [
        { text: 'run', value: 0 },
        { text: 'duathlon', value: 1 },
      ],
      raceSubtypes: [
        { text: 'Mountain', value: 0 },
        { text: 'Up to 10 KM', value: 1 },
        { text: 'Half Marathon', value: 2 },
        { text: 'Marathon', value: 3 },
        { text: 'Ultra Road', value: 4 },
        { text: 'Ultra Trail', value: 5 },
        { text: 'Short Trail', value: 6 },
        { text: 'Road other distance', value: 7 },
        { text: 'Backyard', value: 8 },
      ],
    }
  },
  watch: {
    value(newVal) {
      this.item = Object.assign({}, newVal)
    }
  },
  computed: {
    dateObj: {
      get() {
        return this.item.Date ? this.item.Date.split(' ')[0] : ''
      },
      set(val) {
        let time = this.timeObj || '00:00:00'
        this.item.Date = val ? `${val} ${time}` : ''
      }
    },
    timeObj: {
      get() {
        if (!this.item.Date) return '00:00:00'
        let parts = this.item.Date.split(' ')
        return parts.length > 1 ? parts[1] : '00:00:00'
      },
      set(val) {
        let date = this.dateObj || new Date().toISOString().substr(0, 10)
        let parts = val.split(':')
        if (parts.length === 2) val += ':00'
        this.item.Date = `${date} ${val}`
      }
    }
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
    <v-card-title>Race Details</v-card-title>
    <v-card-text>
      <v-form>
        <v-row>
          <v-col cols="12" sm="6">
            <v-text-field
              v-model="item.Name"
              label="Name"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="6">
            <v-text-field
              v-model="item.Title"
              label="Title"
              outlined
              dense
            ></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="item.Distance"
              label="Distance (in meter)"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-menu
              v-model="dateMenu"
              :close-on-content-click="false"
              :nudge-right="40"
              transition="scale-transition"
              offset-y
              min-width="auto"
            >
              <template v-slot:activator="{ on, attrs }">
                <v-text-field
                  v-model="dateObj"
                  label="Date"
                  prepend-icon="mdi-calendar"
                  readonly
                  outlined
                  dense
                  v-bind="attrs"
                  v-on="on"
                ></v-text-field>
              </template>
              <v-date-picker
                v-model="dateObj"
                @input="dateMenu = false"
              ></v-date-picker>
            </v-menu>
          </v-col>
          <v-col cols="12" sm="4">
            <v-menu
              v-model="timeMenu"
              :close-on-content-click="false"
              :nudge-right="40"
              transition="scale-transition"
              offset-y
              min-width="auto"
            >
              <template v-slot:activator="{ on, attrs }">
                <v-text-field
                  v-model="timeObj"
                  label="Time"
                  prepend-icon="mdi-clock-outline"
                  readonly
                  outlined
                  dense
                  v-bind="attrs"
                  v-on="on"
                ></v-text-field>
              </template>
              <v-time-picker
                v-model="timeObj"
                use-seconds
                format="24hr"
              ></v-time-picker>
            </v-menu>
          </v-col>
        </v-row>

        <v-divider class="mb-4"></v-divider>

        <v-row>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.km_length"
              label="KM Length"
              type="number"
              step="0.1"
              readonly
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.ascending_meter"
              label="Ascending Meter"
              type="number"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.descending_meter"
              label="Descending Meter"
              type="number"
              outlined
              dense
            ></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="item.result_time"
              label="Result Time (hh:mm:ss)"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.pace_kmh"
              label="Pace (km/h)"
              type="number"
              step="0.1"
              readonly
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="item.pace_minkm"
              label="Pace (min/km)"
              readonly
              outlined
              dense
            ></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.rank_global"
              label="Rank Global"
              type="number"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.rank_gender"
              label="Rank Gender"
              type="number"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.rank_class"
              label="Rank Age Group"
              type="number"
              outlined
              dense
            ></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="item.class_name"
              label="Age Group"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-select
              v-model.number="item.sport_type_id"
              :items="sportTypes"
              item-text="text"
              item-value="value"
              label="Sport Type"
              outlined
              dense
            ></v-select>
          </v-col>
          <v-col cols="12" sm="4">
            <v-select
              v-model.number="item.race_subtype_id"
              :items="raceSubtypes"
              item-text="text"
              item-value="value"
              label="Race Subtype"
              outlined
              dense
            ></v-select>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12">
            <v-textarea
              v-model="item.comment"
              label="Comment"
              outlined
              dense
              rows="3"
            ></v-textarea>
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
