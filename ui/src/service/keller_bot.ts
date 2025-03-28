import type { AxiosResponse } from 'axios'
import type { BaseResponse } from 'src/models/base_response.ts'
import type {
  ApiKellerBotCommand,
  ApiKellerBotCommandCreateRequest
} from 'src/models/keller_bot_command.ts'
import { api } from 'boot/axios';
import type { ApiKellerBotChatEvent } from 'src/models/keller_bot_event';
import type { ApiKellerBotAlert } from 'src/models/keller_bot_alert';

class KellerBot {
  getCommands():Promise<AxiosResponse<BaseResponse<ApiKellerBotCommand[]>>> {
    return api.get('/api/v1/command')
  }

  createCommand(commandCreateRequest: ApiKellerBotCommandCreateRequest) : Promise<AxiosResponse<BaseResponse<ApiKellerBotCommand>>> {
    console.log('creating')
    return api.post('/api/v1/command', commandCreateRequest);
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
}

export default new KellerBot();
