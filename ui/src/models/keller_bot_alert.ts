import { autoImplement } from 'src/helper/functions';

export enum KellerBotAlertType {
  Sound = 'SOUND',
  Video = 'VIDEO',
  Gif = 'GIF',
  GifSound = 'GIF_SOUND',
  Composition = 'COMPOSITION',
  Text = 'TEXT',
  Chat = 'CHAT',
}

export interface ApiKellerBotAlertTypeComposition {
  Alerts: ApiKellerBotAlert[];
}

export interface ApiKellerBotAlertTypeText {
  Text: string;
}

export interface ApiKellerBotAlertTypeChat {
  Chat: string;
}

export interface ApiKellerBotAlertTypeCompositionCreateRequest {
  AlertNames: string[];
}

export interface ApiKellerBotAlert {
  ID: number;
  Name: string;
  Type: KellerBotAlertType;
  Data: ApiKellerBotAlertTypeComposition|ApiKellerBotAlertTypeText|ApiKellerBotAlertTypeChat|null;
  DurationInSeconds: number;
}

export interface ApiKellerBotAlertCreateRequest {
  Name: string;
  Type: KellerBotAlertType;
  Data: ApiKellerBotAlertTypeCompositionCreateRequest|ApiKellerBotAlertTypeText|ApiKellerBotAlertTypeChat|null;
}

export class KellerBotAlert extends autoImplement<ApiKellerBotAlert>() {
  static fromApi(item: ApiKellerBotAlert): KellerBotAlert {
    return new KellerBotAlert(item);
  }

  get composition() : ApiKellerBotAlertTypeComposition {
    return this.Data as ApiKellerBotAlertTypeComposition
  }

  get textAlert() : ApiKellerBotAlertTypeText {
    return this.Data as ApiKellerBotAlertTypeText;
  }

  get duration() : number {
    return this.DurationInSeconds > 0 ? this.DurationInSeconds * 1000 : 2000;
  }
}
