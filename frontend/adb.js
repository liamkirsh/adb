import Vue from 'vue';
import VueRouter from 'vue-router';
import AdbShell from 'AdbShell.vue';
import ActivistList from 'ActivistList.vue';
import CirclesList from 'CirclesList.vue';
import EventEdit from 'EventEdit.vue';
import EventList from 'EventList.vue';
import UserList from 'UserList.vue';
import WorkingGroupList from 'WorkingGroupList.vue';

Vue.use(VueRouter);

// NOTE: If you add a new path here, you have to add it to main.go as well.
const routes = [
    // Events
    { path: '/', component: EventEdit, props: { id: "0" } },
    { path: '/list_events', component: EventList },
    { path: '/update_event/:id', component: EventEdit, props: route => ({ id: route.params.id }) },

    // Connections
    { path: '/new_connection', component: EventEdit, props: { connections: true, id: "0" } },
    { path: '/list_connections', component: EventList, props: { connections: true } },
    { path: '/update_connection/:id', component: EventEdit, props: route => ({ connections: true, id: route.params.id }) },
    { path: '/activist_pool', component: ActivistList, props: { title: "Recruitment Connections", view: "activist_pool" } },

    // Circles
    { path: '/circle_member_prospects', component: ActivistList, props: { title: "Circle Member Prospects", view: "circle_member_prospects" } },
    { path: '/list_circles', component: CirclesList },

    // Chapter Members
    { path: '/chapter_member_prospects', component: ActivistList, props: { title: "Chapter Member Prospects", view: "chapter_member_prospects" } },
    { path: '/chapter_member_development', component: ActivistList, props: { title: "Chapter Member Development", view: "chapter_member_development" } },

    // Organizers
    { path: '/organizer_prospects', component: ActivistList, props: { title: "Organizer Prospects", view: "organizer_prospects" } },
    { path: '/activist_development', component: ActivistList, props: { title: "Organizer Development", view: "development" } },
    { path: '/list_working_groups', component: WorkingGroupList },

    // All
    { path: '/list_activists', component: ActivistList, props: { title: "All Activists", view: "all_activists" } },
    { path: '/leaderboard', component: ActivistList, props: { title: "Leaderboard", view: "leaderboard" } },
    { path: '/activist_recruitment', component: ActivistList, props: { title: "Activist Recruitment", view: "activist_recruitment" } },
    { path: '/activist_actionteam', component: ActivistList, props: { title: "Action Team", view: "action_team" } },

    // Users
    { path: '/admin/users', component: UserList },
];

new Vue({
    el: "#adb",
    render: (h) => h(AdbShell),
    router: new VueRouter({
        mode: "history",
        routes,
    }),
});
