import API from '../../apicaller.js?version=101'

const buildselectedParam = (that) => {
  let arr = []
  that.$store.state.admin.cmtSelected.forEach(element => {
    arr.push(element.Id)
  });
  console.log("selection list", arr)
  return arr
}

export default {
  data() {
    return {
      loadingCmt: false,
      loadingPage: false,
      dialogApprove: false,
      current_command: '',
      current_selection: 0,
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
      this.current_command = 'approve'
      const items = buildselectedParam(this)
      if (items.length <= 0){
        return
      }
      this.current_selection = items.length
      this.dialogApprove = true
    },
    rejectCmtList() {
      this.current_command = 'reject'
      const items = buildselectedParam(this)
      if (items.length <= 0){
        return
      }
      this.current_selection = items.length
      this.dialogApprove = true
    },
    okDoIt() {
      const cmd_todo = this.current_command
      console.log('command to do: ', cmd_todo)
      this.dialogApprove = false
      const items = buildselectedParam(this)
      this.loadingCmt = true
      let para = { cmd: cmd_todo, list: items }
      console.log('do on selected comment list', cmd_todo, items)
      API.DoCmt(this, para, () => {
        this.loadingCmt = false
      }, () => {
        console.log('something went wrong')
        this.loadingCmt = false
      })
    }
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
        item-key="Id"
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
      <v-card-actions>
        <v-btn text :loading="loadingCmt" @click="approveCmtList">
          Approve
        </v-btn>
        <v-btn color="red" text :loading="loadingCmt" @click="rejectCmtList">
          Reject
        </v-btn>
      </v-card-actions>
    </v-card>
    <v-dialog v-model="dialogApprove" persistent max-width="290">
      <v-card>
        <v-card-title class="headline">Question</v-card-title>
        <v-card-text
          >Do you want {{current_command}} on {{current_selection}} selected comments?</v-card-text
        >
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" text @click="okDoIt">OK</v-btn>
          <v-btn color="green darken-1" text @click="dialogApprove = false"
            >Cancel</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
`
}