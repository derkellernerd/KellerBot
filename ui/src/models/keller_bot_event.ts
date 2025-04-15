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

export interface ApiKellerBotEvent {
  ID: number;
  CreatedAt: Date;
  EventName: string;
  ExecutedActionName: string;
  Source: string;
}

export class KellerBotEvent extends autoImplement<ApiKellerBotEvent>() {
  static fromApi(item: ApiKellerBotEvent) : KellerBotEvent {
    return new KellerBotEvent(item);
  }
}
