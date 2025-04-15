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

export enum ApiKellerBotEventSource {
  TWITCH= 'TWITCH',
}

export interface ApiKellerBotEvent {
  ID: number;
  CreatedAt: Date;
  EventName: string;
  ExecutedActionName: string;
  Source: ApiKellerBotEventSource;
  Payload: never;
}

export class KellerBotEvent extends autoImplement<ApiKellerBotEvent>() {
  static fromApi(item: ApiKellerBotEvent) : KellerBotEvent {
    return new KellerBotEvent(item);
  }

  get eventInfo() {
    if (this.Source == ApiKellerBotEventSource.TWITCH) {
      switch (this.EventName) {
        case 'channel.follow':
          // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
          return `username: ${this.Payload['user_name']}`
      }
    }
  }

  get sourceIcon() {
    switch (this.Source) {
      case ApiKellerBotEventSource.TWITCH:
        return 'fa-solid fa-twitch';
    }
  }
}

