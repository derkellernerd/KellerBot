import type { AxiosResponse } from 'axios'
import type { BaseResponse } from 'src/models/base_response.ts'
import { api } from 'boot/axios';
import type { ApiKellerBotChatEvent, ApiKellerBotEvent } from 'src/models/keller_bot_event';
import type { ApiKellerBotTwitchEvent } from 'src/models/keller_bot_twitch_event';
import type {
  ApiKellerBotAction,
} from 'src/models/keller_bot_action';
import { type ApiKellerBotChatCommand } from 'src/models/keller_bot_chat_command';

class KellerBot {
  twitchLogin() : Promise<AxiosResponse<BaseResponse<never>>> {
    return api.get('/api/v1/twitch/login')
  }

  getChatCommands():Promise<AxiosResponse<BaseResponse<ApiKellerBotChatCommand[]>>> {
    return api.get('/api/v1/chat_command')
  }

  createChatCommand(chatCommand: ApiKellerBotChatCommand) : Promise<AxiosResponse<BaseResponse<ApiKellerBotChatCommand>>> {
    return api.post('/api/v1/chat_command', chatCommand);
  }

  deleteChatCommand(chatCommandID: number) : Promise<AxiosResponse<never>> {
    return api.delete(`/api/v1/chat_command/${chatCommandID}`)
  }

  getTwitchEvents():Promise<AxiosResponse<BaseResponse<ApiKellerBotTwitchEvent[]>>> {
    return api.get('/api/v1/event/twitch')
  }

  createTwitchEvent(twitchEventCreateRequest: ApiKellerBotTwitchEvent) : Promise<AxiosResponse<BaseResponse<ApiKellerBotTwitchEvent>>> {
    return api.post('/api/v1/event/twitch', twitchEventCreateRequest);
  }

  testTwitchEvent(twitchEventId: number) : Promise<AxiosResponse<never>> {
    return api.post(`/api/v1/event/twitch/${twitchEventId}/action/test`)
  }

  getChatStream() : Promise<AxiosResponse> {
    return api.get('/api/v1/event/chat', {
      headers: {
        'Accept': 'text/event-stream',
      },
      responseType: 'stream',
      adapter: 'fetch'
    })
  }

  getAlertStream() : Promise<AxiosResponse> {
    return api.get('/api/v1/event/alert', {
      headers: {
        'Accept': 'text/event-stream',
      },
      responseType: 'stream',
      adapter: 'fetch'
    })
  }

  createChatEventTest(message: ApiKellerBotChatEvent) : Promise<AxiosResponse<unknown>> {
    return api.post('/api/v1/event/chat', message)
  }

  getActions():Promise<AxiosResponse<BaseResponse<ApiKellerBotAction[]>>> {
    return api.get('/api/v1/action')
  }

  getAction(actionId: number):Promise<AxiosResponse<BaseResponse<ApiKellerBotAction>>> {
    return api.get(`/api/v1/action/${actionId}`)
  }

  createAction(action: ApiKellerBotAction):Promise<AxiosResponse<BaseResponse<ApiKellerBotAction>>> {
    return api.post('/api/v1/action', action)
  }

  updateAction(actionId: number, action: ApiKellerBotAction) : Promise<AxiosResponse<BaseResponse<ApiKellerBotAction>>> {
    return api.put(`/api/v1/action/${actionId}`, action)
  }

  uploadActionFile(actionId: number, file: File) : Promise<AxiosResponse<BaseResponse<ApiKellerBotAction>>> {
    const fileData = new FormData()
    fileData.append('file', file);

    return api.post(`/api/v1/action/${actionId}/upload`, fileData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }

  getEvents() : Promise<AxiosResponse<BaseResponse<ApiKellerBotEvent[]>>> {
    return api.get('/api/v1/event')
  }

  replayEvent(eventId: number) : Promise<AxiosResponse<never>> {
    return api.post(`/api/v1/event/${eventId}/action/replay`)
  }
}

export default new KellerBot();
