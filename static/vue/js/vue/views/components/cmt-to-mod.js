import API from '../../apicaller.js?version=101'

export default {
    data() {
        return {
        }
    },
    created() {
        console.log('cmt-to-mod created')
        let para = { cmd: 'list', type: 'to_moderate' }
        API.DoCmt(this, para)
    },
    computed: {
        ...Vuex.mapState({
          Comments: state => {
            return state.admin.comments_to_mod
          },
        })
    },
    methods: {

    },
    template: `
  <v-container>
    <v-card>
      <v-card-subtitle>Comments to moderate</v-card-subtitle>

      <v-divider></v-divider>

      <v-list
        lines="three"
        select-strategy="leaf"
      >
        <v-list-item
          v-for="cmt in Comments"
          :key="cmt.Id"
          :subtitle="cmt.PostId"
          :title="cmt.Comment"
        >
          <template v-slot:append>
            <v-btn
              color="grey-lighten-1"
              icon="mdi-information"
              variant="text"
            ></v-btn>
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </v-container>`
}