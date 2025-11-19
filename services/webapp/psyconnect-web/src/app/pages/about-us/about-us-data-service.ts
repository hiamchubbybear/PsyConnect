import { Injectable } from '@angular/core';

export interface Partner {
  id: string;
  name: string;
  role: string;
  specialization: string[];
  verified: boolean;
  avatar: string;
  experience: number;
  rating: number;
  reviewCount: number;
  education: string[];
  certificates: Certificate[];
  languages: string[];
  location: string;
}

export interface Certificate {
  id: string;
  name: string;
  issuer: string;
  issueDate: string;
  expiryDate?: string;
  credentialId?: string;
  verificationUrl?: string;
}

export interface Dev {
  id: string;
  name: string;
  role: string;
  github?: string;
  linkedin?: string;
  avatar?: string;
  bio?: string;
}

export interface Service {
  id: string;
  name: string;
  desc: string;
  icon: string;
  duration?: string;
  priceRange?: string;
}

export interface Statistic {
  value: string;
  label: string;
  icon: string;
}

export interface TeamMember {
  id: string;
  name: string;
  role: string;
  bio: string;
  avatar: string;
  social?: {
    linkedin?: string;
    email?: string;
  };
}

@Injectable({
  providedIn: 'root',
})
export class AboutUsDataService {
  getBrandInfo() {
    return {
      name: 'PsyConnect',
      logo: 'assets/about/logo.svg',
      tagline: '',
      established: '2025',
    };
  }

  getStatistics(): Statistic[] {
    return [
      { value: '500+', label: 'verified_experts', icon: 'users' },
      { value: '10,000+', label: 'successful_sessions', icon: 'heart' },
      { value: '4.8/5', label: 'average_rating', icon: 'star' },
      { value: '24/7', label: 'support_available', icon: 'clock' },
    ];
  }

  getPartners(): Partner[] {
    return [
      {
        id: 'p1',
        name: 'TS. Nguyễn Thị Minh An',
        role: 'Clinical Psychologist',
        specialization: ['anxiety_disorders', 'depression', 'trauma_therapy'],
        verified: true,
        avatar: 'assets/partners/dr-an.jpg',
        experience: 8,
        rating: 4.9,
        reviewCount: 127,
        education: ['PhD Psychology - VNU', 'MSc Clinical Psychology - HMU'],
        certificates: [
          {
            id: 'cert1',
            name: 'Certified CBT Therapist',
            issuer: 'International CBT Institute',
            issueDate: '2020-03-15',
            expiryDate: '2025-03-15',
            credentialId: 'CBT-2020-AN-001',
            verificationUrl: 'https://verify.cbt-institute.org/CBT-2020-AN-001',
          },
          {
            id: 'cert2',
            name: 'Trauma-Informed Care Specialist',
            issuer: 'Trauma Recovery Institute',
            issueDate: '2021-07-20',
            credentialId: 'TIC-2021-AN-045',
          },
        ],
        languages: ['Vietnamese', 'English'],
        location: 'Ho Chi Minh City',
      },
      {
        id: 'p2',
        name: 'ThS. Trần Minh Đức',
        role: 'Family Therapist',
        specialization: ['family_counseling', 'couples_therapy', 'child_psychology'],
        verified: true,
        avatar: 'assets/partners/duc-tran.jpg',
        experience: 6,
        rating: 4.7,
        reviewCount: 89,
        education: ['MSc Family Therapy - UEH', 'BA Psychology - HCMUS'],
        certificates: [
          {
            id: 'cert3',
            name: 'Licensed Family Therapist',
            issuer: 'Vietnam Psychology Association',
            issueDate: '2019-11-10',
            credentialId: 'LFT-2019-TD-123',
          },
        ],
        languages: ['Vietnamese'],
        location: 'Hanoi',
      },
    ];
  }

