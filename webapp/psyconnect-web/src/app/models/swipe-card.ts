export interface Therapist {
  id: string;
  name: string;
  title: string;
  rating: number;
  reviews: number;
  specialties: string[];
  availability?: string;
  pricePerHour: string;
  avatarUrl?: string;
  languages?: string[];
  modes: string[];
}
