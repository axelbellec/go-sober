import type {
    ApiError,
    // Auth types
    UserLoginRequest,
    UserLoginResponse,
    UserSignupRequest,
    UserSignupResponse,
    UserMeResponse,
    // Drink options types
    DrinkTemplatesResponse,
    DrinkTemplateResponse,
    // Drink logs types
    DrinkLogsResponse,
    CreateDrinkLogRequest,
    UpdateDrinkLogRequest,
    CreateDrinkLogResponse,
    ParseDrinkLogRequest,
    ParseDrinkLogResponse,
    // BAC types
    BACCalculationResponse,
    CurrentBACResponse,
    // Analytics types
    DrinkStatsResponse,
    DrinkStatsPeriod,
    MonthlyBACStatsResponse,
    DeleteDrinkLogResponse,
    // User types
    UserProfile,
    UpdateUserProfileRequest,
} from '../types/api';


export class ApiService {
    private readonly baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    // Helper function to handle API responses
    private async handleJsonResponse<T>(response: Response): Promise<T> {
        const data = await response.json();

        if (!response.ok) {
            throw data as ApiError;
        }

        return data as T;
    }

    private async fetchWithAuth(url: string, options: RequestInit = {}) {
        const tokenKey = process.env.NEXT_PUBLIC_LOCALSTORAGE_TOKEN_KEY!;
        const token = localStorage.getItem(tokenKey);
        const headers = {
            ...options.headers,
            Authorization: `Bearer ${token}`,
        };

        const response = await fetch(url, {
            ...options,
            headers,
        });

        if (response.status === 401) {
            // Token expired or invalid
            localStorage.removeItem(tokenKey);
            window.location.href = "/login";
            throw new Error("Unauthorized");
        }

        return response;
    }

    // Auth API
    async login(email: string, password: string): Promise<UserLoginResponse> {
        const request: UserLoginRequest = { email, password };
        const response = await fetch(`${this.baseUrl}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(request),
        });

        return this.handleJsonResponse<UserLoginResponse>(response);
    }

    async signup(email: string, password: string): Promise<UserSignupResponse> {
        const request: UserSignupRequest = { email, password };
        const response = await fetch(`${this.baseUrl}/auth/signup`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(request),
        });

        return this.handleJsonResponse<UserSignupResponse>(response);
    }

    async getCurrentUser(): Promise<UserMeResponse> {
        const response = await this.fetchWithAuth(`${this.baseUrl}/auth/me`);
        return this.handleJsonResponse<UserMeResponse>(response);
    }

    // Drink Options API
    async getDrinkTemplates(): Promise<DrinkTemplatesResponse> {
        const response = await fetch(`${this.baseUrl}/drink-templates`);
        return this.handleJsonResponse<DrinkTemplatesResponse>(response);
    }

    async getDrinkTemplate(id: number): Promise<DrinkTemplateResponse> {
        const response = await fetch(`${this.baseUrl}/drink-templates/${id}`);
        return this.handleJsonResponse<DrinkTemplateResponse>(response);
    }


    // Drink Logs API
    async getDrinkLogs(params?: {
        page?: number;
        page_size?: number;
        start_date?: string;
        end_date?: string;
        drink_type?: string;
        min_abv?: number;
        max_abv?: number;
        sort_by?: 'logged_at' | 'abv' | 'size_value' | 'name' | 'type';
        sort_order?: 'asc' | 'desc';
    }): Promise<DrinkLogsResponse> {
        const queryParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined) {
                    queryParams.append(key, value.toString());
                }
            });
        }

        const url = `${this.baseUrl}/drink-logs${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
        const response = await this.fetchWithAuth(url);
        return this.handleJsonResponse<DrinkLogsResponse>(response);
    }

    async createDrinkLog(
        data: CreateDrinkLogRequest
    ): Promise<CreateDrinkLogResponse> {
        const response = await this.fetchWithAuth(`${this.baseUrl}/drink-logs`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        return this.handleJsonResponse<CreateDrinkLogResponse>(response);
    }


    async updateDrinkLog(
        data: UpdateDrinkLogRequest
    ): Promise<void> {
        const response = await this.fetchWithAuth(`${this.baseUrl}/drink-logs`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        return this.handleJsonResponse<void>(response);
    }

    async deleteDrinkLog(id: number): Promise<DeleteDrinkLogResponse> {
        const response = await this.fetchWithAuth(`${this.baseUrl}/drink-logs/${id}`, {
            method: 'DELETE',
        });
        return this.handleJsonResponse<DeleteDrinkLogResponse>(response);
    }

    async parseDrinkLog(text: string): Promise<ParseDrinkLogResponse> {
        const request: ParseDrinkLogRequest = { text };

        const response = await this.fetchWithAuth(`${this.baseUrl}/drink-logs/parse`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(request),
        });

        return this.handleJsonResponse<ParseDrinkLogResponse>(response);
    }

    // Analytics API
    async getCurrentBAC(weightKg: number, gender: string): Promise<CurrentBACResponse> {
        const params = new URLSearchParams({
            weight_kg: weightKg.toString(),
            gender,
        });

        const response = await this.fetchWithAuth(
            `${this.baseUrl}/bac/current?${params.toString()}`
        );

        return this.handleJsonResponse<CurrentBACResponse>(response);
    }

    async getBACTimeline(
        startTime: string,
        endTime: string,
        weightKg: number,
        gender: string,
        timeStepMins: number
    ): Promise<BACCalculationResponse> {
        const params = new URLSearchParams({
            start_time: startTime,
            end_time: endTime,
            weight_kg: weightKg.toString(),
            gender,
            time_step_mins: timeStepMins.toString(),
        });

        const response = await this.fetchWithAuth(
            `${this.baseUrl}/bac/timeline?${params.toString()}`
        );

        return this.handleJsonResponse<BACCalculationResponse>(response);
    }

    async getDrinkStats(period: DrinkStatsPeriod, params?: {
        start_date?: string;
        end_date?: string;
    }): Promise<DrinkStatsResponse> {
        const queryParams = new URLSearchParams({
            period: period,
            ...(params?.start_date && { start_date: params.start_date }),
            ...(params?.end_date && { end_date: params.end_date })
        });

        const response = await this.fetchWithAuth(
            `${this.baseUrl}/analytics/drink-stats?${queryParams.toString()}`
        );
        return this.handleJsonResponse<DrinkStatsResponse>(response);
    }

    async getMonthlyBACStats(params?: {
        start_date?: string;
        end_date?: string;
    }): Promise<MonthlyBACStatsResponse> {
        const queryParams = new URLSearchParams({
            ...(params?.start_date && { start_date: params.start_date }),
            ...(params?.end_date && { end_date: params.end_date })
        });

        const response = await this.fetchWithAuth(
            `${this.baseUrl}/analytics/monthly-bac?${queryParams.toString()}`
        );
        return this.handleJsonResponse<MonthlyBACStatsResponse>(response);
    }

    async getUserProfile(): Promise<UserProfile> {
        const response = await this.fetchWithAuth(`${this.baseUrl}/users/profile`);
        return this.handleJsonResponse<UserProfile>(response);
    }

    async updateUserProfile(data: UpdateUserProfileRequest): Promise<UserProfile> {
        const response = await this.fetchWithAuth(`${this.baseUrl}/users/profile`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        return this.handleJsonResponse<UserProfile>(response);
    }
}


