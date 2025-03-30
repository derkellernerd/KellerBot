import { autoImplement } from 'src/helper/functions';

export enum KellerBotCommandType {
  MessageAction = 'MESSAGE_ACTION',
  HttpAction = 'HTTP_ACTION',
  AlertAction = 'ALERT_ACTION',
}

export interface ApiKellerBotCommand {
  Id: number;
  Command: string;
  Type: KellerBotCommandType;
  Used: number;
  Data: ApiKellerBotCommandAlertAction | ApiKellerBotCommandMessageAction;
  CreatedAt: Date;
  LastUsed: Date;
  TimeoutInSeconds: number;
}

export interface ApiKellerBotCommandCreateRequest {
  Command: string;
  Type: KellerBotCommandType;
  Data: ApiKellerBotCommandAlertAction | ApiKellerBotCommandMessageAction;
  TimeoutInSeconds: number;
}


export interface ApiKellerBotCommandAlertAction {
  Alert: string;
}

export interface ApiKellerBotCommandMessageAction {
  Message: string;
}

export class KellerBotCommand extends autoImplement<ApiKellerBotCommand>() {
  static fromApi(apiItem: ApiKellerBotCommand): KellerBotCommand {
    return new KellerBotCommand(apiItem);
  }

  get typeLabel() {
    switch (this.Type) {
      case KellerBotCommandType.MessageAction:
        return 'Message';
      case KellerBotCommandType.HttpAction:
        return 'HTTP';
      case KellerBotCommandType.AlertAction:
        return 'Alert';
    }
  }
}
