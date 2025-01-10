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


export interface UserProfile {
    id: number;
    email: string;
    gender: 'male' | 'female' | 'unknown';
    weight_kg: number;
    created_at: string;
    updated_at: string;
}

export interface UpdateUserProfileRequest {
    gender: 'male' | 'female' | 'unknown';
    weight_kg: number;
}
export interface DrinkTemplate {
    id: number;
    name: string;
    type: string;
    size_value: number;
    size_unit: string;
    abv: number;
}


export interface DrinkTemplatesResponse {
    drink_templates: DrinkTemplate[];
}

export interface DrinkTemplateResponse {
    drink_template: DrinkTemplate;
}

export interface DrinkLog extends DrinkTemplate {
    logged_at: string;
}

export interface DrinkLogsResponse {
    drink_logs: DrinkLog[];
}

export interface CreateDrinkLogRequest {
    name: string;
    type: string;
    size_value: number;
    size_unit: string;
    abv: number;
    logged_at?: string;
}

export interface UpdateDrinkLogRequest extends CreateDrinkLogRequest {
    id: number;
}

export interface DeleteDrinkLogResponse {
    id: number;
}

export interface CreateDrinkLogResponse {
    id: number;
}



export interface ParseDrinkLogRequest {
    text: string;
}

export interface ParseDrinkLogResponse {
    drink_template: DrinkTemplate;
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
    max_bac: number;
    max_bac_time: string;
    sober_since_time: string;
    total_drinks: number;
    drinking_since_time: string;
    duration_over_bac: number;
    estimated_sober_time: string;
}

export interface BACCalculationResponse {
    timeline: BACPoint[];
    summary: BACSummary;
}

export interface CurrentBACResponse {
    current_bac: number;
    bac_status: BACStatus;
    is_sober: boolean;
    estimated_sober_time: string;
    last_calculated: string;
}

export interface DrinkStatsPoint {
    drink_count: number;
    time_period: string;
    total_standard_drinks: number;
}

export interface DrinkStatsResponse {
    stats: DrinkStatsPoint[];
}

export type DrinkStatsPeriod = 'daily' | 'weekly' | 'monthly' | 'yearly';

export type BACCategories = 'sober' | 'light' | 'heavy';

export interface MonthlyBACStats {
    year: number;
    month: number;
    counts: { [key in BACCategories]: number };
    total: number;
}

export interface MonthlyBACStatsResponse {
    stats: MonthlyBACStats[];
    categories: BACCategories[];
}

