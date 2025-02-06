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
        })
    },
    methods: {

    },
    template: `
  <v-container>
    <v-card>
      <v-card-subtitle>Comments to moderate</v-card-subtitle>
    </v-card>
  </v-container>
`
}