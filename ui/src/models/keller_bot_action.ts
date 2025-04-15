import { autoImplement } from 'src/helper/functions';

export enum KellerBotActionType {
  Sound = 'SOUND',
  Gif = 'GIF',
  Composition = 'COMPOSITION',
  Text = 'TEXT',
  ChatMessage = 'CHAT_MESSAGE',
  ChatAnswer = 'CHAT_ANSWER',
  Http = 'HTTP',
}

export interface ApiKellerBotActionTypeChatMessage {
  ChatMessage: string;
}

export interface ApiKellerBotActionTypeHttp {
  Payload: unknown;
  HttpMethod: string;
  Uri: string;
}

export interface ApiKellerBotActionTypeText {
  Text: string;
  DurationMs: number;
}

export interface ApiKellerBotActionTypeGif {
  FileName: string;
  DurationMs: number;
}

export interface ApiKellerBotActionTypeSound {
  FileName: string;
  Gain: number;
  DurationMs: number;
}

export interface ApiKellerBotActionTypeComposition {
  Actions: string[];
  DurationMs: number;
  PostAction?: string;
}

export type ApiKellerBotActionTypes =
  | ApiKellerBotActionTypeGif
  | ApiKellerBotActionTypeSound
  | ApiKellerBotActionTypeText
  | ApiKellerBotActionTypeComposition
  | ApiKellerBotActionTypeChatMessage
  | ApiKellerBotActionTypeHttp;

export interface ApiKellerBotAction {
  ID: number;
  ActionName: string;
  ActionType: KellerBotActionType;
  Data: ApiKellerBotActionTypes;
}

export interface ApiKellerBotActionCreateRequest {
  Name: string;
  ActionType: KellerBotActionType;
  Data: ApiKellerBotActionTypes | null;
}

export class KellerBotAction extends autoImplement<ApiKellerBotAction>() {
  static fromApi(item: ApiKellerBotAction): KellerBotAction {
    return new KellerBotAction(item);
  }

  static actionTypes() {
    return [
      KellerBotActionType.Gif,
      KellerBotActionType.Composition,
      KellerBotActionType.Sound,
      KellerBotActionType.Text,
      KellerBotActionType.ChatMessage,
      KellerBotActionType.Http
    ];
  }

  get duration() {
    switch (this.ActionType) {
      case KellerBotActionType.Sound:
        return this.actionSound.DurationMs > 0 ? this.actionSound.DurationMs : 2000;
      case KellerBotActionType.Gif:
        return this.actionGif.DurationMs > 0 ? this.actionGif.DurationMs : 2000;
      case KellerBotActionType.Text:
        return this.actionText.DurationMs > 0 ? this.actionSound.DurationMs : 2000;
      default:
        return 2000;
    }
  }

  get actionText() {
    return this.Data as ApiKellerBotActionTypeText;
  }

  get actionGif() {
    return this.Data as ApiKellerBotActionTypeGif;
  }

  get actionSound() {
    return this.Data as ApiKellerBotActionTypeSound;
  }

  get actionComposition() {
    return this.Data as ApiKellerBotActionTypeComposition;
  }

  get actionChatMessage() {
    return this.Data as ApiKellerBotActionTypeChatMessage;
  }

  get actionHttp() {
    return this.Data as ApiKellerBotActionTypeHttp;
  }

  get isNew() {
    return !this.ID || this.ID == 0;
  }

}
