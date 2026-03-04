export interface Client {
  profile_id?: string;
  address: string;
  languages: string[];
  issue_detail?: string[];
  consultation_modes: string[];
  rage_price: number;
  availability: {
    days: string[];
    time_slots: string[];
  };
  preferred_therapist_gender?: string;
  experience_level?: string;
  therapist_specialization?: string[];
  urgency_level?: string;
  session_duration?: number;
  preferred_therapist_language?: string[];
  is_flexible_with_schedule?: boolean;
}

export interface ProfessionalTitle {
  code: string;
  display: string;
}

export interface Degree {
  type: string;
  field: string;
  institution: string;
  year: number;
}

export interface Certification {
  name: string;
  issuer: string;
  year: number;
}

export interface ProfessionalInfo {
  title: ProfessionalTitle;
  degrees: Degree[];
  certifications: Certification[];
  experience_years: number;
}

export interface TherapistV1 {
  profile_id?: string;
  address: string;
  languages: string[];
  specialization: string[];
  consultation_modes: string[];
  experience: number;
  rating?: number;
  currency: string;
  rage_price: number;
  is_available?: boolean;
  availability: {
    days: string[];
    time_slots: string[];
  };
  current_session?: string[];
  matched_clients?: string[];
  avatar_override?: string;
  name: string;
  professional_info?: ProfessionalInfo;
}

export type UserProfile = Client | TherapistV1;

export type ProfileType = 'client' | 'therapist';

export const WEEKDAYS = [
  { value: 'monday', label: 'Monday', labelVi: 'Thứ Hai' },
  { value: 'tuesday', label: 'Tuesday', labelVi: 'Thứ Ba' },
  { value: 'wednesday', label: 'Wednesday', labelVi: 'Thứ Tư' },
  { value: 'thursday', label: 'Thursday', labelVi: 'Thứ Năm' },
  { value: 'friday', label: 'Friday', labelVi: 'Thứ Sáu' },
  { value: 'saturday', label: 'Saturday', labelVi: 'Thứ Bảy' },
  { value: 'sunday', label: 'Sunday', labelVi: 'Chủ Nhật' },
];

export const TIME_SLOTS = [
  {
    value: 'morning',
    label: 'Morning (6AM - 12PM)',
    labelVi: 'Sáng (6h - 12h)',
  },
  {
    value: 'afternoon',
    label: 'Afternoon (12PM - 6PM)',
    labelVi: 'Chiều (12h - 18h)',
  },
  {
    value: 'evening',
    label: 'Evening (6PM - 10PM)',
    labelVi: 'Tối (18h - 22h)',
  },
];

export const CONSULTATION_MODES = [
  { value: 'online', label: 'Online Video', labelVi: 'Video trực tuyến' },
  { value: 'in-person', label: 'In Person', labelVi: 'Trực tiếp' },
  { value: 'phone', label: 'Phone Call', labelVi: 'Điện thoại' },
  { value: 'chat', label: 'Text Chat', labelVi: 'Nhắn tin' },
];

export const LANGUAGES = [
  { value: 'en', label: 'English', labelVi: 'Tiếng Anh' },
  { value: 'vi', label: 'Vietnamese', labelVi: 'Tiếng Việt' },
  { value: 'zh', label: 'Chinese', labelVi: 'Tiếng Trung' },
  { value: 'ja', label: 'Japanese', labelVi: 'Tiếng Nhật' },
  { value: 'ko', label: 'Korean', labelVi: 'Tiếng Hàn' },
];

