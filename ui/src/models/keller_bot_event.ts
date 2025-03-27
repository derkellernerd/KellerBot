import { autoImplement } from 'src/helper/functions';

export interface ApiKellerBotChatEvent {
  User: string;
  Message: string;
}

export class KellerBotChatEvent extends autoImplement<ApiKellerBotChatEvent>() {
  static fromApi(item: ApiKellerBotChatEvent) : KellerBotChatEvent {
    return new KellerBotChatEvent(item);
  }
}
