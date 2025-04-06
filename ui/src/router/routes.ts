import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/IndexPage.vue') },
      { name: 'CommandOverview', path: 'command', component: () => import('pages/CommandsOverviewPage.vue') },
      { name: 'AlertOverview', path: 'alert', component: () => import('pages/AlertOverviewPage.vue') },
      { name: 'Event', path: 'event', component: () => import('pages/EventPage.vue') },
      { name: 'TwitchEvent', path: 'event/twitch', component: () => import('pages/TwitchEventOverviewPage.vue')}
    ],
  },
  {
    path: '/box',
    children: [
      {path: 'chat', component: () => import('pages/boxes/ChatBoxPage.vue')},
      {path: 'alert', component: () => import('pages/boxes/AlertBoxPage.vue')}
    ]
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;

