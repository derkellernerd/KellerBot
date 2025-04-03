import { autoImplement } from 'src/helper/functions';

export enum KellerBotAlertType {
  Sound = 'SOUND',
  Video = 'VIDEO',
  Gif = 'GIF',
  GifSound = 'GIF_SOUND',
  Composition = 'COMPOSITION'
}

export interface ApiKellerBotAlertTypeComposition {
  Alerts: ApiKellerBotAlert[];
}

export interface ApiKellerBotAlertTypeCompositionCreateRequest {
  AlertNames: string[];
}

export interface ApiKellerBotAlert {
  ID: number;
  Name: string;
  Type: KellerBotAlertType;
  Data: ApiKellerBotAlertTypeComposition|null;
}

export interface ApiKellerBotAlertCreateRequest {
  Name: string;
  Type: KellerBotAlertType;
  Data: ApiKellerBotAlertTypeCompositionCreateRequest|null;
}

export class KellerBotAlert extends autoImplement<ApiKellerBotAlert>() {
  static fromApi(item: ApiKellerBotAlert): KellerBotAlert {
    return new KellerBotAlert(item);
  }

  get composition() : ApiKellerBotAlertTypeComposition {
    return this.Data as ApiKellerBotAlertTypeComposition
  }
}
