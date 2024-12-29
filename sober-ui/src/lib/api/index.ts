import { ApiService } from '@/lib/api/service';

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080/api/v1';

export const apiService = new ApiService(API_BASE);