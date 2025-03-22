import { autoImplement } from '@/helper/functions.ts'

export interface ApiKellerBotCommand {
  Id: number;
  Command: string;
  Type: string;
  Used: number;
  Data: unknown;
  CreatedAt: Date;
}

export interface ApiKellerBotCommandCreateRequest {
  Command: string;
  Type: string;
  Data: unknown;
}

export interface ApiKellerBotCommandMessageAction {
  Message: string;
}

export class KellerBotCommand extends autoImplement<ApiKellerBotCommand>() {
  static fromApi(apiItem: ApiKellerBotCommand) : KellerBotCommand {
    return new KellerBotCommand(apiItem);
  }
}
