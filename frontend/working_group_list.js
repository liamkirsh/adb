import WorkingGroupList from './WorkingGroupList.vue';
import Vue from 'vue';

export function initializeApp() {
  var vm = new Vue({
    el: '#app',
    render(h) {
      return h(WorkingGroupList);
    },
  });
}
