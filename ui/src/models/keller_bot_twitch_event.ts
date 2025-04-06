import { autoImplement } from 'src/helper/functions';

export interface ApiKellerBotTwitchEvent {
  TwitchEventSubscription: string;
  AlertName: string;
}

export class KellerBotTwitchEvent extends autoImplement<ApiKellerBotTwitchEvent>() {
  static fromApi(item: ApiKellerBotTwitchEvent) : KellerBotTwitchEvent {
    return new KellerBotTwitchEvent(item);
  }
}
