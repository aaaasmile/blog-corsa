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
</template>