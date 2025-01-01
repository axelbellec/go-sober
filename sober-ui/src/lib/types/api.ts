// Error response type
export interface ApiError {
    code: number;
    type: 'validation' | 'database' | 'entity';
    message: string;
    details?: object[];
    correlation_id: string;
}

// Auth types
export interface UserLoginRequest {
    email: string;
    password: string;
}

export interface UserLoginResponse {
    message: string;
    token: string;
}

export interface UserSignupRequest {
    email: string;
    password: string;
}

export interface UserSignupResponse {
    message: string;
}

export interface UserMeResponse {
    email: string;
    user_id: number;
}

// Drink types
export interface DrinkOption {
    id: number;
    name: string;
    type: string;
    abv: number;
    size_value: number;
    size_unit: string;
}

export interface DrinkOptionsResponse {
    drink_options: DrinkOption[];
}

export interface DrinkOptionResponse {
    drink_option: DrinkOption;
}

export interface UpdateDrinkOptionRequest {
    name: string;
    type: string;
    abv: number;
    size_value: number;
    size_unit: string;
}

export interface DrinkLog {
    id: number;
    user_id: number;
    drink_option_id: number;
    drink_name: string;
    abv: number;
    size_value: number;
    size_unit: string;
    logged_at: string;
}

export interface DrinkLogsResponse {
    drink_logs: DrinkLog[];
}

export interface CreateDrinkLogRequest {
    drink_option_id: number;
    logged_at: string;
}


export interface CreateDrinkLogResponse {
    id: number;
}

export interface ParseDrinkLogRequest {
    text: string;
}

export interface ParseDrinkLogResponse {
    drink_option: DrinkOption;
    confidence: number;
}

// Analytics types
export type BACStatus =
    | 'Sober'
    | 'Minimal'
    | 'Light'
    | 'Mild'
    | 'Significant'
    | 'Severe'
    | 'Dangerous';

export interface BACPoint {
    time: string;
    bac: number;
    is_over_bac: boolean;
    status: BACStatus;
}

export interface BACSummary {
    drinking_since_time: string;
    sober_since_time: string;
    max_bac: number;
    max_bac_time: string;
    duration_over_bac: number;
    total_drinks: number;
}

export interface BACCalculationResponse {
    summary: BACSummary;
    timeline: BACPoint[];
}

export interface CurrentBACResponse {
    current_bac: number;
    bac_status: BACStatus;
    is_sober: boolean;
    estimated_sober_time: string;
    last_calculated: string;
}