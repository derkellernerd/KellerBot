import { autoImplement } from 'src/helper/functions';

export interface ApiKellerBotChatCommand {
  ID: number;
  Command: string;
  Used: number;
  TimeoutInSeconds: number;
  LastUsed: Date;
  Action: string;
}

export class KellerBotChatCommand extends autoImplement<ApiKellerBotChatCommand>() {
  static fromApi(item: ApiKellerBotChatCommand) : KellerBotChatCommand {
    return new KellerBotChatCommand(item);
  }
}
