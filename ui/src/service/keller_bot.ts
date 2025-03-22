import type { AxiosResponse } from 'axios'
import type { BaseResponse } from '@/models/base_response.ts'
import type {
  ApiKellerBotCommand,
  ApiKellerBotCommandCreateRequest
} from '@/models/keller_bot_command.ts'
import api from '@/boot/axios.ts'

class KellerBot {
  getCommands():Promise<AxiosResponse<BaseResponse<ApiKellerBotCommand[]>>> {
    return api.get('/api/v1/command')
  }

  createCommand(commandCreateRequest: ApiKellerBotCommandCreateRequest) : Promise<AxiosResponse<BaseResponse<ApiKellerBotCommand>>> {
    console.log('creating')
    return api.post('/api/v1/command', commandCreateRequest);
  }
}

export default new KellerBot();
