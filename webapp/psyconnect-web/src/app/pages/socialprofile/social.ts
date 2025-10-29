interface ProfileData {
  name: string;
  role: string;
  avatar: string;
  email: string;
  phone: string;
  city: string;
  specialization: string;
  experience: number;
  license: string;
  education: string;
  verified: boolean;
  completedSessions: number;
  rating: number;
}

interface Session {
  id: number;
  clientInitials: string;
  date: string;
  duration: number;
  topic: string;
  status: string;
  notes: string;
}
