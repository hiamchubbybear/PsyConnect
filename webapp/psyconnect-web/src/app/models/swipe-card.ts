export interface Therapist {
  id: string;
  name: string;
  title: string;
  rating: number;
  reviews: number;
  specialties: string[];
  nextAvailable: string;
  pricePerHour: string;
  avatarUrl?: string;
}
