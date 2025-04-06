import type { AxiosResponse } from 'axios'
import type { BaseResponse } from 'src/models/base_response.ts'
import type {
  ApiKellerBotCommand,
  ApiKellerBotCommandCreateRequest
} from 'src/models/keller_bot_command.ts'
import { api } from 'boot/axios';
import type { ApiKellerBotChatEvent } from 'src/models/keller_bot_event';
import type { ApiKellerBotAlert, ApiKellerBotAlertCreateRequest } from 'src/models/keller_bot_alert';
import type { ApiKellerBotTwitchEvent } from 'src/models/keller_bot_twitch_event';

class KellerBot {
  getCommands():Promise<AxiosResponse<BaseResponse<ApiKellerBotCommand[]>>> {
    return api.get('/api/v1/command')
  }

  createCommand(commandCreateRequest: ApiKellerBotCommandCreateRequest) : Promise<AxiosResponse<BaseResponse<ApiKellerBotCommand>>> {
    return api.post('/api/v1/command', commandCreateRequest);
  }

  getTwitchEvents():Promise<AxiosResponse<BaseResponse<ApiKellerBotTwitchEvent[]>>> {
    return api.get('/api/v1/event/twitch')
  }

  createTwitchEvent(twitchEventCreateRequest: ApiKellerBotTwitchEvent) : Promise<AxiosResponse<BaseResponse<ApiKellerBotTwitchEvent>>> {
    return api.post('/api/v1/event/twitch', twitchEventCreateRequest);
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

  getAlerts():Promise<AxiosResponse<BaseResponse<ApiKellerBotAlert[]>>> {
    return api.get('/api/v1/alert')
  }

  createAlert(alert: ApiKellerBotAlertCreateRequest):Promise<AxiosResponse<BaseResponse<ApiKellerBotAlert>>> {
    return api.post('/api/v1/alert', alert)
  }

  uploadAlertFile(alertId: number, file: File) : Promise<AxiosResponse<BaseResponse<ApiKellerBotAlert>>> {
    const fileData = new FormData()
    fileData.append('file', file);

    return api.post(`/api/v1/alert/${alertId}/upload`, fileData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}

export default new KellerBot();
