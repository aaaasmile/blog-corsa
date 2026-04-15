<template>
  <v-container>
    <v-card>
      <v-card-subtitle>Race Editor</v-card-subtitle>

      <v-divider></v-divider>

      <v-card-actions>
        <v-btn color="primary" @click="addRace">
          <v-icon left>mdi-plus</v-icon>
          Add Race
        </v-btn>
      </v-card-actions>

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
        <template v-slot:item.actions="{ item }">
          <v-icon small class="mr-2" color="blue" @click="editLaps(item)">mdi-format-list-numbered</v-icon>
          <v-icon small class="mr-2" @click="editRace(item)">mdi-pencil</v-icon>
          <v-icon small color="red" @click="deleteRace(item)">mdi-delete</v-icon>
        </template>
      </v-data-table>
    </v-card>
    <div v-if="showRaceItem" class="mb-12">
      <RaceItem
        :value="editedRace"
        @save="onSaveRace"
        @cancel="onCancelRace"
      />
    </div>
    <div v-if="showLapEditor" class="mb-12">
      <RaceLapEditor
        :race="selectedRaceForLaps"
        @close="closeLapEditor"
      />
    </div>
    <v-dialog v-model="dialogDelete" persistent max-width="350">
      <v-card>
        <v-card-title class="headline">Confirm Delete</v-card-title>
        <v-card-text>
          Are you sure you want to delete this race?
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" text @click="confirmDelete">OK</v-btn>
          <v-btn color="green darken-1" text @click="dialogDelete = false">Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
