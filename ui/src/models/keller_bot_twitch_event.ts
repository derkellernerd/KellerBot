import { autoImplement } from 'src/helper/functions';

export interface ApiKellerBotTwitchEvent {
  ID: number;
  TwitchEventSubscription: string;
  ActionName: string;
}

export class KellerBotTwitchEvent extends autoImplement<ApiKellerBotTwitchEvent>() {
  static fromApi(item: ApiKellerBotTwitchEvent) : KellerBotTwitchEvent {
    return new KellerBotTwitchEvent(item);
  }
}