  getServices(): Service[] {
    return [
      {
        id: 's1',
        name: 'individual_counseling',
        desc: 'individual_counseling_desc',
        icon: 'user-circle',
        duration: '50-60 minutes',
        priceRange: '300,000 - 800,000 VND',
      },
      {
        id: 's2',
        name: 'group_therapy',
        desc: 'group_therapy_desc',
        icon: 'users',
        duration: '90 minutes',
        priceRange: '200,000 - 400,000 VND',
      },
      {
        id: 's3',
        name: 'couples_therapy',
        desc: 'couples_therapy_desc',
        icon: 'heart',
        duration: '60-90 minutes',
        priceRange: '500,000 - 1,200,000 VND',
      },
      {
        id: 's4',
        name: 'online_consultation',
        desc: 'online_consultation_desc',
        icon: 'video',
        duration: '45-60 minutes',
        priceRange: '250,000 - 600,000 VND',
      },
    ];
  }

  getTeam(): TeamMember[] {
    return [
      {
        id: 'tm1',
        name: 'Dr. Sarah Nguyen',
        role: 'clinical_director',
        bio: 'clinical_director_bio',
        avatar: 'assets/team/sarah-nguyen.jpg',
        social: {
          linkedin: 'https://linkedin.com/in/sarah-nguyen-psych',
          email: 'sarah@psyconnect.com',
        },
      },
      {
        id: 'tm2',
        name: 'Mark Thompson',
        role: 'head_of_operations',
        bio: 'head_of_operations_bio',
        avatar: 'assets/team/mark-thompson.jpg',
        social: {
          linkedin: 'https://linkedin.com/in/mark-thompson-ops',
        },
      },
    ];
  }

  getContact() {
    return {
      hotline: '1800-123-456',
      email: 'hello@psyconnect.com',
      emergency: '113',
      socials: {
        facebook: 'https://facebook.com/psyconnect',
        instagram: 'https://instagram.com/psyconnect',
        linkedin: 'https://linkedin.com/company/psyconnect',
      },
      address: '36 Thanh Hoa St, Thanh Khe ,Da Nang City',
      businessHours: {
        weekdays: '8:00 - 20:00',
        weekends: '9:00 - 18:00',
      },
    };
  }

  getDevs(): Dev[] {
    return [
      {
        id: 'd1',
        name: 'Tran Van Huy',
        role: 'fullstack_lead',
        github: 'github.com/huyfullstack',
        linkedin: 'linkedin.com/in/huyfullstack',
        avatar: 'assets/team/dev1.jpg',
        bio: 'fullstack_lead_bio',
      },
      {
        id: 'd2',
        name: 'Nguyen Thi Nhu Quynh',
        role: 'uiux_designer',
        github: 'github.com/quynhux',
        linkedin: 'linkedin.com/in/quynhux',
        avatar: 'assets/team/dev2.jpg',
        bio: 'uiux_designer_bio',
      },
      {
        id: 'd3',
        name: 'Tran Huy',
        role: 'frontend_developer',
        github: 'github.com/minhfrontend',
        linkedin: 'linkedin.com/in/minhfrontend',
        avatar: 'assets/team/dev3.jpg',
        bio: 'frontend_developer_bio',
      },
      {
        id: 'd4',
        name: 'Phan Cong Danh',
        role: 'backend_developer',
        github: 'github.com/lanbackend',
        linkedin: 'linkedin.com/in/lanbackend',
        avatar: 'assets/team/dev4.jpg',
        bio: 'backend_developer_bio',
      },
      {
        id: 'd5',
        name: 'Phung Dinh Quang Huy',
        role: 'devops_engineer',
        github: 'github.com/hoangdevops',
        linkedin: 'linkedin.com/in/hoangdevops',
        avatar: 'assets/team/dev5.jpg',
        bio: 'devops_engineer_bio',
      },
      {
        id: 'd6',
        name: 'Chu Phuong Anh',
        role: 'qa_engineer',
        github: 'github.com/maitester',
        linkedin: 'linkedin.com/in/maitester',
        avatar: 'assets/team/dev6.jpg',
        bio: 'qa_engineer_bio',
      },
    ];
  }

  getRelease() {
    return {
      version: 'v2.1.0',
      date: '2025-09-09',
      notes: [],
      upcomingFeatures: [],
    };
  }
}
