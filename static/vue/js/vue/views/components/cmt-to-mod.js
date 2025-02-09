import API from '../../apicaller.js?version=101'
const buildselectedParam = (that) => {
  let arr = []
  let stickarr = []
  that.$store.state.admin.cmtSelected.forEach(element => {
    arr.push(element.KeyStore)
  });
  that.$store.state.admin.cmtSelected.forEach(element => {
    if (element.is_sticky) {
      stickarr.push(element.KeyStore)
    }
  });
  console.log("selection list", arr)
  let para = { selected: arr, sticky: stickarr, cmd: 'approve' }
  return para
}

export default {
  data() {
    return {
      loadingCmt: false,
      loadingPage: false,
      dialogApprove: false,
      keyreq: '',
      transition: 'scale-transition',
      search: '',
      loadingData: false,
      headers: [
        { text: 'ID', value: 'Id' },
        { text: 'Comment', value: 'Comment' },
      ],
    }
  },
  created() {
    console.log('cmt-to-mod created')
    let para = { cmd: 'list', type: 'to_moderate' }
    API.DoCmt(this, para)
  },
  computed: {
    cmtSelected: {
      get() {
        return (this.$store.state.admin.cmtSelected)
      },
      set(newVal) {
        this.$store.commit('setCmtSelected', newVal)
      }
    },
    ...Vuex.mapState({
      Comments: state => {
        return state.admin.comments_to_mod
      },
    })
  },
  methods: {
    approveCmtList() {
      this.dialogApprove = false
      this.loadingCmt = true
      console.log('Approve selected comment list')
      let para = buildselectedParam(this)
      API.DoCmt(this, para)
    },
  },
  template: `
  <v-container>
    <v-card>
      <v-card-subtitle>Comments to moderate</v-card-subtitle>

      <v-divider></v-divider>

      <v-data-table
        v-model="cmtSelected"
        :headers="headers"
        :items="Comments"
        :loading="loadingData"
        item-key="KeyStore"
        show-select
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
      </v-data-table>
    </v-card>
  </v-container>
`
}