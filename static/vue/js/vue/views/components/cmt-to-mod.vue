<template>
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
</template>