export const SPECIALIZATIONS = [
  {
    value: 'anxiety',
    label: 'Anxiety & Stress',
    labelVi: 'Lo âu & Căng thẳng',
  },
  { value: 'depression', label: 'Depression', labelVi: 'Trầm cảm' },
  {
    value: 'relationship',
    label: 'Relationship Issues',
    labelVi: 'Vấn đề quan hệ',
  },
  { value: 'trauma', label: 'Trauma & PTSD', labelVi: 'Chấn thương tâm lý' },
  { value: 'family', label: 'Family Therapy', labelVi: 'Trị liệu gia đình' },
  { value: 'addiction', label: 'Addiction', labelVi: 'Nghiện ngập' },
  { value: 'eating', label: 'Eating Disorders', labelVi: 'Rối loạn ăn uống' },
  { value: 'grief', label: 'Grief & Loss', labelVi: 'Đau buồn mất mát' },
];

export const URGENCY_LEVELS = [
  {
    value: 'immediate',
    label: 'Immediate (Within 24 hours)',
    labelVi: 'Khẩn cấp (Trong 24h)',
  },
  {
    value: 'urgent',
    label: 'Urgent (Within 3 days)',
    labelVi: 'Gấp (Trong 3 ngày)',
  },
  {
    value: 'moderate',
    label: 'Moderate (Within a week)',
    labelVi: 'Vừa (Trong tuần)',
  },
  { value: 'flexible', label: 'Flexible', labelVi: 'Linh hoạt' },
];

export const EXPERIENCE_LEVELS = [
  { value: 'any', label: 'No Preference', labelVi: 'Không yêu cầu' },
  { value: 'entry', label: '1-3 years', labelVi: '1-3 năm' },
  { value: 'mid', label: '3-7 years', labelVi: '3-7 năm' },
  { value: 'senior', label: '7+ years', labelVi: '7+ năm' },
];

export const GENDER_PREFERENCES = [
  { value: 'any', label: 'No Preference', labelVi: 'Không yêu cầu' },
  { value: 'male', label: 'Male', labelVi: 'Nam' },
  { value: 'female', label: 'Female', labelVi: 'Nữ' },
  { value: 'non-binary', label: 'Non-binary', labelVi: 'Phi nhị phân' },
];

export const PROFESSIONAL_TITLES = [
  { code: 'MD', display: 'Medical Doctor (Psychiatrist)' },
  { code: 'PhD', display: 'Doctor of Philosophy (Clinical Psychology)' },
  { code: 'PsyD', display: 'Doctor of Psychology' },
  { code: 'LCSW', display: 'Licensed Clinical Social Worker' },
  { code: 'LPC', display: 'Licensed Professional Counselor' },
  { code: 'LMFT', display: 'Licensed Marriage and Family Therapist' },
  { code: 'MSW', display: 'Master of Social Work' },
  { code: 'MA', display: 'Master of Arts (Counseling/Psychology)' },
];

// Session Models
export interface ConsultationSession {
  session_id: string;
  therapist_id: string;
  client_id: string;
  mode: string;
  start_time: string;
  end_time: string;
  scheduled_date?: string;
  time_zone?: string;
  status: string;
  price: number;
  payment_status?: string;
  location_info?: LocationInfo;
  call_session_id?: string;
  created_at: string;
}

export interface LocationInfo {
  link?: string;
  passcode?: string;
  address_line?: string;
  city?: string;
  coordinates?: string[];
}

export interface SessionRequest {
  client_id?: string;
  therapist_id: string;
  start_time: string;
  end_time: string;
  time_zone: string;
  scheduled_date?: string;
  mode: string;
  price: number;
  location_info?: LocationInfo;
}

export interface DeleteSessionRequest {
  sessionID: string;
}

// Matching Models
export interface MatchRequest {
  client_id?: string;
  therapist_id?: string;
  preferences?: any;
}

export interface MatchResponse {
  match_id: string;
  client_id: string;
  therapist_id: string;
  match_score: number;
  created_at: string;
}

// Social Models
export interface Follow {
  follower_id: string;
  following_id: string;
  created_at: string;
}

export interface Bookmark {
  user_id: string;
  post_id: string;
  created_at: string;
}

export interface FollowersResponse {
  followers: string[];
  total: number;
}

export interface FollowingResponse {
  following: string[];
  total: number;
}
