import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    component: () => import('pages/LoginPage.vue')
  },
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { name: 'Dashboard', path: '', component: () => import('pages/DashboardPage.vue') },
      {
        path: 'chat_command',
        children: [
          {
            name: 'ChatCommandOverview',
            path: '',
            component: () => import('pages/chat_command/ChatCommandOverviewPage.vue')
          }
        ]
      },
      {
        path: 'action',
        children: [
          {
            name: 'ActionOverview',
            path: '',
            component: () => import('pages/actions/ActionOverviewPage.vue'),
          },
          {
            name: 'ActionCreate',
            path: 'create',
            component: () => import('pages/actions/ActionDetailPage.vue'),
          },
          {
            name: 'ActionDetail',
            path: ':actionId',
            component: () => import('pages/actions/ActionDetailPage.vue'),
          },
        ],
      },
      { name: 'Event', path: 'event', component: () => import('pages/EventPage.vue') },
      {
        name: 'TwitchEvent',
        path: 'event/twitch',
        component: () => import('pages/TwitchEventOverviewPage.vue'),
      },
    ],
  },
  {
    path: '/box',
    children: [
      { path: 'chat', component: () => import('pages/boxes/ChatBoxPage.vue') },
      { path: 'alert', component: () => import('pages/boxes/AlertBoxPage.vue') },
      { path: 'obs', component: () => import('pages/boxes/ObsControlBoxPage.vue')}
    ],
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
