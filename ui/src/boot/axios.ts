import { defineBoot } from '#q-app/wrappers';
import type { AxiosError} from 'axios';
import axios, { type AxiosInstance } from 'axios';
import { showErrorToast } from 'src/helper/functions';
import type { BaseResponse } from 'src/models/base_response';

declare module 'vue' {
  interface ComponentCustomProperties {
    $axios: AxiosInstance;
    $api: AxiosInstance;
  }
}

// Be careful when using SSR for cross-request state pollution
// due to creating a Singleton instance here;
// If any client changes this (global) instance, it might be a
// good idea to move this instance creation inside of the
// "export default () => {}" function below (which runs individually
// for each client)
const api = axios.create({
  baseURL: 'http://localhost:8080/',
  headers: {
    'Content-Type': 'application/json',
  }
});

api.interceptors.response.use((response) => {
  return response;
}, (error: AxiosError<BaseResponse<never>>) => {
  showErrorToast(error.response?.data?.Error ?? error.message, error.status?.toString() ?? error.name);
  return Promise.reject(error)
})

export default defineBoot(({ app }) => {
  // for use inside Vue files (Options API) through this.$axios and this.$api

  app.config.globalProperties.$axios = axios;
  // ^ ^ ^ this will allow you to use this.$axios (for Vue Options API form)
  //       so you won't necessarily have to import axios in each vue file

  app.config.globalProperties.$api = api;
  // ^ ^ ^ this will allow you to use this.$api (for Vue Options API form)
  //       so you can easily perform requests against your app's API
});

export { api };
