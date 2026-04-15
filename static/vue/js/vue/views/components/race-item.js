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
    }
  },
  watch: {
    value(newVal) {
      this.item = Object.assign({}, newVal)
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
          <v-col cols="12" sm="6">
            <v-text-field
              v-model="item.Distance"
              label="Distance"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="6">
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
                  v-model="item.Date"
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
                v-model="item.Date"
                @input="dateMenu = false"
              ></v-date-picker>
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
              label="Result Time"
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
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="item.pace_minkm"
              label="Pace (min/km)"
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
              label="Rank Class"
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
              label="Class Name"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.sport_type_id"
              label="Sport Type ID"
              type="number"
              outlined
              dense
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model.number="item.race_subtype_id"
              label="Race Subtype ID"
              type="number"
              outlined
              dense
            ></v-text-field>
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
