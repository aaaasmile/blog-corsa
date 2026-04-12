<template>
  <v-container>
    <v-card>
      <v-card-subtitle>Race Editor</v-card-subtitle>

      <v-divider></v-divider>

      <v-data-table
        :headers="headers"
        :items="races"
        :loading="loadingData"
        item-key="Id"
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
