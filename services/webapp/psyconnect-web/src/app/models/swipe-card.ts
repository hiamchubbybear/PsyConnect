export interface Therapist {
  profileId: string;
  address?: string;
  languages?: string[];
  specialization?: string[];
  consultationModes?: string[];
  experience?: number;
  rating?: number;
  currency: string;
  ragePrice?: number;
  isAvailable?: boolean;
  availability?: {
    days?: string[];
    timeSlots?: string[];
  };
  currentSession?: string[];
  matchedClients?: string[];
  avatarOverride?: string;
  name: string;
  professionalInfo?: ProfessionalInfo;
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
  degrees?: Degree[];
  certifications?: Certification[];
  experienceYears?: number;
}
export function mapTherapistResponse(apiData: any): Therapist[] {
  if (!apiData?.data || !Array.isArray(apiData.data)) return [];

  return apiData.data.map((item: any) => ({
    profileId: item.profile_id,
    name: item.name,
    address: item.address,
    languages: item.languages,
    specialization: item.specialization,
    consultationModes: item.consultation_modes,
    experience: item.experience,
    rating: item.rating,
    currency: item.currency,
    ragePrice: item.ragePrice,
    isAvailable: item.is_available,
    avatarOverride: item.avatar_override,
    availability: item.availability
      ? {
          days: item.availability.days,
          timeSlots: item.availability.time_slots,
        }
      : undefined,
    currentSession: [],
    matchedClients: [],
    professionalInfo: item.professional_info
      ? {
          title: item.professional_info.title,
          degrees: item.professional_info.degrees,
          certifications: item.professional_info.certifications,
          experienceYears: item.professional_info.experience_years,
        }
      : undefined,
  }));
}
