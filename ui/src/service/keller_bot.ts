import type { AxiosResponse } from 'axios'
import type { BaseResponse } from 'src/models/base_response.ts'
import type {
  ApiKellerBotCommand,
  ApiKellerBotCommandCreateRequest
} from 'src/models/keller_bot_command.ts'
import { api } from 'boot/axios';

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
}

export default new KellerBot();
