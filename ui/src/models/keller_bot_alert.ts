import { autoImplement } from 'src/helper/functions';

export enum KellerBotAlertType {
  Sound = 'SOUND',
  Video = 'VIDEO',
  Gif = 'GIF',
  GifSound = 'GIF_SOUND',
}

export interface ApiKellerBotAlert {
  ID: number;
  Name: string;
  Type: KellerBotAlertType;
  Data: unknown;
}

export class KellerBotAlert extends autoImplement<ApiKellerBotAlert>() {
  static fromApi(item: ApiKellerBotAlert): KellerBotAlert {
    return new KellerBotAlert(item);
  }
}
