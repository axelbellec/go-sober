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
    CreateDrinkLogResponse,
    ParseDrinkLogRequest,
    ParseDrinkLogResponse,
    // Analytics types
    BACCalculationResponse,
    CurrentBACResponse,
} from '../types/api';
import { fetchWithAuth } from '../utils';

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? 'http://localhost:8080/api/v1';

// Helper function to handle API responses
async function handleResponse<T>(response: Response): Promise<T> {
    const data = await response.json();

    if (!response.ok) {
        throw data as ApiError;
    }

    return data as T;
}

// Auth API
export async function login(email: string, password: string): Promise<UserLoginResponse> {
    const request: UserLoginRequest = { email, password };
    const response = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
    });

    return handleResponse<UserLoginResponse>(response);
}

export async function signup(email: string, password: string): Promise<UserSignupResponse> {
    const request: UserSignupRequest = { email, password };
    const response = await fetch(`${API_BASE}/auth/signup`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
    });

    return handleResponse<UserSignupResponse>(response);
}

export async function getCurrentUser(): Promise<UserMeResponse> {
    const response = await fetchWithAuth(`${API_BASE}/auth/me`);
    return handleResponse<UserMeResponse>(response);
}

// Drink Options API
export async function getDrinkOptions(): Promise<DrinkOptionsResponse> {
    const response = await fetch(`${API_BASE}/drink-options`);
    return handleResponse<DrinkOptionsResponse>(response);
}

export async function getDrinkOption(id: number): Promise<DrinkOptionResponse> {
    const response = await fetch(`${API_BASE}/drink-options/${id}`);
    return handleResponse<DrinkOptionResponse>(response);
}

export async function updateDrinkOption(
    id: number,
    data: UpdateDrinkOptionRequest
): Promise<void> {
    const response = await fetchWithAuth(`${API_BASE}/drink-options/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });

    return handleResponse<void>(response);
}

export async function deleteDrinkOption(id: number) {
    const response = await fetchWithAuth(`${API_BASE}/drink-options/${id}`, {
        method: 'DELETE',
    });

    return handleResponse(response);
}

// Drink Logs API
export async function getDrinkLogs(): Promise<DrinkLogsResponse> {
    const response = await fetchWithAuth(`${API_BASE}/drink-logs`);
    return handleResponse<DrinkLogsResponse>(response);
}

export async function createDrinkLog(
    drinkOptionId: number,
    loggedAt: string
): Promise<CreateDrinkLogResponse> {
    const request: CreateDrinkLogRequest = {
        drink_option_id: drinkOptionId,
        logged_at: loggedAt,
    };

    const response = await fetchWithAuth(`${API_BASE}/drink-logs`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
    });

    return handleResponse<CreateDrinkLogResponse>(response);
}

export async function parseDrinkLog(text: string): Promise<ParseDrinkLogResponse> {
    const request: ParseDrinkLogRequest = { text };

    const response = await fetchWithAuth(`${API_BASE}/drink-logs/parse`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
    });

    return handleResponse<ParseDrinkLogResponse>(response);
}

// Analytics API
export async function getCurrentBAC(weightKg: number, gender: string): Promise<CurrentBACResponse> {
    const params = new URLSearchParams({
        weight_kg: weightKg.toString(),
        gender,
    });

    const response = await fetchWithAuth(
        `${API_BASE}/analytics/current/bac?${params.toString()}`
    );

    return handleResponse<CurrentBACResponse>(response);
}

export async function getBACTimeline(
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

    const response = await fetchWithAuth(
        `${API_BASE}/analytics/timeline/bac?${params.toString()}`
    );

    return handleResponse<BACCalculationResponse>(response);
} 