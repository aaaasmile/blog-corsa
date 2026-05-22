<template>
  <v-card>
    <v-card-title class="subheading grey--text">Races</v-card-title>
    <v-card-subtitle>Race Management</v-card-subtitle>
    <v-divider></v-divider>
    <v-alert v-if="errorMsg" type="error" dismissible>{{ errorMsg }}</v-alert>
    <RaceEditor
      :races="races"
      :loadingData="loadingData"
      @save-race="onSaveRace"
      @delete-race="onDeleteRace"
    />
  </v-card>
</template>
