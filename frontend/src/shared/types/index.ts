export interface PageToken { id: number; created_at: string }
export interface PaginatedResp<T> { list: T[]; page_token: string; has_more: boolean }
