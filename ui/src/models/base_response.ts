export interface BaseResponse<T> {
    Data: T;
    Timestamp: Date;
    Error: string;
}
