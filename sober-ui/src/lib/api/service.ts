import type {
    ApiError,
    // Auth types
    UserLoginRequest,
    UserLoginResponse,
    UserSignupRequest,
    UserSignupResponse,
    UserMeResponse,
    // Drink options types
    DrinkOptionsResponse,
    DrinkOptionResponse,
    UpdateDrinkOptionRequest,
    // Drink logs types
    DrinkLogsResponse,
    CreateDrinkLogRequest,
    UpdateDrinkLogRequest,
    CreateDrinkLogResponse,
    ParseDrinkLogRequest,
    ParseDrinkLogResponse,
    // Analytics types
    BACCalculationResponse,
    CurrentBACResponse,
} from '../types/api';


export class ApiService {
    private readonly baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    // Helper function to handle API responses
    private async handleResponse<T>(response: Response): Promise<T> {
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

        return this.handleResponse<UserLoginResponse>(response);
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

        return this.handleResponse<UserSignupResponse>(response);
    }

    async getCurrentUser(): Promise<UserMeResponse> {
        const response = await this.fetchWithAuth(`${this.baseUrl}/auth/me`);
        return this.handleResponse<UserMeResponse>(response);
    }

    // Drink Options API
    async getDrinkOptions(): Promise<DrinkOptionsResponse> {
        const response = await fetch(`${this.baseUrl}/drink-options`);
        return this.handleResponse<DrinkOptionsResponse>(response);
    }

    async getDrinkOption(id: number): Promise<DrinkOptionResponse> {
        const response = await fetch(`${this.baseUrl}/drink-options/${id}`);
        return this.handleResponse<DrinkOptionResponse>(response);
    }

    async updateDrinkOption(
        id: number,
        data: UpdateDrinkOptionRequest
    ): Promise<void> {
        const response = await this.fetchWithAuth(`${this.baseUrl}/drink-options/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        return this.handleResponse<void>(response);
    }

    async deleteDrinkOption(id: number) {
        const response = await this.fetchWithAuth(`${this.baseUrl}/drink-options/${id}`, {
            method: 'DELETE',
        });

        return this.handleResponse(response);
    }

    // Drink Logs API
    async getDrinkLogs(): Promise<DrinkLogsResponse> {
        const response = await this.fetchWithAuth(`${this.baseUrl}/drink-logs`);
        return this.handleResponse<DrinkLogsResponse>(response);
    }

    async createDrinkLog(
        drinkOptionId: number,
        loggedAt: string
    ): Promise<CreateDrinkLogResponse> {
        const request: CreateDrinkLogRequest = {
            drink_option_id: drinkOptionId,
            logged_at: loggedAt,
        };

        const response = await this.fetchWithAuth(`${this.baseUrl}/drink-logs`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(request),
        });

        return this.handleResponse<CreateDrinkLogResponse>(response);
    }

    async updateDrinkLog(
        id: number,
        data: UpdateDrinkLogRequest
    ): Promise<void> {
        const response = await this.fetchWithAuth(`${this.baseUrl}/drink-logs/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
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

        return this.handleResponse<ParseDrinkLogResponse>(response);
    }

    // Analytics API
    async getCurrentBAC(weightKg: number, gender: string): Promise<CurrentBACResponse> {
        const params = new URLSearchParams({
            weight_kg: weightKg.toString(),
            gender,
        });

        const response = await this.fetchWithAuth(
            `${this.baseUrl}/analytics/current/bac?${params.toString()}`
        );

        return this.handleResponse<CurrentBACResponse>(response);
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
            `${this.baseUrl}/analytics/timeline/bac?${params.toString()}`
        );

        return this.handleResponse<BACCalculationResponse>(response);
    }
}


